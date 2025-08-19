// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package anix

import (
	"context"

	"hotgo/api/anix/comic/v1"
)

type IComicV1 interface {
	// 获取漫画列表
	GetList(ctx context.Context, req *v1.GetListReq) (res *v1.GetListRes, err error)
	// 获取漫画详情
	GetDetail(ctx context.Context, req *v1.GetDetailReq) (res *v1.GetDetailRes, err error)
	// 搜索漫画
	Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error)
	// 获取推荐漫画
	GetRecommend(ctx context.Context, req *v1.GetRecommendReq) (res *v1.GetRecommendRes, err error)
	// 获取热门漫画
	GetHot(ctx context.Context, req *v1.GetHotReq) (res *v1.GetHotRes, err error)
	// 获取最新漫画
	GetLatest(ctx context.Context, req *v1.GetLatestReq) (res *v1.GetLatestRes, err error)
	// 按分类获取漫画
	GetByCategory(ctx context.Context, req *v1.GetByCategoryReq) (res *v1.GetByCategoryRes, err error)
}