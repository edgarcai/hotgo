package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

// CacheService 缓存服务
type CacheService struct {
	redis *gredis.Redis
}

// NewCacheService 创建缓存服务实例
func NewCacheService() *CacheService {
	return &CacheService{
		redis: g.Redis(),
	}
}

// ImageCacheKey 图片缓存键前缀
const (
	ImageCachePrefix     = "anix:image:"
	ChapterCachePrefix   = "anix:chapter:"
	UserCachePrefix      = "anix:user:"
	StatsCachePrefix     = "anix:stats:"
	DefaultCacheExpire   = 3600 * time.Second // 1小时
	ImageCacheExpire     = 7200 * time.Second // 2小时
	ChapterCacheExpire   = 1800 * time.Second // 30分钟
)

// SetImageCache 设置图片缓存
func (s *CacheService) SetImageCache(ctx context.Context, key string, data interface{}, expire time.Duration) error {
	cacheKey := ImageCachePrefix + key
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	if expire == 0 {
		expire = ImageCacheExpire
	}
	
	err = s.redis.Set(ctx, cacheKey, jsonData, expire)
	if err != nil {
		g.Log().Error(ctx, "设置图片缓存失败:", err)
		return err
	}
	
	g.Log().Debug(ctx, "设置图片缓存成功:", cacheKey)
	return nil
}

// GetImageCache 获取图片缓存
func (s *CacheService) GetImageCache(ctx context.Context, key string, result interface{}) error {
	cacheKey := ImageCachePrefix + key
	data, err := s.redis.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Debug(ctx, "获取图片缓存失败:", err)
		return err
	}
	
	if data.IsNil() {
		return fmt.Errorf("缓存不存在")
	}
	
	err = json.Unmarshal(data.Bytes(), result)
	if err != nil {
		g.Log().Error(ctx, "解析图片缓存数据失败:", err)
		return err
	}
	
	g.Log().Debug(ctx, "获取图片缓存成功:", cacheKey)
	return nil
}

// DeleteImageCache 删除图片缓存
func (s *CacheService) DeleteImageCache(ctx context.Context, key string) error {
	cacheKey := ImageCachePrefix + key
	err := s.redis.Del(ctx, cacheKey)
	if err != nil {
		g.Log().Error(ctx, "删除图片缓存失败:", err)
		return err
	}
	
	g.Log().Debug(ctx, "删除图片缓存成功:", cacheKey)
	return nil
}

// SetChapterCache 设置章节缓存
func (s *CacheService) SetChapterCache(ctx context.Context, key string, data interface{}, expire time.Duration) error {
	cacheKey := ChapterCachePrefix + key
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	if expire == 0 {
		expire = ChapterCacheExpire
	}
	
	err = s.redis.Set(ctx, cacheKey, jsonData, expire)
	if err != nil {
		g.Log().Error(ctx, "设置章节缓存失败:", err)
		return err
	}
	
	g.Log().Debug(ctx, "设置章节缓存成功:", cacheKey)
	return nil
}

// GetChapterCache 获取章节缓存
func (s *CacheService) GetChapterCache(ctx context.Context, key string, result interface{}) error {
	cacheKey := ChapterCachePrefix + key
	data, err := s.redis.Get(ctx, cacheKey)
	if err != nil {
		g.Log().Debug(ctx, "获取章节缓存失败:", err)
		return err
	}
	
	if data.IsNil() {
		return fmt.Errorf("缓存不存在")
	}
	
	err = json.Unmarshal(data.Bytes(), result)
	if err != nil {
		g.Log().Error(ctx, "解析章节缓存数据失败:", err)
		return err
	}
	
	g.Log().Debug(ctx, "获取章节缓存成功:", cacheKey)
	return nil
}

// ClearCacheByPattern 根据模式清理缓存
func (s *CacheService) ClearCacheByPattern(ctx context.Context, pattern string) error {
	keys, err := s.redis.Keys(ctx, pattern)
	if err != nil {
		g.Log().Error(ctx, "获取缓存键失败:", err)
		return err
	}
	
	if len(keys) == 0 {
		g.Log().Debug(ctx, "没有找到匹配的缓存键:", pattern)
		return nil
	}
	
	err = s.redis.Del(ctx, keys...)
	if err != nil {
		g.Log().Error(ctx, "批量删除缓存失败:", err)
		return err
	}
	
	g.Log().Info(ctx, fmt.Sprintf("成功清理 %d 个缓存项, 模式: %s", len(keys), pattern))
	return nil
}

// GetCacheStats 获取缓存统计信息
func (s *CacheService) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 获取图片缓存数量
	imageKeys, err := s.redis.Keys(ctx, ImageCachePrefix+"*")
	if err != nil {
		g.Log().Error(ctx, "获取图片缓存键失败:", err)
	} else {
		stats["image_cache_count"] = len(imageKeys)
	}
	
	// 获取章节缓存数量
	chapterKeys, err := s.redis.Keys(ctx, ChapterCachePrefix+"*")
	if err != nil {
		g.Log().Error(ctx, "获取章节缓存键失败:", err)
	} else {
		stats["chapter_cache_count"] = len(chapterKeys)
	}
	
	// 获取Redis信息
	info, err := s.redis.Info(ctx)
	if err != nil {
		g.Log().Error(ctx, "获取Redis信息失败:", err)
	} else {
		stats["redis_info"] = info
	}
	
	return stats, nil
}