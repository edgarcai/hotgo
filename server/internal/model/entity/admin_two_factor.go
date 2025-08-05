// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminTwoFactor is the golang structure for table admin_two_factor.
type AdminTwoFactor struct {
	Id          uint64      `json:"id"          description:"主键ID"`
	MemberId    uint64      `json:"memberId"    description:"管理员ID"`
	SecretKey   string      `json:"secretKey"   description:"加密后的TOTP密钥"`
	BackupCodes *gjson.Json `json:"backupCodes" description:"备用恢复码（哈希后）"`
	IsEnabled   int         `json:"isEnabled"   description:"是否启用（0:未启用 1:已启用）"`
	EnabledAt   *gtime.Time `json:"enabledAt"   description:"启用时间"`
	LastUsedAt  *gtime.Time `json:"lastUsedAt"  description:"最后使用时间"`
	CreatedAt   *gtime.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   description:"更新时间"`
}
