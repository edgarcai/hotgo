// Package auth
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package auth

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/adminin"
)

// Enable2FAReq 启用双因子认证请求
type Enable2FAReq struct {
	g.Meta `path:"/api/admin/auth/enable-2fa" method:"post" tags:"管理员认证" summary:"启用双因子认证"`
	adminin.TwoFactorEnableInp
}

type Enable2FARes struct {
	*adminin.TwoFactorEnableModel
}

// Verify2FASetupReq 验证双因子认证设置请求
type Verify2FASetupReq struct {
	g.Meta `path:"/api/admin/auth/verify-2fa-setup" method:"post" tags:"管理员认证" summary:"验证双因子认证设置"`
	adminin.TwoFactorConfirmEnableInp
}

type Verify2FASetupRes struct{}

// Disable2FAReq 禁用双因子认证请求
type Disable2FAReq struct {
	g.Meta `path:"/api/admin/auth/disable-2fa" method:"post" tags:"管理员认证" summary:"禁用双因子认证"`
	adminin.TwoFactorDisableInp
}

type Disable2FARes struct{}

// Get2FAStatusReq 获取双因子认证状态请求
type Get2FAStatusReq struct {
	g.Meta `path:"/api/admin/auth/2fa-status" method:"get" tags:"管理员认证" summary:"获取双因子认证状态"`
	adminin.TwoFactorGetStatusInp
}

type Get2FAStatusRes struct {
	*adminin.TwoFactorGetStatusModel
}

// RegenerateBackupCodesReq 重新生成备用恢复码请求
type RegenerateBackupCodesReq struct {
	g.Meta `path:"/api/admin/auth/regenerate-backup-codes" method:"post" tags:"管理员认证" summary:"重新生成备用恢复码"`
	adminin.TwoFactorRegenerateBackupCodesInp
}

type RegenerateBackupCodesRes struct {
	*adminin.TwoFactorRegenerateBackupCodesModel
}

// Verify2FAReq 验证双因子认证请求
type Verify2FAReq struct {
	g.Meta `path:"/api/admin/auth/verify-2fa" method:"post" tags:"管理员认证" summary:"验证双因子认证"`
	adminin.TwoFactorVerifyInp
}

type Verify2FARes struct {
	*adminin.TwoFactorVerifyModel
}