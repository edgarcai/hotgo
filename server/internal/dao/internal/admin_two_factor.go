// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminTwoFactorDao is the data access object for table hg_admin_two_factor.
type AdminTwoFactorDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns AdminTwoFactorColumns // columns contains all the column names of Table for convenient usage.
}

// AdminTwoFactorColumns defines and stores column names for table hg_admin_two_factor.
type AdminTwoFactorColumns struct {
	Id          string // 主键ID
	MemberId    string // 管理员ID
	SecretKey   string // 加密后的TOTP密钥
	BackupCodes string // 备用恢复码（哈希后）
	IsEnabled   string // 是否启用（0:未启用 1:已启用）
	EnabledAt   string // 启用时间
	LastUsedAt  string // 最后使用时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// adminTwoFactorColumns holds the columns for table hg_admin_two_factor.
var adminTwoFactorColumns = AdminTwoFactorColumns{
	Id:          "id",
	MemberId:    "member_id",
	SecretKey:   "secret_key",
	BackupCodes: "backup_codes",
	IsEnabled:   "is_enabled",
	EnabledAt:   "enabled_at",
	LastUsedAt:  "last_used_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewAdminTwoFactorDao creates and returns a new DAO object for table data access.
func NewAdminTwoFactorDao() *AdminTwoFactorDao {
	return &AdminTwoFactorDao{
		group:   "default",
		table:   "hg_admin_two_factor",
		columns: adminTwoFactorColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AdminTwoFactorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AdminTwoFactorDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AdminTwoFactorDao) Columns() AdminTwoFactorColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AdminTwoFactorDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AdminTwoFactorDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AdminTwoFactorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
