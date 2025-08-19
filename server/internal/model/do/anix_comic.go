// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AnixComic is the golang structure of table anix_comics for DAO operations like Where/Data.
type AnixComic struct {
	g.Meta      `orm:"table:anix_comics, do:true"`
	Id          interface{} // 漫画ID
	Title       interface{} // 漫画标题
	Cover       interface{} // 封面图片
	Description interface{} // 漫画简介
	AuthorId    interface{} // 作者ID
	Status      interface{} // 状态:1=连载中,2=已完结,3=暂停更新
	IsVip       interface{} // 是否VIP:0=否,1=是
	ViewCount   interface{} // 浏览次数
	LikeCount   interface{} // 点赞次数
	Sort        interface{} // 排序
	PublishedAt *gtime.Time // 发布时间
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
	DeletedAt   *gtime.Time // 删除时间
}