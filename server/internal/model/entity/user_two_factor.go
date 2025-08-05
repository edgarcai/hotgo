// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTwoFactor is the golang structure for table user_two_factor.
type UserTwoFactor struct {
	Id               int64       `json:"id"               description:"主键ID"`
	UserId           int64       `json:"userId"           description:"用户ID，关联hg_admin_member.id"`
	SecretKey        string      `json:"secretKey"        description:"加密后的TOTP密钥"`
	IsEnabled        int         `json:"isEnabled"        description:"是否启用2FA，0=禁用，1=启用"`
	BackupCodesCount int         `json:"backupCodesCount" description:"剩余备用恢复码数量"`
	LastUsedAt       *gtime.Time `json:"lastUsedAt"       description:"最后使用时间"`
	CreatedAt        *gtime.Time `json:"createdAt"        description:"创建时间"`
	UpdatedAt        *gtime.Time `json:"updatedAt"        description:"更新时间"`
}
