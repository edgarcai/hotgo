// Package twofactor
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package twofactor

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/adminin"
)

// EnableReq 启用双因子认证
type EnableReq struct {
	g.Meta `path:"/twofactor/enable" method:"post" tags:"双因子认证" summary:"启用双因子认证"`
	adminin.TwoFactorEnableInp
}

type EnableRes struct {
	*adminin.TwoFactorEnableModel
}

// ConfirmEnableReq 确认启用双因子认证
type ConfirmEnableReq struct {
	g.Meta `path:"/twofactor/confirmEnable" method:"post" tags:"双因子认证" summary:"确认启用双因子认证"`
	adminin.TwoFactorConfirmEnableInp
}

type ConfirmEnableRes struct{}

// DisableReq 禁用双因子认证
type DisableReq struct {
	g.Meta `path:"/twofactor/disable" method:"post" tags:"双因子认证" summary:"禁用双因子认证"`
	adminin.TwoFactorDisableInp
}

type DisableRes struct{}

// VerifyReq 验证双因子认证
type VerifyReq struct {
	g.Meta `path:"/twofactor/verify" method:"post" tags:"双因子认证" summary:"验证双因子认证"`
	adminin.TwoFactorVerifyInp
}

type VerifyRes struct {
	*adminin.TwoFactorVerifyModel
}

// GetStatusReq 获取双因子认证状态
type GetStatusReq struct {
	g.Meta `path:"/twofactor/status" method:"get" tags:"双因子认证" summary:"获取双因子认证状态"`
	adminin.TwoFactorGetStatusInp
}

type GetStatusRes struct {
	*adminin.TwoFactorGetStatusModel
}

// RegenerateBackupCodesReq 重新生成备用恢复码
type RegenerateBackupCodesReq struct {
	g.Meta `path:"/twofactor/regenerateBackupCodes" method:"post" tags:"双因子认证" summary:"重新生成备用恢复码"`
	adminin.TwoFactorRegenerateBackupCodesInp
}

type RegenerateBackupCodesRes struct {
	*adminin.TwoFactorRegenerateBackupCodesModel
}
