// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixChapterDao is the data access object for table anix_chapters.
type AnixChapterDao struct {
	table   string             // table is the underlying table name of the DAO.
	group   string             // group is the database configuration group name of current DAO.
	columns AnixChapterColumns // columns contains all the column names of Table for convenient usage.
}

// AnixChapterColumns defines and stores column names for table anix_chapters.
type AnixChapterColumns struct {
	Id          string // 章节ID
	ComicId     string // 漫画ID
	Title       string // 章节标题
	ChapterNum  string // 章节序号
	IsVip       string // 是否VIP:0=否,1=是
	ViewCount   string // 浏览次数
	Sort        string // 排序
	Status      string // 状态:1=正常,2=禁用
	PublishedAt string // 发布时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// anixChapterColumns holds the columns for table anix_chapters.
var anixChapterColumns = AnixChapterColumns{
	Id:          "id",
	ComicId:     "comic_id",
	Title:       "title",
	ChapterNum:  "chapter_num",
	IsVip:       "is_vip",
	ViewCount:   "view_count",
	Sort:        "sort",
	Status:      "status",
	PublishedAt: "published_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewAnixChapterDao creates and returns a new DAO object for table data access.
func NewAnixChapterDao() *AnixChapterDao {
	return &AnixChapterDao{
		group:   "default",
		table:   "anix_chapters",
		columns: anixChapterColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixChapterDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixChapterDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixChapterDao) Columns() AnixChapterColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixChapterDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixChapterDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}