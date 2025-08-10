// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"hotgo/api/admin/admin2fa"
	"hotgo/internal/service"
)

var Admin2fa = cAdmin2fa{}

type cAdmin2fa struct{}

// Setup 设置2FA
func (c *cAdmin2fa) Setup(ctx context.Context, req *admin2fa.SetupReq) (res *admin2fa.SetupRes, err error) {
	model, err := service.Admin2fa().Setup(ctx, &req.Admin2faSetupInp)
	if err != nil {
		return nil, err
	}

	res = &admin2fa.SetupRes{
		SecretKey: model.SecretKey,
		QrCode:    model.QRCodeURL,
	}
	return
}

// VerifySetup 验证2FA设置
func (c *cAdmin2fa) VerifySetup(ctx context.Context, req *admin2fa.VerifySetupReq) (res *admin2fa.VerifySetupRes, err error) {
	err = service.Admin2fa().VerifySetup(ctx, &req.Admin2faVerifySetupInp)
	return
}

// Verify 验证2FA
func (c *cAdmin2fa) Verify(ctx context.Context, req *admin2fa.VerifyReq) (res *admin2fa.VerifyRes, err error) {
	err = service.Admin2fa().Verify(ctx, &req.Admin2faVerifyInp)
	return
}

// Disable 禁用2FA
func (c *cAdmin2fa) Disable(ctx context.Context, req *admin2fa.DisableReq) (res *admin2fa.DisableRes, err error) {
	err = service.Admin2fa().Disable(ctx, &req.Admin2faDisableInp)
	return
}

// Status 获取2FA状态
func (c *cAdmin2fa) Status(ctx context.Context, req *admin2fa.StatusReq) (res *admin2fa.StatusRes, err error) {
	model, err := service.Admin2fa().Status(ctx, &req.Admin2faStatusInp)
	if err != nil {
		return nil, err
	}

	res = &admin2fa.StatusRes{
		IsEnabled:  model.IsEnabled,
		IsVerified: model.IsVerified,
	}
	return
}

// BackupCodes 获取备用码
func (c *cAdmin2fa) BackupCodes(ctx context.Context, req *admin2fa.BackupCodesReq) (res *admin2fa.BackupCodesRes, err error) {
	model, err := service.Admin2fa().BackupCodes(ctx, &req.Admin2faBackupCodesInp)
	if err != nil {
		return nil, err
	}

	res = &admin2fa.BackupCodesRes{
		BackupCodes: model.BackupCodes,
	}
	return
}

// RegenerateBackupCodes 重新生成备用码
func (c *cAdmin2fa) RegenerateBackupCodes(ctx context.Context, req *admin2fa.RegenerateBackupCodesReq) (res *admin2fa.RegenerateBackupCodesRes, err error) {
	model, err := service.Admin2fa().RegenerateBackupCodes(ctx, &req.Admin2faRegenerateBackupCodesInp)
	if err != nil {
		return nil, err
	}

	res = &admin2fa.RegenerateBackupCodesRes{
		BackupCodes: model.BackupCodes,
	}
	return
}

// UseBackupCode 使用备用码
func (c *cAdmin2fa) UseBackupCode(ctx context.Context, req *admin2fa.UseBackupCodeReq) (res *admin2fa.UseBackupCodeRes, err error) {
	err = service.Admin2fa().UseBackupCode(ctx, &req.Admin2faUseBackupCodeInp)
	return
}

// LogList 获取2FA日志列表
func (c *cAdmin2fa) LogList(ctx context.Context, req *admin2fa.LogListReq) (res *admin2fa.LogListRes, err error) {
	model, err := service.Admin2fa().LogList(ctx, &req.Admin2faLogListInp)
	if err != nil {
		return nil, err
	}

	res = &admin2fa.LogListRes{
		List:     model.List,
		Page:     model.Page,
		PageSize: model.PageSize,
		Total:    model.Total,
	}
	return
}