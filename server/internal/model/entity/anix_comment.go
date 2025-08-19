// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixComment is the golang structure for table anix_comments.
type AnixComment struct {
	Id        uint64      `json:"id"        orm:"id"        description:"评论ID"`
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID"`
	ComicId   uint64      `json:"comicId"   orm:"comic_id"   description:"漫画ID"`
	ChapterId uint64      `json:"chapterId" orm:"chapter_id" description:"章节ID"`
	ParentId  uint64      `json:"parentId"  orm:"parent_id"  description:"父评论ID"`
	Content   string      `json:"content"   orm:"content"   description:"评论内容"`
	LikeCount uint64      `json:"likeCount" orm:"like_count" description:"点赞数"`
	Status    int         `json:"status"    orm:"status"    description:"状态:1=正常,2=隐藏,3=删除"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`
}