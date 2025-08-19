// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnixUserActionDao is the data access object for table anix_user_actions.
type AnixUserActionDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns AnixUserActionColumns // columns contains all the column names of Table for convenient usage.
}

// AnixUserActionColumns defines and stores column names for table anix_user_actions.
type AnixUserActionColumns struct {
	Id         string // 行为ID
	UserId     string // 用户ID
	ComicId    string // 漫画ID
	ActionType string // 行为类型:1=收藏,2=点赞,3=评分
	ActionData string // 行为数据(JSON格式)
	CreatedAt  string // 创建时间
	UpdatedAt  string // 更新时间
}

// anixUserActionColumns holds the columns for table anix_user_actions.
var anixUserActionColumns = AnixUserActionColumns{
	Id:         "id",
	UserId:     "user_id",
	ComicId:    "comic_id",
	ActionType: "action_type",
	ActionData: "action_data",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewAnixUserActionDao creates and returns a new DAO object for table data access.
func NewAnixUserActionDao() *AnixUserActionDao {
	return &AnixUserActionDao{
		group:   "default",
		table:   "anix_user_actions",
		columns: anixUserActionColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnixUserActionDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnixUserActionDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnixUserActionDao) Columns() AnixUserActionColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnixUserActionDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnixUserActionDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}