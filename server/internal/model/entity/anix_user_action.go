// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixUserAction is the golang structure for table anix_user_actions.
type AnixUserAction struct {
	Id         uint64      `json:"id"         orm:"id"         description:"行为ID"`
	UserId     uint64      `json:"userId"     orm:"user_id"     description:"用户ID"`
	ComicId    uint64      `json:"comicId"    orm:"comic_id"    description:"漫画ID"`
	ActionType int         `json:"actionType" orm:"action_type" description:"行为类型:1=点赞,2=收藏,3=评分"`
	ActionData string      `json:"actionData" orm:"action_data" description:"行为数据(JSON格式)"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:"更新时间"`
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"删除时间"`
}