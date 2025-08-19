package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SaveReadingProgressReq 保存阅读进度请求
type SaveReadingProgressReq struct {
	g.Meta    `path:"/history/progress" method:"post" summary:"保存阅读进度" tags:"阅读历史"`
	ComicId   int64   `json:"comic_id" v:"required" dc:"漫画ID"`
	ChapterId int64   `json:"chapter_id" v:"required" dc:"章节ID"`
	PageIndex int     `json:"page_index" v:"min:0" dc:"当前页面索引"`
	Progress  float64 `json:"progress" v:"between:0,100" dc:"阅读进度百分比"`
	ReadingTime int64 `json:"reading_time" dc:"本次阅读时长(秒)"`
}

type SaveReadingProgressRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// GetReadingHistoryReq 获取阅读历史请求
type GetReadingHistoryReq struct {
	g.Meta `path:"/history/list" method:"get" summary:"获取阅读历史" tags:"阅读历史"`
	Page   int `json:"page" v:"min:1" dc:"页码" d:"1"`
	Limit  int `json:"limit" v:"between:1,100" dc:"每页数量" d:"20"`
}

type ReadingHistoryItem struct {
	Id          int64       `json:"id" dc:"记录ID"`
	ComicId     int64       `json:"comic_id" dc:"漫画ID"`
	ComicTitle  string      `json:"comic_title" dc:"漫画标题"`
	ComicCover  string      `json:"comic_cover" dc:"漫画封面"`
	ChapterId   int64       `json:"chapter_id" dc:"章节ID"`
	ChapterTitle string     `json:"chapter_title" dc:"章节标题"`
	PageIndex   int         `json:"page_index" dc:"当前页面索引"`
	Progress    float64     `json:"progress" dc:"阅读进度百分比"`
	ReadingTime int64       `json:"reading_time" dc:"累计阅读时长(秒)"`
	LastReadAt  *gtime.Time `json:"last_read_at" dc:"最后阅读时间"`
}

type GetReadingHistoryRes struct {
	List  []ReadingHistoryItem `json:"list" dc:"阅读历史列表"`
	Total int                 `json:"total" dc:"总数量"`
	Page  int                 `json:"page" dc:"当前页码"`
	Limit int                 `json:"limit" dc:"每页数量"`
}

// GetReadingProgressReq 获取阅读进度请求
type GetReadingProgressReq struct {
	g.Meta    `path:"/history/progress/:comic_id/:chapter_id" method:"get" summary:"获取阅读进度" tags:"阅读历史"`
	ComicId   int64 `json:"comic_id" v:"required" dc:"漫画ID"`
	ChapterId int64 `json:"chapter_id" v:"required" dc:"章节ID"`
}

type GetReadingProgressRes struct {
	PageIndex   int         `json:"page_index" dc:"当前页面索引"`
	Progress    float64     `json:"progress" dc:"阅读进度百分比"`
	ReadingTime int64       `json:"reading_time" dc:"累计阅读时长(秒)"`
	LastReadAt  *gtime.Time `json:"last_read_at" dc:"最后阅读时间"`
}

// AddBookmarkReq 添加书签请求
type AddBookmarkReq struct {
	g.Meta    `path:"/history/bookmark" method:"post" summary:"添加书签" tags:"阅读历史"`
	ComicId   int64  `json:"comic_id" v:"required" dc:"漫画ID"`
	ChapterId int64  `json:"chapter_id" v:"required" dc:"章节ID"`
	PageIndex int    `json:"page_index" v:"min:0" dc:"书签页面索引"`
	Note      string `json:"note" dc:"书签备注"`
}

type AddBookmarkRes struct {
	Id int64 `json:"id" dc:"书签ID"`
}

// GetBookmarksReq 获取书签列表请求
type GetBookmarksReq struct {
	g.Meta  `path:"/history/bookmarks" method:"get" summary:"获取书签列表" tags:"阅读历史"`
	ComicId int64 `json:"comic_id" dc:"漫画ID，为空则获取所有书签"`
	Page    int   `json:"page" v:"min:1" dc:"页码" d:"1"`
	Limit   int   `json:"limit" v:"between:1,100" dc:"每页数量" d:"20"`
}

type BookmarkItem struct {
	Id           int64       `json:"id" dc:"书签ID"`
	ComicId      int64       `json:"comic_id" dc:"漫画ID"`
	ComicTitle   string      `json:"comic_title" dc:"漫画标题"`
	ComicCover   string      `json:"comic_cover" dc:"漫画封面"`
	ChapterId    int64       `json:"chapter_id" dc:"章节ID"`
	ChapterTitle string      `json:"chapter_title" dc:"章节标题"`
	PageIndex    int         `json:"page_index" dc:"书签页面索引"`
	Note         string      `json:"note" dc:"书签备注"`
	CreatedAt    *gtime.Time `json:"created_at" dc:"创建时间"`
}

type GetBookmarksRes struct {
	List  []BookmarkItem `json:"list" dc:"书签列表"`
	Total int           `json:"total" dc:"总数量"`
	Page  int           `json:"page" dc:"当前页码"`
	Limit int           `json:"limit" dc:"每页数量"`
}

// DeleteBookmarkReq 删除书签请求
type DeleteBookmarkReq struct {
	g.Meta `path:"/history/bookmark/:id" method:"delete" summary:"删除书签" tags:"阅读历史"`
	Id     int64 `json:"id" v:"required" dc:"书签ID"`
}

type DeleteBookmarkRes struct {
	Success bool `json:"success" dc:"是否成功"`
}

// GetReadingStatsReq 获取阅读统计请求
type GetReadingStatsReq struct {
	g.Meta  `path:"/history/stats" method:"get" summary:"获取阅读统计" tags:"阅读历史"`
	ComicId int64 `json:"comic_id" dc:"漫画ID，为空则获取总体统计"`
}

type ReadingStatsItem struct {
	ComicId          int64       `json:"comic_id" dc:"漫画ID"`
	ComicTitle       string      `json:"comic_title" dc:"漫画标题"`
	ComicCover       string      `json:"comic_cover" dc:"漫画封面"`
	TotalChapters    int         `json:"total_chapters" dc:"总章节数"`
	ReadChapters     int         `json:"read_chapters" dc:"已读章节数"`
	TotalReadingTime int64       `json:"total_reading_time" dc:"总阅读时长(秒)"`
	FirstReadAt      *gtime.Time `json:"first_read_at" dc:"首次阅读时间"`
	LastReadAt       *gtime.Time `json:"last_read_at" dc:"最后阅读时间"`
}

type GetReadingStatsRes struct {
	Stats []ReadingStatsItem `json:"stats" dc:"阅读统计列表"`
	TotalReadingTime int64  `json:"total_reading_time" dc:"总阅读时长(秒)"`
	TotalComics      int    `json:"total_comics" dc:"总漫画数"`
}