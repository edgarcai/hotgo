// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Admin2faConfigDao is the data access object for table hg_admin_2fa_config.
type Admin2faConfigDao struct {
	table   string               // table is the underlying table name of the DAO.
	group   string               // group is the database configuration group name of current DAO.
	columns Admin2faConfigColumns // columns contains all the column names of Table for convenient usage.
}

// Admin2faConfigColumns defines and stores column names for table hg_admin_2fa_config.
type Admin2faConfigColumns struct {
	Id         string // 主键ID
	AdminId    string // 管理员ID
	SecretKey  string // TOTP密钥
	IsEnabled  string // 是否启用(1:启用 2:禁用)
	IsVerified string // 是否已验证(1:已验证 2:未验证)
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// admin2faConfigColumns holds the columns for table hg_admin_2fa_config.
var admin2faConfigColumns = Admin2faConfigColumns{
	Id:         "id",
	AdminId:    "admin_id",
	SecretKey:  "secret_key",
	IsEnabled:  "is_enabled",
	IsVerified: "is_verified",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewAdmin2faConfigDao creates and returns a new DAO object for table data access.
func NewAdmin2faConfigDao() *Admin2faConfigDao {
	return &Admin2faConfigDao{
		group:   "default",
		table:   "hg_admin_2fa_config",
		columns: admin2faConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *Admin2faConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *Admin2faConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *Admin2faConfigDao) Columns() Admin2faConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *Admin2faConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *Admin2faConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *Admin2faConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}