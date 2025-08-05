// ================================================================================
// 双因子认证输入输出模型定义
// ================================================================================

package adminin

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TwoFactorEnableInp 启用双因子认证输入
type TwoFactorEnableInp struct {
	UserId int64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
}

// TwoFactorEnableModel 启用双因子认证输出
type TwoFactorEnableModel struct {
	Secret      string   `json:"secret"      dc:"密钥（格式化显示）"`
	QRCodeURL   string   `json:"qrCodeURL"   dc:"二维码URL"`
	BackupCodes []string `json:"backupCodes" dc:"备用恢复码"`
}

// TwoFactorConfirmEnableInp 确认启用双因子认证输入
type TwoFactorConfirmEnableInp struct {
	UserId int64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Code   string `json:"code" v:"required|length:6,6#请输入验证码|验证码必须为6位" dc:"TOTP验证码"`
}

// TwoFactorDisableInp 禁用双因子认证输入
type TwoFactorDisableInp struct {
	UserId int64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Code   string `json:"code" v:"required#验证码不能为空" dc:"TOTP验证码或备用恢复码"`
}

// TwoFactorVerifyInp 验证双因子认证输入
type TwoFactorVerifyInp struct {
	UserId     int64  `json:"userId"     v:"required|min:1#请指定用户ID|用户ID无效" dc:"用户ID"`
	Code       string `json:"code"       dc:"TOTP验证码"`
	BackupCode string `json:"backupCode" dc:"备用恢复码"`
}

// TwoFactorVerifyModel 验证双因子认证输出
type TwoFactorVerifyModel struct {
	IsValid bool `json:"isValid" dc:"验证是否通过"`
}

// TwoFactorGetStatusInp 获取双因子认证状态输入
type TwoFactorGetStatusInp struct {
	UserId int64 `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
}

// TwoFactorGetStatusModel 获取双因子认证状态输出
type TwoFactorGetStatusModel struct {
	IsEnabled        bool        `json:"isEnabled"        dc:"是否已启用"`
	BackupCodesCount int         `json:"backupCodesCount" dc:"剩余备用恢复码数量"`
	LastUsedAt       *gtime.Time `json:"lastUsedAt"       dc:"最后使用时间"`
}

// TwoFactorRegenerateBackupCodesInp 重新生成备用恢复码输入
type TwoFactorRegenerateBackupCodesInp struct {
	UserId int64  `json:"userId" v:"required#用户ID不能为空" dc:"用户ID"`
	Code   string `json:"code" v:"required|length:6,6#请输入验证码|验证码必须为6位" dc:"TOTP验证码"`
}

// TwoFactorRegenerateBackupCodesModel 重新生成备用恢复码输出
type TwoFactorRegenerateBackupCodesModel struct {
	BackupCodes []string `json:"backupCodes" dc:"新的备用恢复码列表"`
}
