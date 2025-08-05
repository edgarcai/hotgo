// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminTwoFactor is the golang structure of table hg_admin_two_factor for DAO operations like Where/Data.
type AdminTwoFactor struct {
	g.Meta      `orm:"table:hg_admin_two_factor, do:true"`
	Id          interface{} // 主键ID
	MemberId    interface{} // 管理员ID
	SecretKey   interface{} // 加密后的TOTP密钥
	BackupCodes *gjson.Json // 备用恢复码（哈希后）
	IsEnabled   interface{} // 是否启用（0:未启用 1:已启用）
	EnabledAt   *gtime.Time // 启用时间
	LastUsedAt  *gtime.Time // 最后使用时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
