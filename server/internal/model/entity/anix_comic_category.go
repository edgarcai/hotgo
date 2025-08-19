// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixComicCategory is the golang structure for table anix_comic_categories.
type AnixComicCategory struct {
	Id         uint64      `json:"id"         orm:"id"         description:"关联ID"`
	ComicId    uint64      `json:"comicId"    orm:"comic_id"    description:"漫画ID"`
	CategoryId uint64      `json:"categoryId" orm:"category_id" description:"分类ID"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"删除时间"`
}