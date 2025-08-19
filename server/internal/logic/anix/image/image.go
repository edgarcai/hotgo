package image

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/webp"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nfnt/resize"

	"hotgo/internal/dao"
	"hotgo/internal/logic/anix/cache"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

// sImage 图片服务实现
type sImage struct {
	cacheService *cache.CacheService
}

func init() {
	service.RegisterImage(New())
}

// New 创建图片服务实例
func New() service.IImage {
	return &sImage{
		cacheService: cache.NewCacheService(),
	}
}

// 缓存目录
const (
	CacheDir     = "storage/cache/images"
	MaxCacheSize = 1024 * 1024 * 1024 // 1GB
	CacheExpiry  = 24 * time.Hour     // 24小时
)

// 统计信息缓存
var (
	statsCache = struct {
		sync.RWMutex
		data      *service.ImageStatsOutput
		lastUpdate time.Time
	}{}
	cacheExpiry = 5 * time.Minute // 统计缓存5分钟
)

// ProcessImage 处理图片
func (s *sImage) ProcessImage(ctx context.Context, in *service.ProcessImageInp) (*service.ProcessImageOut, error) {
	g.Log().Info(ctx, "处理图片请求:", in)
	
	// 生成缓存键
	cacheKey := fmt.Sprintf("processed:%s:%d:%d:%s", in.ImageUrl, in.Width, in.Height, in.Format)
	
	// 尝试从缓存获取
	var cachedResult service.ProcessImageOut
	if err := s.cacheService.GetImageCache(ctx, cacheKey, &cachedResult); err == nil {
		g.Log().Debug(ctx, "从缓存获取处理后的图片:", cacheKey)
		return &cachedResult, nil
	}
	
	// 缓存未命中，执行图片处理
	processedUrl, fileSize, format, width, height, err := s.processImageFile(ctx, in)
	if err != nil {
		g.Log().Error(ctx, "图片处理失败:", err)
		return nil, err
	}
	
	result := &service.ProcessImageOut{
		ProcessedUrl: processedUrl,
		FileSize:     fileSize,
		Format:       format,
		Width:        width,
		Height:       height,
	}
	
	// 保存到缓存
	if err := s.cacheService.SetImageCache(ctx, cacheKey, result, 0); err != nil {
		g.Log().Error(ctx, "保存图片处理结果到缓存失败:", err)
	}
	
	return result, nil
}

// PreloadImages 预加载图片
func (s *sImage) PreloadImages(ctx context.Context, in *service.PreloadImagesInp) (*service.PreloadImagesOut, error) {
	g.Log().Info(ctx, "预加载图片请求:", in)
	
	// 获取章节图片列表
	images, err := s.GetChapterImages(ctx, &service.GetChapterImagesInp{
		ChapterId: in.ChapterId,
	})
	if err != nil {
		return nil, err
	}
	
	var preloadedCount int
	var cacheKeys []string
	
	// 批量预加载图片
	for _, img := range images.Images {
		cacheKey := fmt.Sprintf("preload:%d:%d", in.ChapterId, img.PageIndex)
		
		// 检查是否已缓存
		var cached interface{}
		if err := s.cacheService.GetImageCache(ctx, cacheKey, &cached); err == nil {
			preloadedCount++
			cacheKeys = append(cacheKeys, cacheKey)
			continue
		}
		
		// 下载并缓存图片
		if err := s.preloadSingleImage(ctx, img.ImageUrl, cacheKey); err != nil {
			g.Log().Error(ctx, "预加载图片失败:", img.ImageUrl, err)
			continue
		}
		
		preloadedCount++
		cacheKeys = append(cacheKeys, cacheKey)
	}
	
	return &service.PreloadImagesOut{
		PreloadedCount: preloadedCount,
		TotalCount:     len(images.Images),
		CacheKeys:      cacheKeys,
	}, nil
}

// GetCacheInfo 获取缓存信息
func (s *sImage) GetCacheInfo(ctx context.Context, in *service.GetCacheInfoInp) (*service.GetCacheInfoOut, error) {
	g.Log().Info(ctx, "获取缓存信息请求:", in)
	
	// 获取缓存统计信息
	stats, err := s.cacheService.GetCacheStats(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取缓存统计失败:", err)
		return nil, err
	}
	
	return &service.GetCacheInfoOut{
		TotalSize:   stats.TotalSize,
		TotalFiles:  stats.TotalKeys,
		HitRate:     stats.HitRate,
		LastCleanup: stats.LastCleanup,
	}, nil
}

// ClearCache 清理缓存
func (s *sImage) ClearCache(ctx context.Context, in *service.ClearCacheInp) (*service.ClearCacheOut, error) {
	g.Log().Info(ctx, "清理缓存请求:", in)
	
	// 执行缓存清理
	clearedCount, freedSpace, err := s.cacheService.ClearCacheByPattern(ctx, in.Pattern)
	if err != nil {
		g.Log().Error(ctx, "清理缓存失败:", err)
		return nil, err
	}
	
	return &service.ClearCacheOut{
		ClearedFiles: clearedCount,
		FreedSpace:   freedSpace,
	}, nil
}

// GetStats 获取统计信息
func (s *sImage) GetStats(ctx context.Context) (*service.ImageStatsOutput, error) {
	statsCache.RLock()
	if statsCache.data != nil && time.Since(statsCache.lastUpdate) < cacheExpiry {
		defer statsCache.RUnlock()
		return statsCache.data, nil
	}
	statsCache.RUnlock()
	
	// 重新计算统计信息
	stats := &service.ImageStatsOutput{}
	
	// 获取缓存目录信息
	if gfile.Exists(CacheDir) {
		files, _ := filepath.Glob(filepath.Join(CacheDir, "*"))
		stats.CacheCount = len(files)
		
		var totalSize int64
		for _, file := range files {
			if fileInfo, err := os.Stat(file); err == nil {
				totalSize += fileInfo.Size()
			}
		}
		stats.TotalCacheSize = totalSize
		
		if stats.CacheCount > 0 {
			stats.AverageFileSize = totalSize / int64(stats.CacheCount)
		}
	}
	
	// 计算缓存命中率（这里使用模拟数据，实际应该从数据库或Redis获取）
	if stats.CacheHits+stats.CacheMisses > 0 {
		stats.CacheHitRate = float64(stats.CacheHits) / float64(stats.CacheHits+stats.CacheMisses)
	}
	
	// 更新缓存
	statsCache.Lock()
	statsCache.data = stats
	statsCache.lastUpdate = time.Now()
	statsCache.Unlock()
	
	return stats, nil
}

// GetChapterImages 获取章节图片列表
func (s *sImage) GetChapterImages(ctx context.Context, in *service.GetChapterImagesInp) (*service.GetChapterImagesOut, error) {
	g.Log().Info(ctx, "获取章节图片列表请求:", in)
	
	// 从数据库查询章节图片
	var images []entity.AnixChapterImage
	err := dao.AnixChapterImage.Ctx(ctx).Where("chapter_id = ?", in.ChapterId).OrderAsc("page_index").Scan(&images)
	if err != nil {
		return nil, err
	}
	
	// 转换为输出格式
	var imageList []service.ChapterImage
	for _, img := range images {
		imageList = append(imageList, service.ChapterImage{
			PageIndex: img.PageIndex,
			ImageUrl:  img.ImageUrl,
			Width:     img.Width,
			Height:    img.Height,
			FileSize:  img.FileSize,
			Format:    img.Format,
		})
	}
	
	return &service.GetChapterImagesOut{
		Images: imageList,
	}, nil
}

// processImageFile 处理图片文件
func (s *sImage) processImageFile(ctx context.Context, in *service.ProcessImageInp) (string, int64, string, int, int, error) {
	// 下载原始图片
	imageData, err := s.downloadImage(in.ImageUrl)
	if err != nil {
		return "", 0, "", 0, 0, err
	}
	
	// 处理图片数据
	processedData, outputFormat, err := s.processImageData(imageData, in.Width, in.Height, 80, in.Format)
	if err != nil {
		return "", 0, "", 0, 0, err
	}
	
	// 解码处理后的图片获取尺寸
	img, _, err := image.Decode(bytes.NewReader(processedData))
	if err != nil {
		return "", 0, "", 0, 0, err
	}
	
	targetWidth := img.Bounds().Dx()
	targetHeight := img.Bounds().Dy()
	
	// 生成缓存键并保存
	cacheKey := s.generateCacheKey(in.ImageUrl, in.Width, in.Height, 80, in.Format)
	processedUrl, err := s.saveToCache(cacheKey, processedData)
	if err != nil {
		return "", 0, "", 0, 0, err
	}
	
	return s.getCacheUrl(cacheKey), int64(len(processedData)), outputFormat, targetWidth, targetHeight, nil
}

// preloadSingleImage 预加载单张图片
func (s *sImage) preloadSingleImage(ctx context.Context, imageUrl, cacheKey string) error {
	// 下载图片
	imageData, err := s.downloadImage(imageUrl)
	if err != nil {
		return err
	}
	
	// 保存到缓存
	return s.cacheService.SetImageCache(ctx, cacheKey, imageData, 0)
}

// downloadImage 下载图片
func (s *sImage) downloadImage(imageUrl string) ([]byte, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}
	
	return io.ReadAll(resp.Body)
}

// calculateNewSize 计算新的图片尺寸
func (s *sImage) calculateNewSize(originalWidth, originalHeight, targetWidth, targetHeight int) (int, int) {
	if targetWidth == 0 && targetHeight == 0 {
		return originalWidth, originalHeight
	}
	
	if targetWidth == 0 {
		// 只指定高度，按比例缩放
		ratio := float64(targetHeight) / float64(originalHeight)
		return int(float64(originalWidth) * ratio), targetHeight
	}
	
	if targetHeight == 0 {
		// 只指定宽度，按比例缩放
		ratio := float64(targetWidth) / float64(originalWidth)
		return targetWidth, int(float64(originalHeight) * ratio)
	}
	
	// 同时指定宽高
	return targetWidth, targetHeight
}

// processImageData 处理图片数据
func (s *sImage) processImageData(imageData []byte, width, height, quality int, format string) ([]byte, string, error) {
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, "", err
	}
	
	// 调整尺寸
	if width > 0 || height > 0 {
		newWidth, newHeight := s.calculateNewSize(img.Bounds().Dx(), img.Bounds().Dy(), width, height)
		img = resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	}
	
	// 编码图片
	var buf bytes.Buffer
	var outputFormat string
	
	switch format {
	case "webp":
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
		outputFormat = "webp"
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		outputFormat = "jpeg"
	case "png":
		err = png.Encode(&buf, img)
		outputFormat = "png"
	default:
		// 默认使用webp格式
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
		outputFormat = "webp"
	}
	
	if err != nil {
		return nil, "", err
	}
	
	return buf.Bytes(), outputFormat, nil
}

// generateCacheKey 生成缓存键
func (s *sImage) generateCacheKey(url string, width, height, quality int, format string) string {
	return fmt.Sprintf("%s_%d_%d_%d_%s", url, width, height, quality, format)
}

// saveToCache 保存到缓存
func (s *sImage) saveToCache(cacheKey string, data []byte) (string, error) {
	// 这里可以保存到本地文件系统或云存储
	// 返回访问URL
	return fmt.Sprintf("/cache/%s", cacheKey), nil
}

// getCacheUrl 获取缓存URL
func (s *sImage) getCacheUrl(cacheKey string) string {
	return fmt.Sprintf("/cache/%s", cacheKey)
}

// 辅助方法

// generateCacheKey 生成缓存键
func (s *sImage) generateCacheKey(url string, width, height, quality int, format string) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d_%d_%d_%s", url, width, height, quality, format)))
	return hex.EncodeToString(hash[:]) + "." + format
}

// checkCache 检查缓存是否存在
func (s *sImage) checkCache(cacheKey string) (string, bool) {
	cachePath := filepath.Join(CacheDir, cacheKey)
	if gfile.Exists(cachePath) {
		// 检查是否过期
		if fileInfo, err := os.Stat(cachePath); err == nil {
			if time.Since(fileInfo.ModTime()) < CacheExpiry {
				return cachePath, true
			}
		}
	}
	return "", false
}

// getCacheUrl 获取缓存URL
func (s *sImage) getCacheUrl(cacheKey string) string {
	return "/api/cache/images/" + cacheKey
}

// downloadImage 下载图片
func (s *sImage) downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, gerror.New("下载图片失败，状态码: " + strconv.Itoa(resp.StatusCode))
	}
	
	return io.ReadAll(resp.Body)
}

// processImageData 处理图片数据
func (s *sImage) processImageData(data []byte, width, height, quality int, format string) ([]byte, int, int, error) {
	// 解码图片
	img, _, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, 0, 0, err
	}
	
	originalBounds := img.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()
	
	// 计算新尺寸
	newWidth, newHeight := s.calculateNewSize(originalWidth, originalHeight, width, height)
	
	// 调整大小
	if newWidth != originalWidth || newHeight != originalHeight {
		img = resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	}
	
	// 编码图片
	var buf strings.Builder
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	case "png":
		err = png.Encode(&buf, img)
	case "webp":
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
	default:
		return nil, 0, 0, gerror.New("不支持的图片格式: " + format)
	}
	
	if err != nil {
		return nil, 0, 0, err
	}
	
	return []byte(buf.String()), newWidth, newHeight, nil
}

// calculateNewSize 计算新尺寸
func (s *sImage) calculateNewSize(originalWidth, originalHeight, targetWidth, targetHeight int) (int, int) {
	if targetWidth <= 0 && targetHeight <= 0 {
		return originalWidth, originalHeight
	}
	
	if targetWidth <= 0 {
		// 只指定高度
		ratio := float64(targetHeight) / float64(originalHeight)
		return int(float64(originalWidth) * ratio), targetHeight
	}
	
	if targetHeight <= 0 {
		// 只指定宽度
		ratio := float64(targetWidth) / float64(originalWidth)
		return targetWidth, int(float64(originalHeight) * ratio)
	}
	
	// 指定了宽度和高度
	return targetWidth, targetHeight
}

// saveToCache 保存到缓存
func (s *sImage) saveToCache(cacheKey string, data []byte) (string, error) {
	// 确保缓存目录存在
	if err := gfile.Mkdir(CacheDir); err != nil {
		return "", err
	}
	
	cachePath := filepath.Join(CacheDir, cacheKey)
	return cachePath, gfile.PutBytes(cachePath, data)
}

// getImageDimensions 获取图片尺寸
func (s *sImage) getImageDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	
	return img.Width, img.Height, nil
}

// getFormatFromCacheKey 从缓存键获取格式
func (s *sImage) getFormatFromCacheKey(cacheKey string) string {
	parts := strings.Split(cacheKey, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

// updateStats 更新统计信息
func (s *sImage) updateStats(statType string) {
	// 这里应该更新到数据库或Redis，暂时跳过实现
	// 可以考虑使用异步方式更新统计信息
}