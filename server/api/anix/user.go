// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package anix

import (
	"context"

	"hotgo/api/anix/user/v1"
)

type IUserV1 interface {
	// 用户注册
	Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error)
	// 用户登录
	Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error)
	// 获取用户信息
	GetProfile(ctx context.Context, req *v1.GetProfileReq) (res *v1.GetProfileRes, err error)
	// 更新用户信息
	UpdateProfile(ctx context.Context, req *v1.UpdateProfileReq) (res *v1.UpdateProfileRes, err error)
	// 修改密码
	ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (res *v1.ChangePasswordRes, err error)
	// 用户登出
	Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error)
}