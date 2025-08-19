// Package user
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package user

import (
	"context"
	"hotgo/api/anix/user/v1"
	"hotgo/internal/model/input/anixin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// UserV1 用户控制器v1
type UserV1 struct{}

// NewV1 创建用户控制器v1实例
func NewV1() *UserV1 {
	return &UserV1{}
}

// Register 用户注册
func (c *UserV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	inp := &anixin.RegisterInp{
		Username:        req.Username,
		Email:          req.Email,
		Nickname:       req.Nickname,
		Password:       req.Password,
		ConfirmPassword: req.ConfirmPassword,
		AgreeTerms:     req.AgreeTerms,
	}

	err = service.AnixUser().Register(ctx, inp)
	if err != nil {
		return nil, err
	}

	res = &v1.RegisterRes{}
	return res, nil
}

// Login 用户登录
func (c *UserV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	// 获取客户端IP
	r := g.RequestFromCtx(ctx)
	clientIP := r.GetClientIp()

	inp := &anixin.LoginInp{
		Account:    req.Account,
		Password:   req.Password,
		RememberMe: req.RememberMe,
		Ip:         clientIP,
	}

	loginModel, err := service.AnixUser().Login(ctx, inp)
	if err != nil {
		return nil, err
	}

	// 转换用户信息
	userInfo := &v1.UserInfo{
		Id:          int64(loginModel.User.Id),
		Username:    loginModel.User.Username,
		Nickname:    loginModel.User.Nickname,
		Avatar:      loginModel.User.Avatar,
		Email:       loginModel.User.Email,
		Gender:      loginModel.User.Gender,
		Birthday:    gconv.String(loginModel.User.Birthday),
		Mobile:      loginModel.User.Phone, // 实体模型中使用phone字段
		Status:      loginModel.User.Status,
		LastLoginAt: gconv.String(loginModel.User.LastLoginAt),
		LastLoginIp: loginModel.User.LastLoginIp,
		CreatedAt:   gconv.String(loginModel.User.CreatedAt),
		UpdatedAt:   gconv.String(loginModel.User.UpdatedAt),
	}

	res = &v1.LoginRes{
		Token: loginModel.Token,
		User:  userInfo,
	}

	return res, nil
}

// Profile 获取用户信息
func (c *UserV1) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	// 从请求中获取用户信息
	r := g.RequestFromCtx(ctx)
	user, err := service.AnixUser().GetUserFromRequest(r)
	if err != nil {
		return nil, err
	}
// 获取用户信息
	userInfo, err := service.AnixUser().GetProfile(ctx, int64(user.Id))
	if err != nil {
		return nil, err
	}

	res = &v1.ProfileRes{
		User: userInfo,
	}

	return res, nil
}

// UpdateProfile 更新用户信息
func (c *UserV1) UpdateProfile(ctx context.Context, req *v1.UpdateProfileReq) (res *v1.UpdateProfileRes, err error) {
	// 从请求中获取用户信息
	r := g.RequestFromCtx(ctx)
	user, err := service.AnixUser().GetUserFromRequest(r)
	if err != nil {
		return nil, err
	}

	// 转换生日字符串为时间
	var birthday *gtime.Time
	if req.Birthday != "" {
		birthday = gtime.NewFromStr(req.Birthday)
	}

	inp := &anixin.UpdateProfileInp{
		UserId:   int64(user.Id),
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Birthday: birthday,
		Mobile:   req.Mobile,
	}

	err = service.AnixUser().UpdateProfile(ctx, inp)
	if err != nil {
		return nil, err
	}

	res = &v1.UpdateProfileRes{}
	return res, nil
}

// ChangePassword 修改密码
func (c *UserV1) ChangePassword(ctx context.Context, req *v1.ChangePasswordReq) (res *v1.ChangePasswordRes, err error) {
	// 从请求中获取用户信息
	r := g.RequestFromCtx(ctx)
	user, err := service.AnixUser().GetUserFromRequest(r)
	if err != nil {
		return nil, err
	}

	inp := &anixin.ChangePasswordInp{
		UserId:          int64(user.Id),
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	err = service.AnixUser().ChangePassword(ctx, inp)
	if err != nil {
		return nil, err
	}

	res = &v1.ChangePasswordRes{}
	return res, nil
}

// Logout 用户登出
func (c *UserV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	// 从Header获取Token
	token := g.RequestFromCtx(ctx).Header.Get("Authorization")
	if token != "" && gstr.HasPrefix(token, "Bearer ") {
		token = gstr.SubStr(token, 7)
	}
	
	// 调用服务层登出
	err = service.AnixUser().Logout(ctx, token)
	if err != nil {
		return nil, err
	}
	
	return &v1.LogoutRes{}, nil
}