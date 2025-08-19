// Package entity AniX漫画相关实体定义
package entity

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/gogf/gf/v2/os/gtime"
)

// AnixComic 漫画实体
type AnixComic struct {
	Id           uint64      `json:"id"           orm:"id,primary"       description:"漫画ID"`
	Title        string      `json:"title"        orm:"title"            description:"漫画标题"`
	Description  string      `json:"description"  orm:"description"      description:"漫画简介"`
	CoverUrl     string      `json:"coverUrl"     orm:"cover_url"        description:"封面图片URL"`
	AuthorId     uint64      `json:"authorId"     orm:"author_id"        description:"作者ID"`
	Status       string      `json:"status"       orm:"status"           description:"连载状态"`
	Tags         TagsArray   `json:"tags"         orm:"tags"             description:"标签数组"`
	Views        uint64      `json:"views"        orm:"views"            description:"总浏览量"`
	Likes        uint64      `json:"likes"        orm:"likes"            description:"点赞数"`
	Favorites    uint64      `json:"favorites"    orm:"favorites"        description:"收藏数"`
	Rating       float64     `json:"rating"       orm:"rating"           description:"评分"`
	RatingCount  uint64      `json:"ratingCount"  orm:"rating_count"     description:"评分人数"`
	ChapterCount uint32      `json:"chapterCount" orm:"chapter_count"    description:"章节总数"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"       description:"创建时间"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"       description:"更新时间"`

	// 关联数据
	Author     *AnixAuthor      `json:"author,omitempty"     orm:"with:author_id=id"     description:"作者信息"`
	Categories []*AnixCategory  `json:"categories,omitempty" orm:"with:id=comic_id"       description:"分类信息"`
	Chapters   []*AnixChapter   `json:"chapters,omitempty"   orm:"with:id=comic_id"       description:"章节列表"`
}











// TagsArray 标签数组类型，用于JSON序列化
type TagsArray []string

// Scan 实现 sql.Scanner 接口，用于从数据库读取JSON数据
func (t *TagsArray) Scan(value interface{}) error {
	if value == nil {
		*t = TagsArray{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		*t = TagsArray{}
		return nil
	}
}

// Value 实现 driver.Valuer 接口，用于向数据库写入JSON数据
func (t TagsArray) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "[]", nil
	}
	return json.Marshal(t)
}

// String 返回标签的字符串表示
func (t TagsArray) String() string {
	data, _ := json.Marshal(t)
	return string(data)
}

// Contains 检查是否包含指定标签
func (t TagsArray) Contains(tag string) bool {
	for _, v := range t {
		if v == tag {
			return true
		}
	}
	return false
}

// Add 添加标签
func (t *TagsArray) Add(tag string) {
	if !t.Contains(tag) {
		*t = append(*t, tag)
	}
}

// Remove 移除标签
func (t *TagsArray) Remove(tag string) {
	for i, v := range *t {
		if v == tag {
			*t = append((*t)[:i], (*t)[i+1:]...)
			break
		}
	}
}