// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixAuthor is the golang structure of table anix_authors for DAO operations like Where/Data.
type AnixAuthor struct {
	g.Meta      `orm:"table:anix_authors, do:true"`
	Id          interface{} // 作者ID
	Name        interface{} // 作者名称
	Avatar      interface{} // 作者头像
	Description interface{} // 作者简介
	Status      interface{} // 状态:1=正常,2=禁用
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}