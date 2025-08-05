// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTwoFactorBackupCodes is the golang structure of table hg_user_two_factor_backup_codes for DAO operations like Where/Data.
type UserTwoFactorBackupCodes struct {
	g.Meta    `orm:"table:hg_user_two_factor_backup_codes, do:true"`
	Id        interface{} // 主键ID
	UserId    interface{} // 用户ID，关联hg_admin_member.id
	CodeHash  interface{} // 备用恢复码的哈希值
	IsUsed    interface{} // 是否已使用，0=未使用，1=已使用
	UsedAt    *gtime.Time // 使用时间
	CreatedAt *gtime.Time // 创建时间
}
