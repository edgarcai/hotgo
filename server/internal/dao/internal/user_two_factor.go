// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTwoFactorDao is the data access object for table hg_user_two_factor.
type UserTwoFactorDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns UserTwoFactorColumns // columns contains all the column names of Table for convenient usage.
}

// UserTwoFactorColumns defines and stores column names for table hg_user_two_factor.
type UserTwoFactorColumns struct {
	Id               string // 主键ID
	UserId           string // 用户ID，关联hg_admin_member.id
	SecretKey        string // 加密后的TOTP密钥
	IsEnabled        string // 是否启用2FA，0=禁用，1=启用
	BackupCodesCount string // 剩余备用恢复码数量
	LastUsedAt       string // 最后使用时间
	CreatedAt        string // 创建时间
	UpdatedAt        string // 更新时间
}

// userTwoFactorColumns holds the columns for table hg_user_two_factor.
var userTwoFactorColumns = UserTwoFactorColumns{
	Id:               "id",
	UserId:           "user_id",
	SecretKey:        "secret_key",
	IsEnabled:        "is_enabled",
	BackupCodesCount: "backup_codes_count",
	LastUsedAt:       "last_used_at",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
}

// NewUserTwoFactorDao creates and returns a new DAO object for table data access.
func NewUserTwoFactorDao() *UserTwoFactorDao {
	return &UserTwoFactorDao{
		group:   "default",
		table:   "hg_user_two_factor",
		columns: userTwoFactorColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserTwoFactorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserTwoFactorDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserTwoFactorDao) Columns() UserTwoFactorColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserTwoFactorDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserTwoFactorDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserTwoFactorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
