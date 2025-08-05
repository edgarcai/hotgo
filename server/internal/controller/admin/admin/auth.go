// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/admin/auth"
	"hotgo/internal/library/contexts"
	"hotgo/internal/service"
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

	req.UserId = memberId
	result, err := service.AdminTwoFactor().Enable(ctx, &req.TwoFactorEnableInp)
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

	req.UserId = memberId
	err = service.AdminTwoFactor().ConfirmEnable(ctx, &req.TwoFactorConfirmEnableInp)
	return
}

// Disable2FA 禁用双因子认证
func (c *cAuth) Disable2FA(ctx context.Context, req *auth.Disable2FAReq) (res *auth.Disable2FARes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	err = service.AdminTwoFactor().Disable(ctx, &req.TwoFactorDisableInp)
	return
}

// Get2FAStatus 获取双因子认证状态
func (c *cAuth) Get2FAStatus(ctx context.Context, req *auth.Get2FAStatusReq) (res *auth.Get2FAStatusRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	result, err := service.AdminTwoFactor().GetStatus(ctx, &req.TwoFactorGetStatusInp)
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