// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixReadingHistoryDao is the data access object for table anix_reading_history.
type AnixReadingHistoryDao struct {
	table   string                    // table is the underlying table name of the DAO.
	group   string                    // group is the database configuration group name of current DAO.
	columns AnixReadingHistoryColumns // columns contains all the column names of Table for convenient usage.
}

// AnixReadingHistoryColumns defines and stores column names for table anix_reading_history.
type AnixReadingHistoryColumns struct {
	Id        string // 历史ID
	UserId    string // 用户ID
	ComicId   string // 漫画ID
	ChapterId string // 章节ID
	PageNum   string // 阅读到的页码
	CreatedAt string // 创建时间
	UpdatedAt string // 更新时间
}

// anixReadingHistoryColumns holds the columns for table anix_reading_history.
var anixReadingHistoryColumns = AnixReadingHistoryColumns{
	Id:        "id",
	UserId:    "user_id",
	ComicId:   "comic_id",
	ChapterId: "chapter_id",
	PageNum:   "page_num",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewAnixReadingHistoryDao creates and returns a new DAO object for table data access.
func NewAnixReadingHistoryDao() *AnixReadingHistoryDao {
	return &AnixReadingHistoryDao{
		group:   "default",
		table:   "anix_reading_history",
		columns: anixReadingHistoryColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixReadingHistoryDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixReadingHistoryDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixReadingHistoryDao) Columns() AnixReadingHistoryColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixReadingHistoryDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixReadingHistoryDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}