// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixAuthor is the golang structure for table anix_authors.
type AnixAuthor struct {
	Id          uint64      `json:"id"          orm:"id"          description:"作者ID"`
	Name        string      `json:"name"        orm:"name"        description:"作者名称"`
	Avatar      string      `json:"avatar"      orm:"avatar"      description:"作者头像"`
	Description string      `json:"description" orm:"description" description:"作者简介"`
	Status      int         `json:"status"      orm:"status"      description:"状态:1=正常,2=禁用"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`
}