// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserTwoFactorBackupCodesDao is the data access object for table hg_user_two_factor_backup_codes.
type UserTwoFactorBackupCodesDao struct {
	table   string                          // table is the underlying table name of the DAO.
	group   string                          // group is the database configuration group name of current DAO.
	columns UserTwoFactorBackupCodesColumns // columns contains all the column names of Table for convenient usage.
}

// UserTwoFactorBackupCodesColumns defines and stores column names for table hg_user_two_factor_backup_codes.
type UserTwoFactorBackupCodesColumns struct {
	Id        string // 主键ID
	UserId    string // 用户ID，关联hg_admin_member.id
	CodeHash  string // 备用恢复码的哈希值
	IsUsed    string // 是否已使用，0=未使用，1=已使用
	UsedAt    string // 使用时间
	CreatedAt string // 创建时间
}

// userTwoFactorBackupCodesColumns holds the columns for table hg_user_two_factor_backup_codes.
var userTwoFactorBackupCodesColumns = UserTwoFactorBackupCodesColumns{
	Id:        "id",
	UserId:    "user_id",
	CodeHash:  "code_hash",
	IsUsed:    "is_used",
	UsedAt:    "used_at",
	CreatedAt: "created_at",
}

// NewUserTwoFactorBackupCodesDao creates and returns a new DAO object for table data access.
func NewUserTwoFactorBackupCodesDao() *UserTwoFactorBackupCodesDao {
	return &UserTwoFactorBackupCodesDao{
		group:   "default",
		table:   "hg_user_two_factor_backup_codes",
		columns: userTwoFactorBackupCodesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UserTwoFactorBackupCodesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UserTwoFactorBackupCodesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UserTwoFactorBackupCodesDao) Columns() UserTwoFactorBackupCodesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UserTwoFactorBackupCodesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UserTwoFactorBackupCodesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UserTwoFactorBackupCodesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
