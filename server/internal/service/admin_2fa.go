// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ================================================================================

package service

import (
	"context"

	"hotgo/internal/model/input/adminin"
)

type (
	IAdmin2fa interface {
		// Setup 设置2FA
		Setup(ctx context.Context, in *adminin.Admin2faSetupInp) (res *adminin.Admin2faSetupModel, err error)
		// VerifySetup 验证2FA设置
		VerifySetup(ctx context.Context, in *adminin.Admin2faVerifySetupInp) (err error)
		// Verify 验证2FA
		Verify(ctx context.Context, in *adminin.Admin2faVerifyInp) (err error)
		// Disable 禁用2FA
		Disable(ctx context.Context, in *adminin.Admin2faDisableInp) (err error)
		// Status 获取2FA状态
		Status(ctx context.Context, in *adminin.Admin2faStatusInp) (res *adminin.Admin2faStatusModel, err error)
		// BackupCodes 获取备用码
		BackupCodes(ctx context.Context, in *adminin.Admin2faBackupCodesInp) (res *adminin.Admin2faBackupCodesModel, err error)
		// RegenerateBackupCodes 重新生成备用码
		RegenerateBackupCodes(ctx context.Context, in *adminin.Admin2faRegenerateBackupCodesInp) (res *adminin.Admin2faBackupCodesModel, err error)
		// UseBackupCode 使用备用码
		UseBackupCode(ctx context.Context, in *adminin.Admin2faUseBackupCodeInp) (err error)
		// LogList 获取2FA日志列表
		LogList(ctx context.Context, in *adminin.Admin2faLogListInp) (res *adminin.Admin2faLogListModel, err error)
		// IsEnabled 检查用户是否启用了2FA
		IsEnabled(ctx context.Context, adminId uint64) (bool, error)
	}
)

var (
	localAdmin2fa IAdmin2fa
)

func Admin2fa() IAdmin2fa {
	if localAdmin2fa == nil {
		panic("implement not found for interface IAdmin2fa, forgot register?")
	}
	return localAdmin2fa
}

func RegisterAdmin2fa(i IAdmin2fa) {
	localAdmin2fa = i
}