// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package anix

import (
	"context"

	"hotgo/api/anix/chapter/v1"
)

type IChapterV1 interface {
	// 获取章节详情
	GetDetail(ctx context.Context, req *v1.GetDetailReq) (res *v1.GetDetailRes, err error)
	// 获取章节页面列表
	GetPages(ctx context.Context, req *v1.GetPagesReq) (res *v1.GetPagesRes, err error)
	// 记录阅读历史
	RecordReading(ctx context.Context, req *v1.RecordReadingReq) (res *v1.RecordReadingRes, err error)
	// 获取上一章节
	GetPrevious(ctx context.Context, req *v1.GetPreviousReq) (res *v1.GetPreviousRes, err error)
	// 获取下一章节
	GetNext(ctx context.Context, req *v1.GetNextReq) (res *v1.GetNextRes, err error)
}