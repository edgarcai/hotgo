// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Admin2faLogDao is the data access object for table hg_admin_2fa_log.
type Admin2faLogDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns Admin2faLogColumns  // columns contains all the column names of Table for convenient usage.
}

// Admin2faLogColumns defines and stores column names for table hg_admin_2fa_log.
type Admin2faLogColumns struct {
	Id        string // 主键ID
	AdminId   string // 管理员ID
	Action    string // 操作类型(setup:设置 verify:验证 disable:禁用 backup_used:使用备用码)
	Result    string // 操作结果(success:成功 failed:失败)
	Ip        string // IP地址
	UserAgent string // 用户代理
	Remark    string // 备注信息
	CreatedAt string // 创建时间
}

// admin2faLogColumns holds the columns for table hg_admin_2fa_log.
var admin2faLogColumns = Admin2faLogColumns{
	Id:        "id",
	AdminId:   "admin_id",
	Action:    "action",
	Result:    "result",
	Ip:        "ip",
	UserAgent: "user_agent",
	Remark:    "remark",
	CreatedAt: "created_at",
}

// NewAdmin2faLogDao creates and returns a new DAO object for table data access.
func NewAdmin2faLogDao() *Admin2faLogDao {
	return &Admin2faLogDao{
		group:   "default",
		table:   "hg_admin_2fa_log",
		columns: admin2faLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *Admin2faLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *Admin2faLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *Admin2faLogDao) Columns() Admin2faLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *Admin2faLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *Admin2faLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *Admin2faLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}