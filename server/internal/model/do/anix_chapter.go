// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixChapter is the golang structure of table anix_chapters for DAO operations like Where/Data.
type AnixChapter struct {
	g.Meta      `orm:"table:anix_chapters, do:true"`
	Id          interface{} // 章节ID
	ComicId     interface{} // 漫画ID
	Title       interface{} // 章节标题
	ChapterNum  interface{} // 章节序号
	IsVip       interface{} // 是否VIP:0=否,1=是
	ViewCount   interface{} // 浏览次数
	Sort        interface{} // 排序
	Status      interface{} // 状态:1=正常,2=禁用
	PublishedAt *gtime.Time // 发布时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}