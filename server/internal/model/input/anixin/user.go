// Package anixin AniX用户输入模型
package anixin

import (
	"hotgo/internal/model/entity"

	"github.com/gogf/gf/v2/os/gtime"
)

// RegisterInp 用户注册输入
type RegisterInp struct {
	Username        string `json:"username" v:"required|length:3,20#用户名不能为空|用户名长度为3-20位" dc:"用户名"`
	Email          string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确" dc:"邮箱"`
	Nickname       string `json:"nickname" v:"required|length:2,20#昵称不能为空|昵称长度为2-20位" dc:"昵称"`
	Password       string `json:"password" v:"required|length:6,20#密码不能为空|密码长度为6-20位" dc:"密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:password#确认密码不能为空|两次密码输入不一致" dc:"确认密码"`
	AgreeTerms     bool   `json:"agreeTerms" v:"required|eq:true#请同意用户协议和隐私政策" dc:"同意条款"`
}

// LoginInp 用户登录输入
type LoginInp struct {
	Account    string `json:"account" v:"required#账号不能为空" dc:"账号(用户名或邮箱)"`
	Password   string `json:"password" v:"required#密码不能为空" dc:"密码"`
	RememberMe bool   `json:"rememberMe" dc:"记住登录"`
	Ip         string `json:"-" dc:"登录IP"`
}

// LoginModel 登录返回模型
type LoginModel struct {
	User  *entity.AnixUser `json:"user" dc:"用户信息"`
	Token string           `json:"token" dc:"访问令牌"`
}

// UpdateProfileInp 更新用户信息输入
type UpdateProfileInp struct {
	UserId   int64       `json:"-" dc:"用户ID"`
	Nickname string      `json:"nickname" v:"length:2,20#昵称长度为2-20位" dc:"昵称"`
	Gender   int         `json:"gender" v:"in:0,1,2#性别值不正确" dc:"性别:0=未知,1=男,2=女"`
	Birthday *gtime.Time `json:"birthday" dc:"生日"`
	Mobile   string      `json:"mobile" v:"phone#手机号格式不正确" dc:"手机号"`
	Avatar   string      `json:"avatar" dc:"头像"`
}

// ChangePasswordInp 修改密码输入
type ChangePasswordInp struct {
	UserId          int64  `json:"-" dc:"用户ID"`
	OldPassword     string `json:"oldPassword" v:"required#原密码不能为空" dc:"原密码"`
	NewPassword     string `json:"newPassword" v:"required|length:6,20#新密码不能为空|新密码长度为6-20位" dc:"新密码"`
	ConfirmPassword string `json:"confirmPassword" v:"required|same:newPassword#确认密码不能为空|两次密码输入不一致" dc:"确认密码"`
}

// GetProfileInp 获取用户信息输入
type GetProfileInp struct {
	UserId int64 `json:"-" dc:"用户ID"`
}

// LogoutInp 用户登出输入
type LogoutInp struct {
	Token string `json:"-" dc:"访问令牌"`
}

// VerifyTokenInp 验证Token输入
type VerifyTokenInp struct {
	Token string `json:"token" v:"required#Token不能为空" dc:"访问令牌"`
}

// GetUserInp 获取用户输入
type GetUserInp struct {
	UserId   int64  `json:"userId" dc:"用户ID"`
	Username string `json:"username" dc:"用户名"`
	Email    string `json:"email" dc:"邮箱"`
}

// UpdateLastLoginInp 更新最后登录信息输入
type UpdateLastLoginInp struct {
	UserId int64  `json:"userId" dc:"用户ID"`
	Ip     string `json:"ip" dc:"登录IP"`
}