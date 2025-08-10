// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin2faConfig is the golang structure for table admin_2fa_config.
type Admin2faConfig struct {
	Id         int64       `json:"id"         orm:"id"          description:"主键ID"`
	AdminId    int64       `json:"adminId"    orm:"admin_id"     description:"管理员ID"`
	SecretKey  string      `json:"secretKey"  orm:"secret_key"   description:"TOTP密钥"`
	IsEnabled  int         `json:"isEnabled"  orm:"is_enabled"   description:"是否启用(1:启用 2:禁用)"`
	IsVerified int         `json:"isVerified" orm:"is_verified"  description:"是否已验证(1:已验证 2:未验证)"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   description:"创建时间"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"   description:"更新时间"`
}