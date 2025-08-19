// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixChapterPageDao is the data access object for table anix_chapter_pages.
type AnixChapterPageDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns AnixChapterPageColumns // columns contains all the column names of Table for convenient usage.
}

// AnixChapterPageColumns defines and stores column names for table anix_chapter_pages.
type AnixChapterPageColumns struct {
	Id        string // 页面ID
	ChapterId string // 章节ID
	PageNum   string // 页码
	ImageUrl  string // 图片URL
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// anixChapterPageColumns holds the columns for table anix_chapter_pages.
var anixChapterPageColumns = AnixChapterPageColumns{
	Id:        "id",
	ChapterId: "chapter_id",
	PageNum:   "page_num",
	ImageUrl:  "image_url",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewAnixChapterPageDao creates and returns a new DAO object for table data access.
func NewAnixChapterPageDao() *AnixChapterPageDao {
	return &AnixChapterPageDao{
		group:   "default",
		table:   "anix_chapter_pages",
		columns: anixChapterPageColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixChapterPageDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixChapterPageDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixChapterPageDao) Columns() AnixChapterPageColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixChapterPageDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixChapterPageDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}