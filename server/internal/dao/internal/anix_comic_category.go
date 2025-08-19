// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixComicCategoryDao is the data access object for table anix_comic_categories.
type AnixComicCategoryDao struct {
	table   string                   // table is the underlying table name of the DAO.
	group   string                   // group is the database configuration group name of current DAO.
	columns AnixComicCategoryColumns // columns contains all the column names of Table for convenient usage.
}

// AnixComicCategoryColumns defines and stores column names for table anix_comic_categories.
type AnixComicCategoryColumns struct {
	Id         string // ID
	ComicId    string // 漫画ID
	CategoryId string // 分类ID
	CreatedAt  string // 创建时间
}

// anixComicCategoryColumns holds the columns for table anix_comic_categories.
var anixComicCategoryColumns = AnixComicCategoryColumns{
	Id:         "id",
	ComicId:    "comic_id",
	CategoryId: "category_id",
	CreatedAt:  "created_at",
}

// NewAnixComicCategoryDao creates and returns a new DAO object for table data access.
func NewAnixComicCategoryDao() *AnixComicCategoryDao {
	return &AnixComicCategoryDao{
		group:   "default",
		table:   "anix_comic_categories",
		columns: anixComicCategoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixComicCategoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixComicCategoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixComicCategoryDao) Columns() AnixComicCategoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixComicCategoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixComicCategoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}