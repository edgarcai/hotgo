package history

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/anix/history/v1"
	"hotgo/internal/library/contexts"
	"hotgo/internal/service"
)

type HistoryV1 struct{}

func NewV1() *HistoryV1 {
	return &HistoryV1{}
}

// SaveReadingProgress 保存阅读进度
func (c *HistoryV1) SaveReadingProgress(ctx context.Context, req *v1.SaveReadingProgressReq) (res *v1.SaveReadingProgressRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层保存阅读进度
	err = service.History().SaveReadingProgress(ctx, userId, req)
	if err != nil {
		return nil, err
	}

	return &v1.SaveReadingProgressRes{}, nil
}

// GetReadingHistory 获取阅读历史
func (c *HistoryV1) GetReadingHistory(ctx context.Context, req *v1.GetReadingHistoryReq) (res *v1.GetReadingHistoryRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 调用服务层获取阅读历史
	return service.History().GetReadingHistory(ctx, userId, req)
}

// GetReadingProgress 获取阅读进度
func (c *HistoryV1) GetReadingProgress(ctx context.Context, req *v1.GetReadingProgressReq) (res *v1.GetReadingProgressRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层获取阅读进度
	return service.History().GetReadingProgress(ctx, userId, req.ComicId, req.ChapterId)
}

// AddBookmark 添加书签
func (c *HistoryV1) AddBookmark(ctx context.Context, req *v1.AddBookmarkReq) (res *v1.AddBookmarkRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层添加书签
	id, err := service.History().AddBookmark(ctx, userId, req)
	if err != nil {
		return nil, err
	}

	return &v1.AddBookmarkRes{Id: id}, nil
}

// GetBookmarks 获取书签列表
func (c *HistoryV1) GetBookmarks(ctx context.Context, req *v1.GetBookmarksReq) (res *v1.GetBookmarksRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 调用服务层获取书签列表
	return service.History().GetBookmarks(ctx, userId, req)
}

// DeleteBookmark 删除书签
func (c *HistoryV1) DeleteBookmark(ctx context.Context, req *v1.DeleteBookmarkReq) (res *v1.DeleteBookmarkRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层删除书签
	err = service.History().DeleteBookmark(ctx, userId, req.Id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteBookmarkRes{}, nil
}

// GetReadingStats 获取阅读统计
func (c *HistoryV1) GetReadingStats(ctx context.Context, req *v1.GetReadingStatsReq) (res *v1.GetReadingStatsRes, err error) {
	// 获取当前用户ID
	userId := contexts.GetUserId(ctx)
	if userId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 调用服务层获取阅读统计
	return service.History().GetReadingStats(ctx, userId, req.ComicId)
}