// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/api/admin/twofactor"
	"hotgo/internal/library/contexts"
	"hotgo/internal/service"
)

var (
	TwoFactor = cTwoFactor{}
)

type cTwoFactor struct{}

// Enable 启用双因子认证
func (c *cTwoFactor) Enable(ctx context.Context, req *twofactor.EnableReq) (res *twofactor.EnableRes, err error) {
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

	res = &twofactor.EnableRes{
		TwoFactorEnableModel: result,
	}
	return
}

// ConfirmEnable 确认启用双因子认证
func (c *cTwoFactor) ConfirmEnable(ctx context.Context, req *twofactor.ConfirmEnableReq) (res *twofactor.ConfirmEnableRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	err = service.AdminTwoFactor().ConfirmEnable(ctx, &req.TwoFactorConfirmEnableInp)
	return
}

// Disable 禁用双因子认证
func (c *cTwoFactor) Disable(ctx context.Context, req *twofactor.DisableReq) (res *twofactor.DisableRes, err error) {
	// 获取当前用户ID
	var memberId = contexts.Get(ctx).User.Id
	if memberId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	req.UserId = memberId
	err = service.AdminTwoFactor().Disable(ctx, &req.TwoFactorDisableInp)
	return
}

// Verify 验证双因子认证
func (c *cTwoFactor) Verify(ctx context.Context, req *twofactor.VerifyReq) (res *twofactor.VerifyRes, err error) {
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

	res = &twofactor.VerifyRes{}
	return
}

// GetStatus 获取双因子认证状态
func (c *cTwoFactor) GetStatus(ctx context.Context, req *twofactor.GetStatusReq) (res *twofactor.GetStatusRes, err error) {
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

	res = &twofactor.GetStatusRes{
		TwoFactorGetStatusModel: result,
	}
	return
}

// RegenerateBackupCodes 重新生成备用恢复码
func (c *cTwoFactor) RegenerateBackupCodes(ctx context.Context, req *twofactor.RegenerateBackupCodesReq) (res *twofactor.RegenerateBackupCodesRes, err error) {
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

	res = &twofactor.RegenerateBackupCodesRes{
		TwoFactorRegenerateBackupCodesModel: result,
	}
	return
}
