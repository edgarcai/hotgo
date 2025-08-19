// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixComicDao is the data access object for table anix_comics.
type AnixComicDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns AnixComicColumns // columns contains all the column names of Table for convenient usage.
}

// AnixComicColumns defines and stores column names for table anix_comics.
type AnixComicColumns struct {
	Id          string // 漫画ID
	Title       string // 漫画标题
	Cover       string // 封面图片
	Description string // 漫画简介
	AuthorId    string // 作者ID
	Status      string // 状态:1=连载中,2=已完结,3=暂停更新
	IsVip       string // 是否VIP:0=否,1=是
	ViewCount   string // 浏览次数
	LikeCount   string // 点赞次数
	Sort        string // 排序
	PublishedAt string // 发布时间
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
	DeletedAt   string // 删除时间
}

// anixComicColumns holds the columns for table anix_comics.
var anixComicColumns = AnixComicColumns{
	Id:          "id",
	Title:       "title",
	Cover:       "cover",
	Description: "description",
	AuthorId:    "author_id",
	Status:      "status",
	IsVip:       "is_vip",
	ViewCount:   "view_count",
	LikeCount:   "like_count",
	Sort:        "sort",
	PublishedAt: "published_at",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewAnixComicDao creates and returns a new DAO object for table data access.
func NewAnixComicDao() *AnixComicDao {
	return &AnixComicDao{
		group:   "default",
		table:   "anix_comics",
		columns: anixComicColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixComicDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixComicDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixComicDao) Columns() AnixComicColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixComicDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixComicDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}