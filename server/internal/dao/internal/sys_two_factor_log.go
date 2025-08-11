// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysTwoFactorLogDao is the data access object for table sys_two_factor_log.
type SysTwoFactorLogDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns SysTwoFactorLogColumns // columns contains all the column names of Table for convenient usage.
}

// SysTwoFactorLogColumns defines and stores column names for table sys_two_factor_log.
type SysTwoFactorLogColumns struct {
	Id        string // 主键ID
	MemberId  string // 用户ID
	Username  string // 用户名
	Action    string // 操作类型：enable,disable,verify,backup_used
	Method    string // 2FA方法：totp,backup_code
	Result    string // 操作结果：success,failed
	Ip        string // IP地址
	UserAgent string // 用户代理
	Remark    string // 备注
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// sysTwoFactorLogColumns holds the columns for table sys_two_factor_log.
var sysTwoFactorLogColumns = SysTwoFactorLogColumns{
	Id:        "id",
	MemberId:  "member_id",
	Username:  "username",
	Action:    "action",
	Method:    "method",
	Result:    "result",
	Ip:        "ip",
	UserAgent: "user_agent",
	Remark:    "remark",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewSysTwoFactorLogDao creates and returns a new DAO object for table data access.
func NewSysTwoFactorLogDao() *SysTwoFactorLogDao {
	return &SysTwoFactorLogDao{
		group:   "default",
		table:   "sys_two_factor_log",
		columns: sysTwoFactorLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysTwoFactorLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *SysTwoFactorLogDao) Table() string {
	return dao.table
}

// Columns returns the columns of current dao.
func (dao *SysTwoFactorLogDao) Columns() SysTwoFactorLogColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *SysTwoFactorLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysTwoFactorLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
func (dao *SysTwoFactorLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}