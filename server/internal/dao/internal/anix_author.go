// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixAuthorDao is the data access object for table anix_authors.
type AnixAuthorDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns AnixAuthorColumns // columns contains all the column names of Table for convenient usage.
}

// AnixAuthorColumns defines and stores column names for table anix_authors.
type AnixAuthorColumns struct {
	Id          string // 作者ID
	Name        string // 作者姓名
	Avatar      string // 作者头像
	Description string // 作者简介
	Status      string // 状态:1=正常,2=禁用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// anixAuthorColumns holds the columns for table anix_authors.
var anixAuthorColumns = AnixAuthorColumns{
	Id:          "id",
	Name:        "name",
	Avatar:      "avatar",
	Description: "description",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewAnixAuthorDao creates and returns a new DAO object for table data access.
func NewAnixAuthorDao() *AnixAuthorDao {
	return &AnixAuthorDao{
		group:   "default",
		table:   "anix_authors",
		columns: anixAuthorColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixAuthorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixAuthorDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixAuthorDao) Columns() AnixAuthorColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixAuthorDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixAuthorDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}