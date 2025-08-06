// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"hotgo/api/admin/auth"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	Auth = cAuth{}
)

type cAuth struct{}

// Enable2FA 启用双因子认证
func (c *cAuth) Enable2FA(ctx context.Context, req *auth.Enable2FAReq) (res *auth.Enable2FARes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 验证当前密码
	var mb entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(memberId).Scan(&mb); err != nil {
		return nil, gerror.Wrap(err, "获取用户信息失败，请稍后重试！")
	}

	if gmd5.MustEncryptString(req.Password+mb.Salt) != mb.PasswordHash {
		return nil, gerror.New("密码不正确")
	}

	// 创建输入参数
	input := &adminin.TwoFactorEnableInp{
		UserId: memberId,
	}
	result, err := service.AdminTwoFactor().Enable(ctx, input)
	if err != nil {
		return nil, err
	}

	res = &auth.Enable2FARes{
		TwoFactorEnableModel: result,
	}
	return
}

// Verify2FASetup 验证双因子认证设置
func (c *cAuth) Verify2FASetup(ctx context.Context, req *auth.Verify2FASetupReq) (res *auth.Verify2FASetupRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 创建输入参数
	input := &adminin.TwoFactorConfirmEnableInp{
		UserId: memberId,
		Code:   req.Code,
	}
	result, err := service.AdminTwoFactor().ConfirmEnable(ctx, input)
	if err != nil {
		return nil, err
	}

	res = &auth.Verify2FASetupRes{
		TwoFactorConfirmEnableModel: result,
	}
	return
}

// Disable2FA 禁用双因子认证
func (c *cAuth) Disable2FA(ctx context.Context, req *auth.Disable2FAReq) (res *auth.Disable2FARes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 创建输入参数
	input := &adminin.TwoFactorDisableInp{
		UserId: memberId,
		Code:   req.Code,
	}
	err = service.AdminTwoFactor().Disable(ctx, input)
	return
}

// Get2FAStatus 获取双因子认证状态
func (c *cAuth) Get2FAStatus(ctx context.Context, req *auth.Get2FAStatusReq) (res *auth.Get2FAStatusRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 创建输入参数
	input := &adminin.TwoFactorGetStatusInp{
		UserId: memberId,
	}
	result, err := service.AdminTwoFactor().GetStatus(ctx, input)
	if err != nil {
		return nil, err
	}

	res = &auth.Get2FAStatusRes{
		TwoFactorGetStatusModel: result,
	}
	return
}

// RegenerateBackupCodes 重新生成备用恢复码
func (c *cAuth) RegenerateBackupCodes(ctx context.Context, req *auth.RegenerateBackupCodesReq) (res *auth.RegenerateBackupCodesRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	result, err := service.AdminTwoFactor().RegenerateBackupCodes(ctx, &req.TwoFactorRegenerateBackupCodesInp)
	if err != nil {
		return nil, err
	}

	res = &auth.RegenerateBackupCodesRes{
		TwoFactorRegenerateBackupCodesModel: result,
	}
	return
}

// Verify2FA 验证双因子认证
func (c *cAuth) Verify2FA(ctx context.Context, req *auth.Verify2FAReq) (res *auth.Verify2FARes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	err = service.AdminTwoFactor().Verify(ctx, &req.TwoFactorVerifyInp)
	if err != nil {
		return nil, err
	}

	res = &auth.Verify2FARes{}
	return
}

// verifyTempToken 验证临时token
func (c *cAuth) verifyTempToken(tempToken string) (*model.Identity, error) {
	// 创建一个临时请求来解析token
	req := &ghttp.Request{}
	req.Header.Set("Authorization", "Bearer "+tempToken)

	// 解析临时token
	claims, err := token.ParseLoginUser(req)
	if err != nil {
		return nil, gerror.New("临时token无效或已过期")
	}

	// 检查是否为临时2FA token
	if claims.RoleKey != "temp_2fa" {
		return nil, gerror.New("无效的临时token")
	}

	return claims, nil
}

// VerifyLogin2FA 验证登录双因子认证
// VerifyLogin2FA 验证登录2FA
func (c *cAuth) VerifyLogin2FA(ctx context.Context, req *auth.VerifyLogin2FAReq) (res *auth.VerifyLogin2FARes, err error) {
	// 验证临时token
	userInfo, err := c.verifyTempToken(req.TempToken)
	if err != nil {
		// 记录2FA验证失败日志 - 临时token验证失败
		service.SysLoginLog().Push(ctx, &sysin.LoginLogPushInp{
			Response:         &adminin.LoginModel{Username: "unknown"},
			Err:              gerror.New("2FA验证失败: 临时token无效"),
			TwoFactorEnabled: true,
			TwoFactorMethod:  req.CodeType,
			TwoFactorSuccess: false,
			LoginStep:        "2fa_verification_failed",
		})
		return
	}

	// 验证2FA代码
	verifyInp := &adminin.TwoFactorVerifyInp{
		UserId: userInfo.Id,
		Code:   req.Code, // AdminTwoFactor.Verify方法会根据代码长度自动判断是TOTP还是备用码
	}

	err = service.AdminTwoFactor().Verify(ctx, verifyInp)
	if err != nil {
		// 记录2FA验证失败日志 - 验证码错误
		service.SysLoginLog().Push(ctx, &sysin.LoginLogPushInp{
			Response:         &adminin.LoginModel{Id: userInfo.Id, Username: userInfo.Username},
			Err:              gerror.Newf("2FA验证失败: %s", err.Error()),
			TwoFactorEnabled: true,
			TwoFactorMethod:  req.CodeType,
			TwoFactorSuccess: false,
			LoginStep:        "2fa_verification_failed",
		})
		return
	}

	// 2FA验证成功，完成登录流程
	var member *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(userInfo.Id).Scan(&member); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if member == nil {
		err = gerror.New("用户不存在")
		return
	}

	// 调用AdminSite服务的HandleLogin方法
	loginData, err := service.AdminSite().HandleLogin(ctx, member)
	if err != nil {
		return
	}

	// 记录2FA验证成功日志
	service.SysLoginLog().Push(ctx, &sysin.LoginLogPushInp{
		Response:         loginData,
		Err:              nil,
		TwoFactorEnabled: true,
		TwoFactorMethod:  req.CodeType,
		TwoFactorSuccess: true,
		LoginStep:        "login_complete",
	})

	res = &auth.VerifyLogin2FARes{LoginModel: loginData}
	return
}
