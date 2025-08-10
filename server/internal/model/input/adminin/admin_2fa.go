package adminin

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/entity"
)

// Admin2faSetupInp 2FA设置输入
type Admin2faSetupInp struct {
	g.Meta `path:"/admin2fa/setup" method:"post" tags:"管理员2FA" summary:"设置2FA"`
}

// Admin2faSetupModel 2FA设置输出
type Admin2faSetupModel struct {
	SecretKey string `json:"secretKey" dc:"TOTP密钥"`
	QRCodeURL string `json:"qrCodeURL" dc:"二维码URL"`
	BackupCodes []string `json:"backupCodes" dc:"备用恢复码"`
}

// Admin2faVerifySetupInp 验证2FA设置输入
type Admin2faVerifySetupInp struct {
	g.Meta `path:"/admin2fa/verify-setup" method:"post" tags:"管理员2FA" summary:"验证2FA设置"`
	Code   string `json:"code" v:"required|length:6,6#验证码不能为空|验证码长度必须为6位" dc:"验证码"`
}

// Admin2faVerifyInp 2FA验证输入
type Admin2faVerifyInp struct {
	g.Meta `path:"/admin2fa/verify" method:"post" tags:"管理员2FA" summary:"2FA验证"`
	Code   string `json:"code" v:"required|length:6,6#验证码不能为空|验证码长度必须为6位" dc:"验证码"`
	AdminId int64 `json:"-" dc:"管理员ID"`
}

// Admin2faDisableInp 禁用2FA输入
type Admin2faDisableInp struct {
	g.Meta `path:"/admin2fa/disable" method:"post" tags:"管理员2FA" summary:"禁用2FA"`
	Code   string `json:"code" v:"required|length:6,6#验证码不能为空|验证码长度必须为6位" dc:"验证码"`
}

// Admin2faStatusInp 获取2FA状态输入
type Admin2faStatusInp struct {
	g.Meta `path:"/admin2fa/status" method:"get" tags:"管理员2FA" summary:"获取2FA状态"`
}

// Admin2faStatusModel 2FA状态输出
type Admin2faStatusModel struct {
	IsEnabled  bool `json:"isEnabled" dc:"是否启用"`
	IsVerified bool `json:"isVerified" dc:"是否已验证"`
	HasBackupCodes bool `json:"hasBackupCodes" dc:"是否有备用码"`
}

// Admin2faBackupCodesInp 获取备用码输入
type Admin2faBackupCodesInp struct {
	g.Meta `path:"/admin2fa/backup-codes" method:"get" tags:"管理员2FA" summary:"获取备用码"`
}

// Admin2faBackupCodesModel 备用码输出
type Admin2faBackupCodesModel struct {
	BackupCodes []string `json:"backupCodes" dc:"备用恢复码"`
}

// Admin2faRegenerateBackupCodesInp 重新生成备用码输入
type Admin2faRegenerateBackupCodesInp struct {
	g.Meta `path:"/admin2fa/regenerate-backup-codes" method:"post" tags:"管理员2FA" summary:"重新生成备用码"`
	Code   string `json:"code" v:"required|length:6,6#验证码不能为空|验证码长度必须为6位" dc:"验证码"`
}

// Admin2faUseBackupCodeInp 使用备用码输入
type Admin2faUseBackupCodeInp struct {
	g.Meta `path:"/admin2fa/use-backup-code" method:"post" tags:"管理员2FA" summary:"使用备用码"`
	Code    string `json:"code" v:"required|length:8,8#备用码不能为空|备用码长度必须为8位" dc:"备用码"`
	AdminId int64  `json:"-" dc:"管理员ID"`
}

// Admin2faLogListInp 2FA日志列表输入
type Admin2faLogListInp struct {
	g.Meta `path:"/admin2fa/log/list" method:"get" tags:"管理员2FA" summary:"2FA日志列表"`
	Page     int    `json:"page" v:"min:0#分页号码错误" dc:"分页号码，默认1"`
	PageSize int    `json:"pageSize" v:"max:50#分页数量最大50条" dc:"分页数量，默认10"`
	Action   string `json:"action" dc:"操作类型"`
	Result   string `json:"result" dc:"操作结果"`
}

// Admin2faLogListModel 2FA日志列表输出
type Admin2faLogListModel struct {
	List []*entity.Admin2faLog `json:"list" dc:"日志列表"`
	Page int                   `json:"page" dc:"分页号码"`
	PageSize int              `json:"pageSize" dc:"分页数量"`
	Total    int64            `json:"total" dc:"总数"`
}