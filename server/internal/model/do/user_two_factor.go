// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserTwoFactor is the golang structure of table hg_user_two_factor for DAO operations like Where/Data.
type UserTwoFactor struct {
	g.Meta           `orm:"table:hg_user_two_factor, do:true"`
	Id               interface{} // 主键ID
	UserId           interface{} // 用户ID，关联hg_admin_member.id
	SecretKey        interface{} // 加密后的TOTP密钥
	IsEnabled        interface{} // 是否启用2FA，0=禁用，1=启用
	BackupCodesCount interface{} // 剩余备用恢复码数量
	LastUsedAt       *gtime.Time // 最后使用时间
	CreatedAt        *gtime.Time // 创建时间
	UpdatedAt        *gtime.Time // 更新时间
}
