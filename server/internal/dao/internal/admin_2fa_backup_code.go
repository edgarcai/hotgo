// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Admin2faBackupCodeDao is the data access object for table hg_admin_2fa_backup_code.
type Admin2faBackupCodeDao struct {
	table   string                      // table is the underlying table name of the DAO.
	group   string                      // group is the database configuration group name of current DAO.
	columns Admin2faBackupCodeColumns   // columns contains all the column names of Table for convenient usage.
}

// Admin2faBackupCodeColumns defines and stores column names for table hg_admin_2fa_backup_code.
type Admin2faBackupCodeColumns struct {
	Id        string // 主键ID
	AdminId   string // 管理员ID
	Code      string // 备用恢复码
	IsUsed    string // 是否已使用(1:已使用 2:未使用)
	UsedAt    string // 使用时间
	CreatedAt string // 创建时间
}

// admin2faBackupCodeColumns holds the columns for table hg_admin_2fa_backup_code.
var admin2faBackupCodeColumns = Admin2faBackupCodeColumns{
	Id:        "id",
	AdminId:   "admin_id",
	Code:      "code",
	IsUsed:    "is_used",
	UsedAt:    "used_at",
	CreatedAt: "created_at",
}

// NewAdmin2faBackupCodeDao creates and returns a new DAO object for table data access.
func NewAdmin2faBackupCodeDao() *Admin2faBackupCodeDao {
	return &Admin2faBackupCodeDao{
		group:   "default",
		table:   "hg_admin_2fa_backup_code",
		columns: admin2faBackupCodeColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *Admin2faBackupCodeDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *Admin2faBackupCodeDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *Admin2faBackupCodeDao) Columns() Admin2faBackupCodeColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *Admin2faBackupCodeDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *Admin2faBackupCodeDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *Admin2faBackupCodeDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}