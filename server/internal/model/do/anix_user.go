// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixUser is the golang structure of table anix_users for DAO operations like Where/Data.
type AnixUser struct {
	g.Meta      `orm:"table:anix_users, do:true"`
	Id          interface{} // 用户ID
	Username    interface{} // 用户名
	Email       interface{} // 邮箱
	Password    interface{} // 密码
	Nickname    interface{} // 昵称
	Avatar      interface{} // 头像
	Gender      interface{} // 性别:0=未知,1=男,2=女
	Birthday    *gtime.Time // 生日
	Phone       interface{} // 手机号
	Status      interface{} // 状态:1=正常,2=禁用
	LastLoginAt *gtime.Time // 最后登录时间
	LastLoginIp interface{} // 最后登录IP
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}