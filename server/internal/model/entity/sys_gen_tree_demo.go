// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysGenTreeDemo is the golang structure for table sys_gen_tree_demo.
type SysGenTreeDemo struct {
	Id          int64       `json:"id"          description:"ID"`
	Pid         int64       `json:"pid"         description:"上级ID"`
	Level       int         `json:"level"       description:"关系树级别"`
	Tree        string      `json:"tree"        description:"关系树"`
	CategoryId  int64       `json:"categoryId"  description:"分类ID"`
	Title       string      `json:"title"       description:"标题"`
	Description string      `json:"description" description:"描述"`
	Sort        int         `json:"sort"        description:"排序"`
	Status      int         `json:"status"      description:"状态"`
	CreatedBy   int64       `json:"createdBy"   description:"创建者"`
	UpdatedBy   int64       `json:"updatedBy"   description:"更新者"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"修改时间"`
	DeletedAt   *gtime.Time `json:"deletedAt"   description:"删除时间"`
}
