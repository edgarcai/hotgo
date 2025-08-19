// Package v1
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ComicInfo 漫画信息
type ComicInfo struct {
	Id          int64  `json:"id" dc:"漫画ID"`
	Title       string `json:"title" dc:"漫画标题"`
	Cover       string `json:"cover" dc:"封面图片"`
	Description string `json:"description" dc:"漫画描述"`
	Author      string `json:"author" dc:"作者"`
	Status      int    `json:"status" dc:"状态：1连载中，2已完结，3暂停更新"`
	Category    string `json:"category" dc:"分类"`
	Tags        string `json:"tags" dc:"标签"`
	ViewCount   int64  `json:"view_count" dc:"浏览次数"`
	LikeCount   int64  `json:"like_count" dc:"点赞次数"`
}

// ListReq 获取漫画列表请求
type ListReq struct {
	g.Meta   `path:"/comic/list" method:"get" tags:"漫画" summary:"获取漫画列表"`
	Page     int    `json:"page" d:"1" v:"min:1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"20" v:"min:1|max:100" dc:"每页数量"`
	Keyword  string `json:"keyword" dc:"搜索关键词"`
	Category string `json:"category" dc:"分类"`
	Status   int    `json:"status" dc:"状态"`
}

// ListRes 获取漫画列表响应
type ListRes struct {
	List  []ComicInfo `json:"list" dc:"漫画列表"`
	Total int64       `json:"total" dc:"总数量"`
}

// DetailReq 获取漫画详情请求
type DetailReq struct {
	g.Meta `path:"/comic/detail/:id" method:"get" tags:"漫画" summary:"获取漫画详情"`
	Id     int64 `json:"id" v:"required#漫画ID不能为空" dc:"漫画ID"`
}

// DetailRes 获取漫画详情响应
type DetailRes struct {
	Comic ComicInfo `json:"comic" dc:"漫画信息"`
}

// SearchReq 搜索漫画请求
type SearchReq struct {
	g.Meta   `path:"/comic/search" method:"get" tags:"漫画" summary:"搜索漫画"`
	Keyword  string `json:"keyword" v:"required#搜索关键词不能为空" dc:"搜索关键词"`
	Page     int    `json:"page" d:"1" v:"min:1" dc:"页码"`
	PageSize int    `json:"pageSize" d:"20" v:"min:1|max:100" dc:"每页数量"`
}

// SearchRes 搜索漫画响应
type SearchRes struct {
	List  []ComicInfo `json:"list" dc:"漫画列表"`
	Total int64       `json:"total" dc:"总数量"`
}