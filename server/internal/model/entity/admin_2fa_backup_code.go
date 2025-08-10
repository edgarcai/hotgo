// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin2faBackupCode is the golang structure for table admin_2fa_backup_code.
type Admin2faBackupCode struct {
	Id        int64       `json:"id"        orm:"id"         description:"主键ID"`
	AdminId   int64       `json:"adminId"   orm:"admin_id"   description:"管理员ID"`
	Code      string      `json:"code"      orm:"code"       description:"备用恢复码"`
	IsUsed    int         `json:"isUsed"    orm:"is_used"    description:"是否已使用(1:已使用 2:未使用)"`
	UsedAt    *gtime.Time `json:"usedAt"    orm:"used_at"    description:"使用时间"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}