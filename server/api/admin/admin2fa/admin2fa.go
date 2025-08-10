// Package admin2fa
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin2fa

import (
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"

	"github.com/gogf/gf/v2/frame/g"
)

// SetupReq 设置2FA请求
type SetupReq struct {
	g.Meta `path:"/admin2fa/setup" method:"post" tags:"2FA管理" summary:"设置2FA"`
	adminin.Admin2faSetupInp
}

type SetupRes struct {
	SecretKey string `json:"secretKey" dc:"密钥"`
	QrCode    string `json:"qrCode" dc:"二维码"`
}

// VerifySetupReq 验证2FA设置请求
type VerifySetupReq struct {
	g.Meta `path:"/admin2fa/verify-setup" method:"post" tags:"2FA管理" summary:"验证2FA设置"`
	adminin.Admin2faVerifySetupInp
}

type VerifySetupRes struct{}

// VerifyReq 验证2FA请求
type VerifyReq struct {
	g.Meta `path:"/admin2fa/verify" method:"post" tags:"2FA管理" summary:"验证2FA"`
	adminin.Admin2faVerifyInp
}

type VerifyRes struct{}

// DisableReq 禁用2FA请求
type DisableReq struct {
	g.Meta `path:"/admin2fa/disable" method:"post" tags:"2FA管理" summary:"禁用2FA"`
	adminin.Admin2faDisableInp
}

type DisableRes struct{}

// StatusReq 获取2FA状态请求
type StatusReq struct {
	g.Meta `path:"/admin2fa/status" method:"get" tags:"2FA管理" summary:"获取2FA状态"`
	adminin.Admin2faStatusInp
}

type StatusRes struct {
	IsEnabled  bool `json:"isEnabled" dc:"是否启用"`
	IsVerified bool `json:"isVerified" dc:"是否已验证"`
}

// BackupCodesReq 获取备用码请求
type BackupCodesReq struct {
	g.Meta `path:"/admin2fa/backup-codes" method:"get" tags:"2FA管理" summary:"获取备用码"`
	adminin.Admin2faBackupCodesInp
}

type BackupCodesRes struct {
	BackupCodes []string `json:"backupCodes" dc:"备用码列表"`
}

// RegenerateBackupCodesReq 重新生成备用码请求
type RegenerateBackupCodesReq struct {
	g.Meta `path:"/admin2fa/regenerate-backup-codes" method:"post" tags:"2FA管理" summary:"重新生成备用码"`
	adminin.Admin2faRegenerateBackupCodesInp
}

type RegenerateBackupCodesRes struct {
	BackupCodes []string `json:"backupCodes" dc:"备用码列表"`
}

// UseBackupCodeReq 使用备用码请求
type UseBackupCodeReq struct {
	g.Meta `path:"/admin2fa/use-backup-code" method:"post" tags:"2FA管理" summary:"使用备用码"`
	adminin.Admin2faUseBackupCodeInp
}

type UseBackupCodeRes struct{}

// LogListReq 获取2FA日志列表请求
type LogListReq struct {
	g.Meta `path:"/admin2fa/log/list" method:"get" tags:"2FA管理" summary:"获取2FA日志列表"`
	adminin.Admin2faLogListInp
}

type LogListRes struct {
	List     []*entity.Admin2faLog `json:"list" dc:"日志列表"`
	Page     int                   `json:"page" dc:"分页号码"`
	PageSize int                   `json:"pageSize" dc:"分页数量"`
	Total    int64                 `json:"total" dc:"总数"`
}