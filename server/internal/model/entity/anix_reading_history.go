// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixReadingHistory is the golang structure for table anix_reading_history.
type AnixReadingHistory struct {
	Id           uint64      `json:"id"           orm:"id"           description:"历史ID"`
	UserId       uint64      `json:"userId"       orm:"user_id"       description:"用户ID"`
	ComicId      uint64      `json:"comicId"      orm:"comic_id"      description:"漫画ID"`
	ChapterId    uint64      `json:"chapterId"    orm:"chapter_id"    description:"章节ID"`
	PageNum      int         `json:"pageNum"      orm:"page_num"      description:"阅读到的页面"`
	ReadProgress float64     `json:"readProgress" orm:"read_progress" description:"阅读进度(0-100)"`
	ReadTime     int         `json:"readTime"     orm:"read_time"     description:"阅读时长(秒)"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`
}