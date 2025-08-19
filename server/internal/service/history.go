package service

import (
	"context"
	"hotgo/api/anix/history/v1"
)

type (
	// IHistory 阅读历史服务接口
	IHistory interface {
		// SaveReadingProgress 保存阅读进度
		SaveReadingProgress(ctx context.Context, userId int64, req *v1.SaveReadingProgressReq) error
		// GetReadingHistory 获取阅读历史
		GetReadingHistory(ctx context.Context, userId int64, req *v1.GetReadingHistoryReq) (*v1.GetReadingHistoryRes, error)
		// GetReadingProgress 获取阅读进度
		GetReadingProgress(ctx context.Context, userId int64, comicId, chapterId int64) (*v1.GetReadingProgressRes, error)
		// AddBookmark 添加书签
		AddBookmark(ctx context.Context, userId int64, req *v1.AddBookmarkReq) (int64, error)
		// GetBookmarks 获取书签列表
		GetBookmarks(ctx context.Context, userId int64, req *v1.GetBookmarksReq) (*v1.GetBookmarksRes, error)
		// DeleteBookmark 删除书签
		DeleteBookmark(ctx context.Context, userId int64, bookmarkId int64) error
		// GetReadingStats 获取阅读统计
		GetReadingStats(ctx context.Context, userId int64, comicId int64) (*v1.GetReadingStatsRes, error)
	}
)

var (
	localHistory IHistory
)

// History 返回阅读历史服务实例
func History() IHistory {
	if localHistory == nil {
		panic("implement not found for interface IHistory, forgot register?")
	}
	return localHistory
}

// RegisterHistory 注册阅读历史服务实现
func RegisterHistory(i IHistory) {
	localHistory = i
}