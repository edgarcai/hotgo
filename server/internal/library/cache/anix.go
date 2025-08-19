// Package cache
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

// AniX应用缓存键前缀
const (
	AnixComicListPrefix     = "anix:comic:list:%s"      // 漫画列表缓存
	AnixComicDetailPrefix   = "anix:comic:detail:%d"    // 漫画详情缓存
	AnixComicSearchPrefix   = "anix:comic:search:%s"    // 搜索结果缓存
	AnixChapterDetailPrefix = "anix:chapter:detail:%d"  // 章节详情缓存
	AnixChapterPagesPrefix  = "anix:chapter:pages:%d"   // 章节页面缓存
	AnixUserProfilePrefix   = "anix:user:profile:%d"    // 用户信息缓存
	AnixReadingHistoryPrefix = "anix:reading:history:%d" // 阅读历史缓存
)

// AniX应用缓存过期时间
const (
	AnixComicListExpire     = 30 * time.Minute  // 漫画列表缓存30分钟
	AnixComicDetailExpire   = 60 * time.Minute  // 漫画详情缓存1小时
	AnixComicSearchExpire   = 15 * time.Minute  // 搜索结果缓存15分钟
	AnixChapterDetailExpire = 60 * time.Minute  // 章节详情缓存1小时
	AnixChapterPagesExpire  = 120 * time.Minute // 章节页面缓存2小时
	AnixUserProfileExpire   = 30 * time.Minute  // 用户信息缓存30分钟
	AnixReadingHistoryExpire = 60 * time.Minute // 阅读历史缓存1小时
)

// AnixCache AniX应用缓存服务
type AnixCache struct {
	cache *gcache.Cache
}

// NewAnixCache 创建AniX缓存服务实例
func NewAnixCache() *AnixCache {
	return &AnixCache{
		cache: Instance(),
	}
}

// SetComicList 设置漫画列表缓存
func (c *AnixCache) SetComicList(ctx context.Context, key string, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixComicListPrefix, key)
	return c.cache.Set(ctx, cacheKey, data, AnixComicListExpire)
}

// GetComicList 获取漫画列表缓存
func (c *AnixCache) GetComicList(ctx context.Context, key string) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixComicListPrefix, key)
	return c.cache.Get(ctx, cacheKey)
}

// SetComicDetail 设置漫画详情缓存
func (c *AnixCache) SetComicDetail(ctx context.Context, comicId int64, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixComicDetailPrefix, comicId)
	return c.cache.Set(ctx, cacheKey, data, AnixComicDetailExpire)
}

// GetComicDetail 获取漫画详情缓存
func (c *AnixCache) GetComicDetail(ctx context.Context, comicId int64) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixComicDetailPrefix, comicId)
	return c.cache.Get(ctx, cacheKey)
}

// SetComicSearch 设置搜索结果缓存
func (c *AnixCache) SetComicSearch(ctx context.Context, keyword string, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixComicSearchPrefix, keyword)
	return c.cache.Set(ctx, cacheKey, data, AnixComicSearchExpire)
}

// GetComicSearch 获取搜索结果缓存
func (c *AnixCache) GetComicSearch(ctx context.Context, keyword string) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixComicSearchPrefix, keyword)
	return c.cache.Get(ctx, cacheKey)
}

// SetChapterDetail 设置章节详情缓存
func (c *AnixCache) SetChapterDetail(ctx context.Context, chapterId int64, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixChapterDetailPrefix, chapterId)
	return c.cache.Set(ctx, cacheKey, data, AnixChapterDetailExpire)
}

// GetChapterDetail 获取章节详情缓存
func (c *AnixCache) GetChapterDetail(ctx context.Context, chapterId int64) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixChapterDetailPrefix, chapterId)
	return c.cache.Get(ctx, cacheKey)
}

// SetChapterPages 设置章节页面缓存
func (c *AnixCache) SetChapterPages(ctx context.Context, chapterId int64, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixChapterPagesPrefix, chapterId)
	return c.cache.Set(ctx, cacheKey, data, AnixChapterPagesExpire)
}

// GetChapterPages 获取章节页面缓存
func (c *AnixCache) GetChapterPages(ctx context.Context, chapterId int64) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixChapterPagesPrefix, chapterId)
	return c.cache.Get(ctx, cacheKey)
}

// SetUserProfile 设置用户信息缓存
func (c *AnixCache) SetUserProfile(ctx context.Context, userId int64, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixUserProfilePrefix, userId)
	return c.cache.Set(ctx, cacheKey, data, AnixUserProfileExpire)
}

// GetUserProfile 获取用户信息缓存
func (c *AnixCache) GetUserProfile(ctx context.Context, userId int64) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixUserProfilePrefix, userId)
	return c.cache.Get(ctx, cacheKey)
}

// SetReadingHistory 设置阅读历史缓存
func (c *AnixCache) SetReadingHistory(ctx context.Context, userId int64, data interface{}) error {
	cacheKey := fmt.Sprintf(AnixReadingHistoryPrefix, userId)
	return c.cache.Set(ctx, cacheKey, data, AnixReadingHistoryExpire)
}

// GetReadingHistory 获取阅读历史缓存
func (c *AnixCache) GetReadingHistory(ctx context.Context, userId int64) (interface{}, error) {
	cacheKey := fmt.Sprintf(AnixReadingHistoryPrefix, userId)
	return c.cache.Get(ctx, cacheKey)
}

// DeleteComicCache 删除漫画相关缓存
func (c *AnixCache) DeleteComicCache(ctx context.Context, comicId int64) error {
	// 删除漫画详情缓存
	detailKey := fmt.Sprintf(AnixComicDetailPrefix, comicId)
	if _, err := c.cache.Remove(ctx, detailKey); err != nil {
		g.Log().Errorf(ctx, "删除漫画详情缓存失败: %v", err)
	}

	// 删除相关的列表缓存（这里可以根据实际需要扩展）
	// 由于列表缓存的key是动态的，这里只是示例
	return nil
}

// DeleteChapterCache 删除章节相关缓存
func (c *AnixCache) DeleteChapterCache(ctx context.Context, chapterId int64) error {
	// 删除章节详情缓存
	detailKey := fmt.Sprintf(AnixChapterDetailPrefix, chapterId)
	if _, err := c.cache.Remove(ctx, detailKey); err != nil {
		g.Log().Errorf(ctx, "删除章节详情缓存失败: %v", err)
	}

	// 删除章节页面缓存
	pagesKey := fmt.Sprintf(AnixChapterPagesPrefix, chapterId)
	if _, err := c.cache.Remove(ctx, pagesKey); err != nil {
		g.Log().Errorf(ctx, "删除章节页面缓存失败: %v", err)
	}

	return nil
}

// DeleteUserCache 删除用户相关缓存
func (c *AnixCache) DeleteUserCache(ctx context.Context, userId int64) error {
	// 删除用户信息缓存
	profileKey := fmt.Sprintf(AnixUserProfilePrefix, userId)
	if _, err := c.cache.Remove(ctx, profileKey); err != nil {
		g.Log().Errorf(ctx, "删除用户信息缓存失败: %v", err)
	}

	// 删除阅读历史缓存
	historyKey := fmt.Sprintf(AnixReadingHistoryPrefix, userId)
	if _, err := c.cache.Remove(ctx, historyKey); err != nil {
		g.Log().Errorf(ctx, "删除阅读历史缓存失败: %v", err)
	}

	return nil
}