// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTwoFactorBackupCodes is the golang structure for table user_two_factor_backup_codes.
type UserTwoFactorBackupCodes struct {
	Id        int64       `json:"id"        description:"主键ID"`
	UserId    int64       `json:"userId"    description:"用户ID，关联hg_admin_member.id"`
	CodeHash  string      `json:"codeHash"  description:"备用恢复码的哈希值"`
	IsUsed    int         `json:"isUsed"    description:"是否已使用，0=未使用，1=已使用"`
	UsedAt    *gtime.Time `json:"usedAt"    description:"使用时间"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
}
