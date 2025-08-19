package image

import (
	"context"
	"strconv"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	v1 "hotgo/api/anix/image/v1"
	"hotgo/internal/service"
)

// ImageV1 图片服务控制器v1
type ImageV1 struct{}

var Image = cImage{}

type cImage struct{}

// ProcessImage 图片处理
func (c *cImage) ProcessImage(ctx context.Context, req *v1.ImageProcessReq) (res *v1.ImageProcessRes, err error) {
	// 参数验证
	if req.Quality <= 0 || req.Quality > 100 {
		req.Quality = 80 // 默认质量
	}
	if req.Format == "" {
		req.Format = "webp" // 默认格式
	}

	// 调用服务层处理图片
	result, err := service.Image().ProcessImage(ctx, &service.ImageProcessInput{
		Url:     req.Url,
		Width:   req.Width,
		Height:  req.Height,
		Quality: req.Quality,
		Format:  req.Format,
	})
	if err != nil {
		return nil, gerror.New("图片处理失败: " + err.Error())
	}

	res = &v1.ImageProcessRes{
		Url:      result.Url,
		Width:    result.Width,
		Height:   result.Height,
		Size:     result.Size,
		Format:   result.Format,
		CacheKey: result.CacheKey,
	}
	return
}

// PreloadImages 图片预加载
func (c *cImage) PreloadImages(ctx context.Context, req *v1.ImagePreloadReq) (res *v1.ImagePreloadRes, err error) {
	if len(req.Urls) == 0 {
		return nil, gerror.New("图片URL列表不能为空")
	}

	// 调用服务层预加载图片
	result, err := service.Image().PreloadImages(ctx, &service.ImagePreloadInput{
		Urls:   req.Urls,
		Width:  req.Width,
		Height: req.Height,
	})
	if err != nil {
		return nil, gerror.New("图片预加载失败: " + err.Error())
	}

	res = &v1.ImagePreloadRes{
		Success: make([]v1.PreloadResult, len(result.Success)),
		Failed:  make([]v1.PreloadResult, len(result.Failed)),
		Total:   result.Total,
	}

	// 转换成功结果
	for i, item := range result.Success {
		res.Success[i] = v1.PreloadResult{
			OriginalUrl:  item.OriginalUrl,
			ProcessedUrl: item.ProcessedUrl,
			CacheKey:     item.CacheKey,
		}
	}

	// 转换失败结果
	for i, item := range result.Failed {
		res.Failed[i] = v1.PreloadResult{
			OriginalUrl: item.OriginalUrl,
			Error:       item.Error,
		}
	}

	return
}

// GetCacheInfo 获取图片缓存信息
func (c *cImage) GetCacheInfo(ctx context.Context, req *v1.ImageCacheInfoReq) (res *v1.ImageCacheInfoRes, err error) {
	if req.CacheKey == "" {
		return nil, gerror.New("缓存键不能为空")
	}

	// 调用服务层获取缓存信息
	result, err := service.Image().GetCacheInfo(ctx, req.CacheKey)
	if err != nil {
		return nil, gerror.New("获取缓存信息失败: " + err.Error())
	}

	res = &v1.ImageCacheInfoRes{
		Exists:    result.Exists,
		Url:       result.Url,
		Size:      result.Size,
		Format:    result.Format,
		CreatedAt: result.CreatedAt,
		ExpiresAt: result.ExpiresAt,
	}
	return
}

// ClearCache 清理图片缓存
func (c *cImage) ClearCache(ctx context.Context, req *v1.ImageCacheClearReq) (res *v1.ImageCacheClearRes, err error) {
	// 调用服务层清理缓存
	result, err := service.Image().ClearCache(ctx, &service.ImageCacheClearInput{
		CacheKeys: req.CacheKeys,
		Expired:   req.Expired,
	})
	if err != nil {
		return nil, gerror.New("清理缓存失败: " + err.Error())
	}

	res = &v1.ImageCacheClearRes{
		Cleared: result.Cleared,
		Errors:  result.Errors,
	}
	return
}

// GetStats 获取图片服务统计信息
func (c *cImage) GetStats(ctx context.Context, req *v1.ImageStatsReq) (res *v1.ImageStatsRes, err error) {
	// 调用服务层获取统计信息
	result, err := service.Image().GetStats(ctx)
	if err != nil {
		return nil, gerror.New("获取统计信息失败: " + err.Error())
	}

	res = &v1.ImageStatsRes{
		TotalProcessed:   result.TotalProcessed,
		CacheHits:        result.CacheHits,
		CacheMisses:      result.CacheMisses,
		CacheHitRate:     result.CacheHitRate,
		TotalCacheSize:   result.TotalCacheSize,
		CacheCount:       result.CacheCount,
		AverageFileSize:  result.AverageFileSize,
		ProcessingErrors: result.ProcessingErrors,
	}
	return
}

// GetChapterImages 获取章节图片列表
func (c *cImage) GetChapterImages(ctx context.Context, req *v1.ChapterImagesReq) (res *v1.ChapterImagesRes, err error) {
	// 从URL路径中获取章节ID
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		idStr := r.Get("id").String()
		if idStr != "" {
			if id, parseErr := strconv.ParseInt(idStr, 10, 64); parseErr == nil {
				req.Id = id
			}
		}
	}

	if req.Id <= 0 {
		return nil, gerror.New("章节ID不能为空")
	}

	// 设置默认参数
	if req.Quality <= 0 || req.Quality > 100 {
		req.Quality = 80
	}
	if req.Format == "" {
		req.Format = "webp"
	}

	// 调用服务层获取章节图片
	result, err := service.Image().GetChapterImages(ctx, &service.ChapterImagesInput{
		ChapterId: req.Id,
		Width:     req.Width,
		Height:    req.Height,
		Quality:   req.Quality,
		Format:    req.Format,
		Preload:   req.Preload,
	})
	if err != nil {
		return nil, gerror.New("获取章节图片失败: " + err.Error())
	}

	res = &v1.ChapterImagesRes{
		ChapterId:    result.ChapterId,
		ChapterTitle: result.ChapterTitle,
		Images:       make([]v1.ImageInfo, len(result.Images)),
		TotalPages:   result.TotalPages,
		Preloaded:    result.Preloaded,
	}

	// 转换图片信息
	for i, img := range result.Images {
		res.Images[i] = v1.ImageInfo{
			PageIndex:    img.PageIndex,
			OriginalUrl:  img.OriginalUrl,
			ProcessedUrl: img.ProcessedUrl,
			Width:        img.Width,
			Height:       img.Height,
			Size:         img.Size,
			Format:       img.Format,
			CacheKey:     img.CacheKey,
		}
	}

	return
}

// NewV1 创建图片服务控制器v1实例
func NewV1() *ImageV1 {
	return &ImageV1{}
}