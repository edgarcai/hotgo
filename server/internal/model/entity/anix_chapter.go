// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixChapter is the golang structure for table anix_chapters.
type AnixChapter struct {
	Id          uint64      `json:"id"          orm:"id"          description:"章节ID"`
	ComicId     uint64      `json:"comicId"     orm:"comic_id"     description:"漫画ID"`
	Title       string      `json:"title"       orm:"title"       description:"章节标题"`
	ChapterNum  int         `json:"chapterNum"  orm:"chapter_num"  description:"章节序号"`
	IsVip       int         `json:"isVip"       orm:"is_vip"       description:"是否VIP:0=否,1=是"`
	ViewCount   uint64      `json:"viewCount"   orm:"view_count"   description:"浏览次数"`
	Sort        int         `json:"sort"        orm:"sort"        description:"排序"`
	Status      int         `json:"status"      orm:"status"      description:"状态:1=正常,2=禁用"`
	PublishedAt *gtime.Time `json:"publishedAt" orm:"published_at" description:"发布时间"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`
}