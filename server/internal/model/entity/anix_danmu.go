// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixDanmu is the golang structure for table anix_danmu.
type AnixDanmu struct {
	Id        uint64      `json:"id"        orm:"id"        description:"弹幕ID"`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID"`
	ComicId   uint64      `json:"comicId"   orm:"comic_id"   description:"漫画ID"`
	ChapterId uint64      `json:"chapterId" orm:"chapter_id" description:"章节ID"`
	PageNum   int         `json:"pageNum"   orm:"page_num"   description:"页面序号"`
	Content   string      `json:"content"   orm:"content"   description:"弹幕内容"`
	Color     string      `json:"color"     orm:"color"     description:"弹幕颜色"`
	PositionX float64     `json:"positionX" orm:"position_x" description:"X坐标位置"`
	PositionY float64     `json:"positionY" orm:"position_y" description:"Y坐标位置"`
	Status    int         `json:"status"    orm:"status"    description:"状态:1=正常,2=隐藏,3=删除"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`
}