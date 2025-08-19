package service

import (
	"context"
	"time"
)

var (
	// localReadingHistory 本地阅读历史记录服务实例
	localReadingHistory IReadingHistory
)

// IReadingHistory 阅读历史记录服务接口
type IReadingHistory interface {
	// SaveReadingProgress 保存阅读进度
	SaveReadingProgress(ctx context.Context, in *SaveReadingProgressInp) (*SaveReadingProgressOut, error)
	// GetReadingHistory 获取阅读历史
	GetReadingHistory(ctx context.Context, in *GetReadingHistoryInp) (*GetReadingHistoryOut, error)
	// AddBookmark 添加书签
	AddBookmark(ctx context.Context, in *AddBookmarkInp) (*AddBookmarkOut, error)
	// RemoveBookmark 删除书签
	RemoveBookmark(ctx context.Context, in *RemoveBookmarkInp) (*RemoveBookmarkOut, error)
	// GetBookmarks 获取书签列表
	GetBookmarks(ctx context.Context, in *GetBookmarksInp) (*GetBookmarksOut, error)
	// UpdateReadingTime 更新阅读时长
	UpdateReadingTime(ctx context.Context, in *UpdateReadingTimeInp) (*UpdateReadingTimeOut, error)
	// GetReadingStats 获取阅读统计
	GetReadingStats(ctx context.Context, in *GetReadingStatsInp) (*GetReadingStatsOut, error)
}

// RegisterReadingHistory 注册阅读历史记录服务
func RegisterReadingHistory(i IReadingHistory) {
	localReadingHistory = i
}

// ReadingHistory 获取阅读历史记录服务实例
func ReadingHistory() IReadingHistory {
	if localReadingHistory == nil {
		panic("implement not found for interface IReadingHistory, forgot register?")
	}
	return localReadingHistory
}

// SaveReadingProgressInp 保存阅读进度输入
type SaveReadingProgressInp struct {
	UserId    int64 `json:"user_id" v:"required#用户ID不能为空"`
	ComicId   int64 `json:"comic_id" v:"required#漫画ID不能为空"`
	ChapterId int64 `json:"chapter_id" v:"required#章节ID不能为空"`
	PageIndex int   `json:"page_index" v:"required|min:0#页码不能为空|页码不能小于0"`
}

// SaveReadingProgressOut 保存阅读进度输出
type SaveReadingProgressOut struct {
	Success bool `json:"success"`
}

// GetReadingHistoryInp 获取阅读历史输入
type GetReadingHistoryInp struct {
	UserId int64 `json:"user_id" v:"required#用户ID不能为空"`
	Page   int   `json:"page" v:"min:1#页码不能小于1"`
	Limit  int   `json:"limit" v:"min:1|max:100#每页数量不能小于1|每页数量不能大于100"`
}

// GetReadingHistoryOut 获取阅读历史输出
type GetReadingHistoryOut struct {
	List  []ReadingHistoryItem `json:"list"`
	Total int64                `json:"total"`
}

// ReadingHistoryItem 阅读历史项
type ReadingHistoryItem struct {
	ComicId      int64     `json:"comic_id"`
	ComicTitle   string    `json:"comic_title"`
	ComicCover   string    `json:"comic_cover"`
	ChapterId    int64     `json:"chapter_id"`
	ChapterTitle string    `json:"chapter_title"`
	PageIndex    int       `json:"page_index"`
	LastReadAt   time.Time `json:"last_read_at"`
	Progress     float64   `json:"progress"` // 阅读进度百分比
}

// AddBookmarkInp 添加书签输入
type AddBookmarkInp struct {
	UserId    int64  `json:"user_id" v:"required#用户ID不能为空"`
	ComicId   int64  `json:"comic_id" v:"required#漫画ID不能为空"`
	ChapterId int64  `json:"chapter_id" v:"required#章节ID不能为空"`
	PageIndex int    `json:"page_index" v:"required|min:0#页码不能为空|页码不能小于0"`
	Note      string `json:"note"` // 书签备注
}

// AddBookmarkOut 添加书签输出
type AddBookmarkOut struct {
	BookmarkId int64 `json:"bookmark_id"`
}

// RemoveBookmarkInp 删除书签输入
type RemoveBookmarkInp struct {
	UserId     int64 `json:"user_id" v:"required#用户ID不能为空"`
	BookmarkId int64 `json:"bookmark_id" v:"required#书签ID不能为空"`
}

// RemoveBookmarkOut 删除书签输出
type RemoveBookmarkOut struct {
	Success bool `json:"success"`
}

// GetBookmarksInp 获取书签列表输入
type GetBookmarksInp struct {
	UserId  int64 `json:"user_id" v:"required#用户ID不能为空"`
	ComicId int64 `json:"comic_id"` // 可选，指定漫画的书签
	Page    int   `json:"page" v:"min:1#页码不能小于1"`
	Limit   int   `json:"limit" v:"min:1|max:100#每页数量不能小于1|每页数量不能大于100"`
}

// GetBookmarksOut 获取书签列表输出
type GetBookmarksOut struct {
	List  []BookmarkItem `json:"list"`
	Total int64          `json:"total"`
}

// BookmarkItem 书签项
type BookmarkItem struct {
	BookmarkId   int64     `json:"bookmark_id"`
	ComicId      int64     `json:"comic_id"`
	ComicTitle   string    `json:"comic_title"`
	ComicCover   string    `json:"comic_cover"`
	ChapterId    int64     `json:"chapter_id"`
	ChapterTitle string    `json:"chapter_title"`
	PageIndex    int       `json:"page_index"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"created_at"`
}

// UpdateReadingTimeInp 更新阅读时长输入
type UpdateReadingTimeInp struct {
	UserId      int64 `json:"user_id" v:"required#用户ID不能为空"`
	ComicId     int64 `json:"comic_id" v:"required#漫画ID不能为空"`
	ReadingTime int64 `json:"reading_time" v:"required|min:1#阅读时长不能为空|阅读时长不能小于1秒"` // 秒
}

// UpdateReadingTimeOut 更新阅读时长输出
type UpdateReadingTimeOut struct {
	Success bool `json:"success"`
}

// GetReadingStatsInp 获取阅读统计输入
type GetReadingStatsInp struct {
	UserId int64 `json:"user_id" v:"required#用户ID不能为空"`
}

// GetReadingStatsOut 获取阅读统计输出
type GetReadingStatsOut struct {
	TotalReadingTime int64   `json:"total_reading_time"` // 总阅读时长（秒）
	TotalChapters    int     `json:"total_chapters"`    // 总阅读章节数
	TotalPages       int     `json:"total_pages"`       // 总阅读页数
	ReadingStreak    int     `json:"reading_streak"`    // 连续阅读天数
	CompletionRate   float64 `json:"completion_rate"`   // 完成率
	LastReadAt       *time.Time `json:"last_read_at"`      // 最后阅读时间
	FavoriteGenres   []string   `json:"favorite_genres"`   // 喜欢的类型
	ReadingReport    ReadingReport `json:"reading_report"`  // 阅读报告
}

// ReadingReport 阅读报告
type ReadingReport struct {
	DailyAverage   float64            `json:"daily_average"`   // 日均阅读时长（分钟）
	WeeklyStats    []WeeklyReadingStat `json:"weekly_stats"`    // 周统计
	MonthlyStats   []MonthlyReadingStat `json:"monthly_stats"`   // 月统计
	TopComics      []TopComicStat      `json:"top_comics"`      // 最常读的漫画
	ReadingHabits  ReadingHabits       `json:"reading_habits"`  // 阅读习惯
}

// WeeklyReadingStat 周阅读统计
type WeeklyReadingStat struct {
	Week        string `json:"week"`        // 周（YYYY-WW）
	ReadingTime int64  `json:"reading_time"` // 阅读时长（秒）
	Chapters    int    `json:"chapters"`     // 章节数
	Pages       int    `json:"pages"`        // 页数
}

// MonthlyReadingStat 月阅读统计
type MonthlyReadingStat struct {
	Month       string `json:"month"`        // 月（YYYY-MM）
	ReadingTime int64  `json:"reading_time"` // 阅读时长（秒）
	Chapters    int    `json:"chapters"`     // 章节数
	Pages       int    `json:"pages"`        // 页数
}

// TopComicStat 热门漫画统计
type TopComicStat struct {
	ComicId     int64  `json:"comic_id"`
	ComicTitle  string `json:"comic_title"`
	ComicCover  string `json:"comic_cover"`
	ReadingTime int64  `json:"reading_time"` // 阅读时长（秒）
	Chapters    int    `json:"chapters"`     // 阅读章节数
	Pages       int    `json:"pages"`        // 阅读页数
}

// ReadingHabits 阅读习惯
type ReadingHabits struct {
	PreferredTime   string  `json:"preferred_time"`   // 偏好阅读时间段
	AverageSession  float64 `json:"average_session"`  // 平均阅读时长（分钟）
	MostActiveDay   string  `json:"most_active_day"`  // 最活跃的星期
	ReadingSpeed    float64 `json:"reading_speed"`    // 阅读速度（页/分钟）
}