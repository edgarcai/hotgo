// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixUserDao is the data access object for table anix_users.
type AnixUserDao struct {
	table   string          // table is the underlying table name of the DAO.
	group   string          // group is the database configuration group name of current DAO.
	columns AnixUserColumns // columns contains all the column names of Table for convenient usage.
}

// AnixUserColumns defines and stores column names for table anix_users.
type AnixUserColumns struct {
	Id          string // 用户ID
	Username    string // 用户名
	Email       string // 邮箱
	Password    string // 密码
	Nickname    string // 昵称
	Avatar      string // 头像
	Gender      string // 性别:0=未知,1=男,2=女
	Birthday    string // 生日
	Phone       string // 手机号
	Status      string // 状态:1=正常,2=禁用
	LastLoginAt string // 最后登录时间
	LastLoginIp string // 最后登录IP
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// anixUserColumns holds the columns for table anix_users.
var anixUserColumns = AnixUserColumns{
	Id:          "id",
	Username:    "username",
	Email:       "email",
	Password:    "password",
	Nickname:    "nickname",
	Avatar:      "avatar",
	Gender:      "gender",
	Birthday:    "birthday",
	Phone:       "phone",
	Status:      "status",
	LastLoginAt: "last_login_at",
	LastLoginIp: "last_login_ip",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewAnixUserDao creates and returns a new DAO object for table data access.
func NewAnixUserDao() *AnixUserDao {
	return &AnixUserDao{
		group:   "default",
		table:   "anix_users",
		columns: anixUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixUserDao) Columns() AnixUserColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}