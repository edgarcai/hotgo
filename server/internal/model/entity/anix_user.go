// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixUser is the golang structure for table anix_users.
type AnixUser struct {
	Id           uint64      `json:"id"           orm:"id"           description:"用户ID"`
	Username     string      `json:"username"     orm:"username"     description:"用户名"`
	Email        string      `json:"email"        orm:"email"        description:"邮箱"`
	Password     string      `json:"password"     orm:"password"     description:"密码"`
	Nickname     string      `json:"nickname"     orm:"nickname"     description:"昵称"`
	Avatar       string      `json:"avatar"       orm:"avatar"       description:"头像"`
	Gender       int         `json:"gender"       orm:"gender"       description:"性别:0=未知,1=男,2=女"`
	Birthday     *gtime.Time `json:"birthday"     orm:"birthday"     description:"生日"`
	Phone        string      `json:"phone"        orm:"phone"        description:"手机号"`
	Status       int         `json:"status"       orm:"status"       description:"状态:1=正常,2=禁用"`
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" description:"最后登录时间"`
	LastLoginIp  string      `json:"lastLoginIp"  orm:"last_login_ip" description:"最后登录IP"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`
}