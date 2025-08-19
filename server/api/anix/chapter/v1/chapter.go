// Package v1
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ChapterInfo 章节信息
type ChapterInfo struct {
	Id          int64  `json:"id" dc:"章节ID"`
	ComicId     int64  `json:"comic_id" dc:"漫画ID"`
	Title       string `json:"title" dc:"章节标题"`
	ChapterNum  int    `json:"chapter_num" dc:"章节序号"`
	PageCount   int    `json:"page_count" dc:"页面数量"`
	ViewCount   int64  `json:"view_count" dc:"浏览次数"`
	IsFree      bool   `json:"is_free" dc:"是否免费"`
	PublishTime string `json:"publish_time" dc:"发布时间"`
}

// PageInfo 页面信息
type PageInfo struct {
	Id       int64  `json:"id" dc:"页面ID"`
	PageNum  int    `json:"page_num" dc:"页面序号"`
	ImageUrl string `json:"image_url" dc:"图片地址"`
	Width    int    `json:"width" dc:"图片宽度"`
	Height   int    `json:"height" dc:"图片高度"`
}

// ListReq 获取章节列表请求
type ListReq struct {
	g.Meta  `path:"/chapter/list" method:"get" tags:"章节" summary:"获取章节列表"`
	ComicId int64 `json:"comicId" v:"required#漫画ID不能为空" dc:"漫画ID"`
	Page    int   `json:"page" d:"1" v:"min:1" dc:"页码"`
	Size    int   `json:"size" d:"20" v:"min:1|max:100" dc:"每页数量"`
}

// ListRes 获取章节列表响应
type ListRes struct {
	List  []ChapterInfo `json:"list" dc:"章节列表"`
	Total int64         `json:"total" dc:"总数量"`
}





// GetDetailReq 获取章节详情请求
type GetDetailReq struct {
	g.Meta `path:"/chapter/detail/:id" method:"get" tags:"章节" summary:"获取章节详情"`
	Id     int64 `json:"id" v:"required#章节ID不能为空" dc:"章节ID"`
}

// GetDetailRes 获取章节详情响应
type GetDetailRes struct {
	Chapter     ChapterInfo `json:"chapter" dc:"章节信息"`
	PrevChapter *ChapterInfo `json:"prev_chapter" dc:"上一章节"`
	NextChapter *ChapterInfo `json:"next_chapter" dc:"下一章节"`
}

// GetPagesReq 获取章节页面请求
type GetPagesReq struct {
	g.Meta    `path:"/chapter/pages/:id" method:"get" tags:"章节" summary:"获取章节页面"`
	ChapterId int64 `json:"chapterId" v:"required#章节ID不能为空" dc:"章节ID"`
}

// GetPagesRes 获取章节页面响应
type GetPagesRes struct {
	Pages []PageInfo `json:"pages" dc:"页面列表"`
	Total int64      `json:"total" dc:"总数量"`
}

// RecordReadingReq 记录阅读历史请求
type RecordReadingReq struct {
	g.Meta      `path:"/chapter/reading" method:"post" tags:"章节" summary:"记录阅读历史"`
	ChapterId   int64 `json:"chapter_id" v:"required#章节ID不能为空" dc:"章节ID"`
	PageIndex   int   `json:"page_index" v:"min:0" dc:"页面索引"`
	ReadingTime int64 `json:"reading_time" dc:"阅读时长(秒)"`
	Progress    float64 `json:"progress" v:"min:0|max:100" dc:"阅读进度百分比"`
}

// RecordReadingRes 记录阅读历史响应
type RecordReadingRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// GetPreviousReq 获取上一章节请求
type GetPreviousReq struct {
	g.Meta    `path:"/chapter/previous/:id" method:"get" tags:"章节" summary:"获取上一章节"`
	ChapterId int64 `json:"chapterId" v:"required#章节ID不能为空" dc:"章节ID"`
}

// GetPreviousRes 获取上一章节响应
type GetPreviousRes struct {
	Chapter *ChapterInfo `json:"chapter" dc:"上一章节信息"`
}

// GetNextReq 获取下一章节请求
type GetNextReq struct {
	g.Meta    `path:"/chapter/next/:id" method:"get" tags:"章节" summary:"获取下一章节"`
	ChapterId int64 `json:"chapterId" v:"required#章节ID不能为空" dc:"章节ID"`
}

// GetNextRes 获取下一章节响应
type GetNextRes struct {
	Chapter *ChapterInfo `json:"chapter" dc:"下一章节信息"`
}