// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixCategoryDao is the data access object for table anix_categories.
type AnixCategoryDao struct {
	table   string              // table is the underlying table name of the DAO.
	group   string              // group is the database configuration group name of current DAO.
	columns AnixCategoryColumns // columns contains all the column names of Table for convenient usage.
}

// AnixCategoryColumns defines and stores column names for table anix_categories.
type AnixCategoryColumns struct {
	Id          string // 分类ID
	Name        string // 分类名称
	Description string // 分类描述
	Sort        string // 排序
	Status      string // 状态:1=正常,2=禁用
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// anixCategoryColumns holds the columns for table anix_categories.
var anixCategoryColumns = AnixCategoryColumns{
	Id:          "id",
	Name:        "name",
	Description: "description",
	Sort:        "sort",
	Status:      "status",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewAnixCategoryDao creates and returns a new DAO object for table data access.
func NewAnixCategoryDao() *AnixCategoryDao {
	return &AnixCategoryDao{
		group:   "default",
		table:   "anix_categories",
		columns: anixCategoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixCategoryDao) Columns() AnixCategoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixCategoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}