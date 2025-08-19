// Package v1
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// UserInfo 用户信息
type UserInfo struct {
	Id          int64  `json:"id" dc:"用户ID"`
	Username    string `json:"username" dc:"用户名"`
	Nickname    string `json:"nickname" dc:"昵称"`
	Avatar      string `json:"avatar" dc:"头像"`
	Email       string `json:"email" dc:"邮箱"`
	Gender      int    `json:"gender" dc:"性别:0=未知,1=男,2=女"`
	Birthday    string `json:"birthday" dc:"生日"`
	Mobile      string `json:"mobile" dc:"手机号"`
	Status      int    `json:"status" dc:"状态:1=正常,2=禁用"`
	LastLoginAt string `json:"lastLoginAt" dc:"最后登录时间"`
	LastLoginIp string `json:"lastLoginIp" dc:"最后登录IP"`
	CreatedAt   string `json:"createdAt" dc:"创建时间"`
	UpdatedAt   string `json:"updatedAt" dc:"更新时间"`
}

// RegisterReq 用户注册请求
type RegisterReq struct {
	g.Meta          `path:"/register" method:"post" tags:"用户" summary:"用户注册"`
	Username        string `json:"username" v:"required|length:3,20|regex:^[a-zA-Z0-9_]+$#用户名不能为空|用户名长度为3-20个字符|用户名只能包含字母、数字和下划线" dc:"用户名"`
	Password        string `json:"password" v:"required|length:6,50#密码不能为空|密码长度为6-50个字符" dc:"密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:Password#确认密码不能为空|两次输入的密码不一致" dc:"确认密码"`
	Email           string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确" dc:"邮箱"`
	Nickname        string `json:"nickname" v:"length:0,50#昵称不能超过50个字符" dc:"昵称"`
	AgreeTerms      bool   `json:"agreeTerms" v:"required#请同意用户协议" dc:"是否同意用户协议"`
}

// RegisterRes 用户注册响应
type RegisterRes struct {
	Message string `json:"message" dc:"注册结果消息"`
}

// LoginReq 用户登录请求
type LoginReq struct {
	g.Meta     `path:"/login" method:"post" tags:"用户" summary:"用户登录"`
	Account    string `json:"account" v:"required#账号不能为空" dc:"账号(用户名或邮箱)"`
	Password   string `json:"password" v:"required#密码不能为空" dc:"密码"`
	RememberMe bool   `json:"rememberMe" dc:"记住登录状态"`
}

// LoginRes 用户登录响应
type LoginRes struct {
	Token string    `json:"token" dc:"访问令牌"`
	User  *UserInfo `json:"user" dc:"用户信息"`
}

// ProfileReq 获取用户信息请求
type ProfileReq struct {
	g.Meta `path:"/info" method:"get" tags:"用户" summary:"获取用户信息"`
}

// ProfileRes 获取用户信息响应
type ProfileRes struct {
	User *UserInfo `json:"user" dc:"用户信息"`
}

// UpdateProfileReq 更新用户信息请求
type UpdateProfileReq struct {
	g.Meta   `path:"/profile" method:"put" tags:"用户" summary:"更新用户信息"`
	Nickname string `json:"nickname" v:"length:0,50#昵称不能超过50个字符" dc:"昵称"`
	Avatar   string `json:"avatar" dc:"头像"`
	Gender   int    `json:"gender" v:"in:0,1,2#性别值不正确" dc:"性别:0=未知,1=男,2=女"`
	Birthday string `json:"birthday" dc:"生日"`
	Mobile   string `json:"mobile" v:"phone#手机号格式不正确" dc:"手机号"`
}

// UpdateProfileRes 更新用户信息响应
type UpdateProfileRes struct {
	Message string `json:"message" dc:"响应消息"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	g.Meta          `path:"/password" method:"put" tags:"用户" summary:"修改密码"`
	OldPassword     string `json:"oldPassword" v:"required#原密码不能为空" dc:"原密码"`
	NewPassword     string `json:"newPassword" v:"required|length:6,50#新密码不能为空|新密码长度为6-50个字符" dc:"新密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:NewPassword#确认密码不能为空|两次输入的密码不一致" dc:"确认密码"`
}

// ChangePasswordRes 修改密码响应
type ChangePasswordRes struct {
	Message string `json:"message" dc:"响应消息"`
}

// LogoutReq 用户登出请求
type LogoutReq struct {
	g.Meta `path:"/logout" method:"post" tags:"用户" summary:"用户登出"`
}

// LogoutRes 用户登出响应
type LogoutRes struct {
	Message string `json:"message" dc:"响应消息"`
}
