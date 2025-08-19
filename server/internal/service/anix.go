// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hotgo/api/anix/user/v1"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/anixin"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	// IAnixUser AniX用户服务接口
	IAnixUser interface {
		// Model 用户ORM模型
		Model(ctx context.Context) *gdb.Model
		// Register 用户注册
		Register(ctx context.Context, in *anixin.RegisterInp) (err error)
		// Login 用户登录
		Login(ctx context.Context, in *anixin.LoginInp) (res *anixin.LoginModel, err error)
		// GetProfile 获取用户信息
		GetProfile(ctx context.Context, userId int64) (res *v1.UserInfo, err error)
		// UpdateProfile 更新用户信息
		UpdateProfile(ctx context.Context, in *anixin.UpdateProfileInp) (err error)
		// ChangePassword 修改密码
		ChangePassword(ctx context.Context, in *anixin.ChangePasswordInp) (err error)
		// Logout 用户登出
		Logout(ctx context.Context, token string) (err error)
		// VerifyToken 验证Token
		VerifyToken(ctx context.Context, token string) (userId int64, err error)
		// GetUserByUsername 根据用户名获取用户
		GetUserByUsername(ctx context.Context, username string) (user *entity.AnixUser, err error)
		// GetUserByEmail 根据邮箱获取用户
		GetUserByEmail(ctx context.Context, email string) (user *entity.AnixUser, err error)
		// IsEmailExists 检查邮箱是否存在
		IsEmailExists(ctx context.Context, email string) (exists bool, err error)
		// IsUsernameExists 检查用户名是否存在
		IsUsernameExists(ctx context.Context, username string) (exists bool, err error)
		// GenerateToken 生成JWT Token
		GenerateToken(ctx context.Context, user *entity.AnixUser) (token string, err error)
		// HashPassword 密码加密
		HashPassword(password string) (hashedPassword string, err error)
		// VerifyPassword 验证密码
		VerifyPassword(password, hashedPassword string) bool
		// UpdateLastLogin 更新最后登录信息
		UpdateLastLogin(ctx context.Context, userId int64, ip string) (err error)
		// GetUserFromContext 从上下文获取用户信息
		GetUserFromContext(ctx context.Context) (user *entity.AnixUser, err error)
		// GetUserFromRequest 从请求获取用户信息
		GetUserFromRequest(r *ghttp.Request) (user *entity.AnixUser, err error)
	}
)

var (
	localAnixUser IAnixUser
)

// AnixUser 获取AniX用户服务实例
func AnixUser() IAnixUser {
	if localAnixUser == nil {
		panic("implement not found for interface IAnixUser, forgot register?")
	}
	return localAnixUser
}

// RegisterAnixUser 注册AniX用户服务实例
func RegisterAnixUser(i IAnixUser) {
	localAnixUser = i
}