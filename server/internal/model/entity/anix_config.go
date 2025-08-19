// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixConfig is the golang structure for table anix_configs.
type AnixConfig struct {
	Id          uint64      `json:"id"          orm:"id"          description:"配置ID"`
	Key         string      `json:"key"         orm:"key"         description:"配置键"`
	Value       string      `json:"value"       orm:"value"       description:"配置值"`
	Description string      `json:"description" orm:"description" description:"配置描述"`
	Type        string      `json:"type"        orm:"type"        description:"配置类型"`
	Group       string      `json:"group"       orm:"group"       description:"配置分组"`
	Sort        int         `json:"sort"        orm:"sort"        description:"排序"`
	Status      int         `json:"status"      orm:"status"      description:"状态:1=正常,2=禁用"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`
}