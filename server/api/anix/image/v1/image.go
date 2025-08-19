package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ImageProcessReq 图片处理请求
type ImageProcessReq struct {
	g.Meta `path:"/image/process" method:"post" tags:"图片服务" summary:"图片处理"`
	Url     string `json:"url" v:"required#图片URL不能为空" dc:"原始图片URL"`
	Width   int    `json:"width" dc:"目标宽度，0表示保持原比例"`
	Height  int    `json:"height" dc:"目标高度，0表示保持原比例"`
	Quality int    `json:"quality" dc:"图片质量，1-100，默认80"`
	Format  string `json:"format" dc:"输出格式：jpeg, png, webp，默认webp"`
}

// ImageProcessRes 图片处理响应
type ImageProcessRes struct {
	Url      string `json:"url" dc:"处理后的图片URL"`
	Width    int    `json:"width" dc:"实际宽度"`
	Height   int    `json:"height" dc:"实际高度"`
	Size     int64  `json:"size" dc:"文件大小（字节）"`
	Format   string `json:"format" dc:"图片格式"`
	CacheKey string `json:"cache_key" dc:"缓存键"`
}

// ImagePreloadReq 图片预加载请求
type ImagePreloadReq struct {
	g.Meta `path:"/image/preload" method:"post" tags:"图片服务" summary:"图片预加载"`
	Urls   []string `json:"urls" v:"required#图片URL列表不能为空" dc:"需要预加载的图片URL列表"`
	Width  int      `json:"width" dc:"预处理宽度"`
	Height int      `json:"height" dc:"预处理高度"`
}

// ImagePreloadRes 图片预加载响应
type ImagePreloadRes struct {
	Success []PreloadResult `json:"success" dc:"成功预加载的图片"`
	Failed  []PreloadResult `json:"failed" dc:"预加载失败的图片"`
	Total   int             `json:"total" dc:"总数"`
}

// PreloadResult 预加载结果
type PreloadResult struct {
	OriginalUrl string `json:"original_url" dc:"原始URL"`
	ProcessedUrl string `json:"processed_url,omitempty" dc:"处理后URL"`
	Error       string `json:"error,omitempty" dc:"错误信息"`
	CacheKey    string `json:"cache_key,omitempty" dc:"缓存键"`
}

// ImageCacheInfoReq 图片缓存信息请求
type ImageCacheInfoReq struct {
	g.Meta   `path:"/image/cache/info" method:"get" tags:"图片服务" summary:"获取图片缓存信息"`
	CacheKey string `json:"cache_key" v:"required#缓存键不能为空" dc:"图片缓存键"`
}

// ImageCacheInfoRes 图片缓存信息响应
type ImageCacheInfoRes struct {
	Exists    bool   `json:"exists" dc:"缓存是否存在"`
	Url       string `json:"url,omitempty" dc:"图片URL"`
	Size      int64  `json:"size,omitempty" dc:"文件大小"`
	Format    string `json:"format,omitempty" dc:"图片格式"`
	CreatedAt string `json:"created_at,omitempty" dc:"创建时间"`
	ExpiresAt string `json:"expires_at,omitempty" dc:"过期时间"`
}

// ImageCacheClearReq 清理图片缓存请求
type ImageCacheClearReq struct {
	g.Meta    `path:"/image/cache/clear" method:"post" tags:"图片服务" summary:"清理图片缓存"`
	CacheKeys []string `json:"cache_keys" dc:"要清理的缓存键列表，为空则清理所有"`
	Expired   bool     `json:"expired" dc:"是否只清理过期缓存"`
}

// ImageCacheClearRes 清理图片缓存响应
type ImageCacheClearRes struct {
	Cleared int   `json:"cleared" dc:"已清理的缓存数量"`
	Errors  []string `json:"errors,omitempty" dc:"清理过程中的错误"`
}

// ImageStatsReq 图片服务统计请求
type ImageStatsReq struct {
	g.Meta `path:"/image/stats" method:"get" tags:"图片服务" summary:"获取图片服务统计信息"`
}

// ImageStatsRes 图片服务统计响应
type ImageStatsRes struct {
	TotalProcessed   int64 `json:"total_processed" dc:"总处理次数"`
	CacheHits        int64 `json:"cache_hits" dc:"缓存命中次数"`
	CacheMisses      int64 `json:"cache_misses" dc:"缓存未命中次数"`
	CacheHitRate     float64 `json:"cache_hit_rate" dc:"缓存命中率"`
	TotalCacheSize   int64 `json:"total_cache_size" dc:"总缓存大小（字节）"`
	CacheCount       int   `json:"cache_count" dc:"缓存文件数量"`
	AverageFileSize  int64 `json:"average_file_size" dc:"平均文件大小（字节）"`
	ProcessingErrors int64 `json:"processing_errors" dc:"处理错误次数"`
}

// ChapterImagesReq 获取章节图片列表请求
type ChapterImagesReq struct {
	g.Meta    `path:"/image/chapter/:id" method:"get" tags:"图片服务" summary:"获取章节图片列表"`
	Id        int64  `json:"id" v:"required#章节ID不能为空" dc:"章节ID"`
	Width     int    `json:"width" dc:"图片宽度"`
	Height    int    `json:"height" dc:"图片高度"`
	Quality   int    `json:"quality" dc:"图片质量"`
	Format    string `json:"format" dc:"图片格式"`
	Preload   bool   `json:"preload" dc:"是否预加载"`
}

// ChapterImagesRes 获取章节图片列表响应
type ChapterImagesRes struct {
	ChapterId    int64       `json:"chapter_id" dc:"章节ID"`
	ChapterTitle string      `json:"chapter_title" dc:"章节标题"`
	Images       []ImageInfo `json:"images" dc:"图片列表"`
	TotalPages   int         `json:"total_pages" dc:"总页数"`
	Preloaded    bool        `json:"preloaded" dc:"是否已预加载"`
}

// ImageInfo 图片信息
type ImageInfo struct {
	PageIndex    int    `json:"page_index" dc:"页码索引"`
	OriginalUrl  string `json:"original_url" dc:"原始图片URL"`
	ProcessedUrl string `json:"processed_url" dc:"处理后图片URL"`
	Width        int    `json:"width" dc:"图片宽度"`
	Height       int    `json:"height" dc:"图片高度"`
	Size         int64  `json:"size" dc:"文件大小"`
	Format       string `json:"format" dc:"图片格式"`
	CacheKey     string `json:"cache_key" dc:"缓存键"`
}