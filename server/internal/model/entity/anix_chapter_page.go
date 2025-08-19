// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixChapterPage is the golang structure for table anix_chapter_pages.
type AnixChapterPage struct {
	Id        uint64      `json:"id"        orm:"id"        description:"页面ID"`
	ChapterId uint64      `json:"chapterId" orm:"chapter_id" description:"章节ID"`
	PageNum   int         `json:"pageNum"   orm:"page_num"   description:"页面序号"`
	ImageUrl  string      `json:"imageUrl"  orm:"image_url"  description:"图片URL"`
	Width     int         `json:"width"     orm:"width"     description:"图片宽度"`
	Height    int         `json:"height"    orm:"height"    description:"图片高度"`
	Sort      int         `json:"sort"      orm:"sort"      description:"排序"`
	Status    int         `json:"status"    orm:"status"    description:"状态:1=正常,2=禁用"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`
}