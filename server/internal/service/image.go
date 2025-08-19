package service

import (
	"context"
)

// IImage 图片服务接口
type IImage interface {
	// ProcessImage 处理图片
	ProcessImage(ctx context.Context, input *ImageProcessInput) (*ImageProcessOutput, error)
	// PreloadImages 预加载图片
	PreloadImages(ctx context.Context, input *ImagePreloadInput) (*ImagePreloadOutput, error)
	// GetCacheInfo 获取缓存信息
	GetCacheInfo(ctx context.Context, cacheKey string) (*ImageCacheInfo, error)
	// ClearCache 清理缓存
	ClearCache(ctx context.Context, input *ImageCacheClearInput) (*ImageCacheClearOutput, error)
	// GetStats 获取统计信息
	GetStats(ctx context.Context) (*ImageStatsOutput, error)
	// GetChapterImages 获取章节图片列表
	GetChapterImages(ctx context.Context, input *ChapterImagesInput) (*ChapterImagesOutput, error)
}

// ImageProcessInput 图片处理输入
type ImageProcessInput struct {
	Url     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Quality int    `json:"quality"`
	Format  string `json:"format"`
}

// ImageProcessOutput 图片处理输出
type ImageProcessOutput struct {
	Url      string `json:"url"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int64  `json:"size"`
	Format   string `json:"format"`
	CacheKey string `json:"cache_key"`
}

// ImagePreloadInput 图片预加载输入
type ImagePreloadInput struct {
	Urls   []string `json:"urls"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
}

// ImagePreloadOutput 图片预加载输出
type ImagePreloadOutput struct {
	Success []PreloadResultItem `json:"success"`
	Failed  []PreloadResultItem `json:"failed"`
	Total   int                 `json:"total"`
}

// PreloadResultItem 预加载结果项
type PreloadResultItem struct {
	OriginalUrl  string `json:"original_url"`
	ProcessedUrl string `json:"processed_url,omitempty"`
	Error        string `json:"error,omitempty"`
	CacheKey     string `json:"cache_key,omitempty"`
}

// ImageCacheInfo 图片缓存信息
type ImageCacheInfo struct {
	Exists    bool   `json:"exists"`
	Url       string `json:"url,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Format    string `json:"format,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// ImageCacheClearInput 图片缓存清理输入
type ImageCacheClearInput struct {
	CacheKeys []string `json:"cache_keys"`
	Expired   bool     `json:"expired"`
}

// ImageCacheClearOutput 图片缓存清理输出
type ImageCacheClearOutput struct {
	Cleared int      `json:"cleared"`
	Errors  []string `json:"errors,omitempty"`
}

// ImageStatsOutput 图片服务统计输出
type ImageStatsOutput struct {
	TotalProcessed   int64   `json:"total_processed"`
	CacheHits        int64   `json:"cache_hits"`
	CacheMisses      int64   `json:"cache_misses"`
	CacheHitRate     float64 `json:"cache_hit_rate"`
	TotalCacheSize   int64   `json:"total_cache_size"`
	CacheCount       int     `json:"cache_count"`
	AverageFileSize  int64   `json:"average_file_size"`
	ProcessingErrors int64   `json:"processing_errors"`
}

// ChapterImagesInput 章节图片输入
type ChapterImagesInput struct {
	ChapterId int64  `json:"chapter_id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Quality   int    `json:"quality"`
	Format    string `json:"format"`
	Preload   bool   `json:"preload"`
}

// ChapterImagesOutput 章节图片输出
type ChapterImagesOutput struct {
	ChapterId    int64           `json:"chapter_id"`
	ChapterTitle string          `json:"chapter_title"`
	Images       []ImageInfoItem `json:"images"`
	TotalPages   int             `json:"total_pages"`
	Preloaded    bool            `json:"preloaded"`
}

// ImageInfoItem 图片信息项
type ImageInfoItem struct {
	PageIndex    int    `json:"page_index"`
	OriginalUrl  string `json:"original_url"`
	ProcessedUrl string `json:"processed_url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Size         int64  `json:"size"`
	Format       string `json:"format"`
	CacheKey     string `json:"cache_key"`
}

var localImage IImage

// Image 获取图片服务实例
func Image() IImage {
	if localImage == nil {
		panic("implement not found for interface IImage, forgot register?")
	}
	return localImage
}

// RegisterImage 注册图片服务实现
func RegisterImage(i IImage) {
	localImage = i
}