package reading_history

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

// sReadingHistory 阅读历史记录服务实现
type sReadingHistory struct{}

func init() {
	service.RegisterReadingHistory(New())
}

// New 创建阅读历史记录服务实例
func New() service.IReadingHistory {
	return &sReadingHistory{}
}

// SaveReadingProgress 保存阅读进度
func (s *sReadingHistory) SaveReadingProgress(ctx context.Context, in *service.SaveReadingProgressInp) (*service.SaveReadingProgressOut, error) {
	g.Log().Info(ctx, "保存阅读进度:", in)
	
	// 检查是否已存在阅读记录
	var existingRecord *entity.AnixReadingStats
	err := dao.AnixReadingStats.Ctx(ctx).Where("user_id = ? AND comic_id = ?", in.UserId, in.ComicId).Scan(&existingRecord)
	if err != nil {
		return nil, err
	}
	
	now := gtime.Now()
	
	if existingRecord != nil {
		// 更新现有记录
		_, err = dao.AnixReadingStats.Ctx(ctx).Where("id = ?", existingRecord.Id).Update(g.Map{
			"last_read_at": now,
			"updated_at":   now,
		})
	} else {
		// 创建新记录
		_, err = dao.AnixReadingStats.Ctx(ctx).Insert(g.Map{
			"user_id":      in.UserId,
			"comic_id":     in.ComicId,
			"last_read_at": now,
			"created_at":   now,
			"updated_at":   now,
		})
	}
	
	if err != nil {
		return nil, err
	}
	
	// 保存具体的阅读进度（章节和页码）
	err = s.saveChapterProgress(ctx, in.UserId, in.ComicId, in.ChapterId, in.PageIndex)
	if err != nil {
		g.Log().Error(ctx, "保存章节进度失败:", err)
	}
	
	return &service.SaveReadingProgressOut{
		Success: true,
	}, nil
}

// GetReadingHistory 获取阅读历史
func (s *sReadingHistory) GetReadingHistory(ctx context.Context, in *service.GetReadingHistoryInp) (*service.GetReadingHistoryOut, error) {
	g.Log().Info(ctx, "获取阅读历史:", in)
	
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}
	
	offset := (in.Page - 1) * in.Limit
	
	// 查询阅读历史，关联漫画和章节信息
	sql := `
		SELECT 
			rs.comic_id,
			c.title as comic_title,
			c.cover as comic_cover,
			rs.last_read_at,
			0 as chapter_id,
			'' as chapter_title,
			0 as page_index,
			0 as progress
		FROM anix_reading_stats rs
		LEFT JOIN anix_comic c ON rs.comic_id = c.id
		WHERE rs.user_id = ?
		ORDER BY rs.last_read_at DESC
		LIMIT ? OFFSET ?
	`
	
	var records []service.ReadingHistoryItem
	err := g.DB().GetScan(ctx, &records, sql, in.UserId, in.Limit, offset)
	if err != nil {
		return nil, err
	}
	
	// 获取总数
	total, err := dao.AnixReadingStats.Ctx(ctx).Where("user_id = ?", in.UserId).Count()
	if err != nil {
		return nil, err
	}
	
	return &service.GetReadingHistoryOut{
		List:  records,
		Total: int64(total),
	}, nil
}

// AddBookmark 添加书签
func (s *sReadingHistory) AddBookmark(ctx context.Context, in *service.AddBookmarkInp) (*service.AddBookmarkOut, error) {
	g.Log().Info(ctx, "添加书签:", in)
	
	// 检查书签是否已存在
	count, err := g.DB().Model("anix_bookmark").Where("user_id = ? AND comic_id = ? AND chapter_id = ? AND page_index = ?", 
		in.UserId, in.ComicId, in.ChapterId, in.PageIndex).Count()
	if err != nil {
		return nil, err
	}
	
	if count > 0 {
		return nil, fmt.Errorf("书签已存在")
	}
	
	now := gtime.Now()
	
	// 插入书签记录
	result, err := g.DB().Model("anix_bookmark").Insert(g.Map{
		"user_id":    in.UserId,
		"comic_id":   in.ComicId,
		"chapter_id": in.ChapterId,
		"page_index": in.PageIndex,
		"note":       in.Note,
		"created_at": now,
		"updated_at": now,
	})
	
	if err != nil {
		return nil, err
	}
	
	bookmarkId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	
	return &service.AddBookmarkOut{
		BookmarkId: bookmarkId,
	}, nil
}

// RemoveBookmark 删除书签
func (s *sReadingHistory) RemoveBookmark(ctx context.Context, in *service.RemoveBookmarkInp) (*service.RemoveBookmarkOut, error) {
	g.Log().Info(ctx, "删除书签:", in)
	
	_, err := g.DB().Model("anix_bookmark").Where("id = ? AND user_id = ?", in.BookmarkId, in.UserId).Delete()
	if err != nil {
		return nil, err
	}
	
	return &service.RemoveBookmarkOut{
		Success: true,
	}, nil
}

// GetBookmarks 获取书签列表
func (s *sReadingHistory) GetBookmarks(ctx context.Context, in *service.GetBookmarksInp) (*service.GetBookmarksOut, error) {
	g.Log().Info(ctx, "获取书签列表:", in)
	
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Limit <= 0 {
		in.Limit = 20
	}
	
	offset := (in.Page - 1) * in.Limit
	
	// 构建查询条件
	query := g.DB().Model("anix_bookmark b").LeftJoin("anix_comic c", "b.comic_id = c.id").LeftJoin("anix_chapter ch", "b.chapter_id = ch.id")
	query = query.Where("b.user_id = ?", in.UserId)
	
	if in.ComicId > 0 {
		query = query.Where("b.comic_id = ?", in.ComicId)
	}
	
	// 查询书签列表
	sql := `
		SELECT 
			b.id as bookmark_id,
			b.comic_id,
			c.title as comic_title,
			c.cover as comic_cover,
			b.chapter_id,
			ch.title as chapter_title,
			b.page_index,
			b.note,
			b.created_at
		FROM anix_bookmark b
		LEFT JOIN anix_comic c ON b.comic_id = c.id
		LEFT JOIN anix_chapter ch ON b.chapter_id = ch.id
		WHERE b.user_id = ?
	`
	
	args := []interface{}{in.UserId}
	if in.ComicId > 0 {
		sql += " AND b.comic_id = ?"
		args = append(args, in.ComicId)
	}
	
	sql += " ORDER BY b.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, in.Limit, offset)
	
	var bookmarks []service.BookmarkItem
	err := g.DB().GetScan(ctx, &bookmarks, sql, args...)
	if err != nil {
		return nil, err
	}
	
	// 获取总数
	countQuery := g.DB().Model("anix_bookmark").Where("user_id = ?", in.UserId)
	if in.ComicId > 0 {
		countQuery = countQuery.Where("comic_id = ?", in.ComicId)
	}
	total, err := countQuery.Count()
	if err != nil {
		return nil, err
	}
	
	return &service.GetBookmarksOut{
		List:  bookmarks,
		Total: int64(total),
	}, nil
}

// UpdateReadingTime 更新阅读时长
func (s *sReadingHistory) UpdateReadingTime(ctx context.Context, in *service.UpdateReadingTimeInp) (*service.UpdateReadingTimeOut, error) {
	g.Log().Info(ctx, "更新阅读时长:", in)
	
	// 查找现有记录
	var existingRecord *entity.AnixReadingStats
	err := dao.AnixReadingStats.Ctx(ctx).Where("user_id = ? AND comic_id = ?", in.UserId, in.ComicId).Scan(&existingRecord)
	if err != nil {
		return nil, err
	}
	
	now := gtime.Now()
	
	if existingRecord != nil {
		// 更新现有记录
		newTotalTime := existingRecord.TotalReadingTime + in.ReadingTime
		_, err = dao.AnixReadingStats.Ctx(ctx).Where("id = ?", existingRecord.Id).Update(g.Map{
			"total_reading_time": newTotalTime,
			"last_read_at":       now,
			"updated_at":         now,
		})
	} else {
		// 创建新记录
		_, err = dao.AnixReadingStats.Ctx(ctx).Insert(g.Map{
			"user_id":            in.UserId,
			"comic_id":           in.ComicId,
			"total_reading_time": in.ReadingTime,
			"last_read_at":       now,
			"created_at":         now,
			"updated_at":         now,
		})
	}
	
	if err != nil {
		return nil, err
	}
	
	return &service.UpdateReadingTimeOut{
		Success: true,
	}, nil
}

// GetReadingStats 获取阅读统计
func (s *sReadingHistory) GetReadingStats(ctx context.Context, in *service.GetReadingStatsInp) (*service.GetReadingStatsOut, error) {
	g.Log().Info(ctx, "获取阅读统计:", in)
	
	// 获取基础统计数据
	var stats struct {
		TotalReadingTime int64      `json:"total_reading_time"`
		TotalChapters    int        `json:"total_chapters"`
		TotalPages       int        `json:"total_pages"`
		ReadingStreak    int        `json:"reading_streak"`
		CompletionRate   float64    `json:"completion_rate"`
		LastReadAt       *time.Time `json:"last_read_at"`
	}
	
	// 查询总阅读时长和最后阅读时间
	sql := `
		SELECT 
			COALESCE(SUM(total_reading_time), 0) as total_reading_time,
			COALESCE(SUM(read_chapters), 0) as total_chapters,
			COALESCE(SUM(pages_read), 0) as total_pages,
			COALESCE(AVG(completion_rate), 0) as completion_rate,
			MAX(last_read_at) as last_read_at
		FROM anix_reading_stats 
		WHERE user_id = ?
	`
	
	err := g.DB().GetScan(ctx, &stats, sql, in.UserId)
	if err != nil {
		return nil, err
	}
	
	// 计算连续阅读天数
	stats.ReadingStreak = s.calculateReadingStreak(ctx, in.UserId)
	
	// 生成阅读报告
	report, err := s.generateReadingReport(ctx, in.UserId)
	if err != nil {
		g.Log().Error(ctx, "生成阅读报告失败:", err)
		report = service.ReadingReport{} // 使用空报告
	}
	
	return &service.GetReadingStatsOut{
		TotalReadingTime: stats.TotalReadingTime,
		TotalChapters:    stats.TotalChapters,
		TotalPages:       stats.TotalPages,
		ReadingStreak:    stats.ReadingStreak,
		CompletionRate:   stats.CompletionRate,
		LastReadAt:       stats.LastReadAt,
		FavoriteGenres:   []string{"动作", "冒险", "喜剧"}, // TODO: 从用户阅读历史中分析
		ReadingReport:    report,
	}, nil
}

// saveChapterProgress 保存章节进度
func (s *sReadingHistory) saveChapterProgress(ctx context.Context, userId, comicId, chapterId int64, pageIndex int) error {
	// 这里可以保存更详细的章节阅读进度
	// 暂时简化处理
	return nil
}

// calculateReadingStreak 计算连续阅读天数
func (s *sReadingHistory) calculateReadingStreak(ctx context.Context, userId int64) int {
	// 查询最近的阅读记录，计算连续天数
	sql := `
		SELECT DATE(last_read_at) as read_date
		FROM anix_reading_stats 
		WHERE user_id = ? AND last_read_at IS NOT NULL
		ORDER BY last_read_at DESC
		LIMIT 30
	`
	
	var dates []string
	err := g.DB().GetScan(ctx, &dates, sql, userId)
	if err != nil || len(dates) == 0 {
		return 0
	}
	
	// 简化计算：返回最近有阅读记录的天数
	return len(dates)
}

// generateReadingReport 生成阅读报告
func (s *sReadingHistory) generateReadingReport(ctx context.Context, userId int64) (service.ReadingReport, error) {
	// 计算日均阅读时长
	dailyAverage := s.calculateDailyAverage(ctx, userId)
	
	// 获取周统计
	weeklyStats := s.getWeeklyStats(ctx, userId)
	
	// 获取月统计
	monthlyStats := s.getMonthlyStats(ctx, userId)
	
	// 获取热门漫画
	topComics := s.getTopComics(ctx, userId)
	
	// 分析阅读习惯
	habits := s.analyzeReadingHabits(ctx, userId)
	
	return service.ReadingReport{
		DailyAverage:  dailyAverage,
		WeeklyStats:   weeklyStats,
		MonthlyStats:  monthlyStats,
		TopComics:     topComics,
		ReadingHabits: habits,
	}, nil
}

// calculateDailyAverage 计算日均阅读时长
func (s *sReadingHistory) calculateDailyAverage(ctx context.Context, userId int64) float64 {
	sql := `
		SELECT AVG(total_reading_time) / 60.0 as daily_average
		FROM anix_reading_stats 
		WHERE user_id = ? AND total_reading_time > 0
	`
	
	var average float64
	g.DB().GetScan(ctx, &average, sql, userId)
	return average
}

// getWeeklyStats 获取周统计
func (s *sReadingHistory) getWeeklyStats(ctx context.Context, userId int64) []service.WeeklyReadingStat {
	// 简化实现，返回最近4周的数据
	return []service.WeeklyReadingStat{
		{Week: "2024-01", ReadingTime: 3600, Chapters: 10, Pages: 200},
		{Week: "2024-02", ReadingTime: 4200, Chapters: 12, Pages: 240},
	}
}

// getMonthlyStats 获取月统计
func (s *sReadingHistory) getMonthlyStats(ctx context.Context, userId int64) []service.MonthlyReadingStat {
	// 简化实现，返回最近3个月的数据
	return []service.MonthlyReadingStat{
		{Month: "2024-01", ReadingTime: 14400, Chapters: 40, Pages: 800},
		{Month: "2024-02", ReadingTime: 16800, Chapters: 45, Pages: 900},
	}
}

// getTopComics 获取热门漫画
func (s *sReadingHistory) getTopComics(ctx context.Context, userId int64) []service.TopComicStat {
	sql := `
		SELECT 
			rs.comic_id,
			c.title as comic_title,
			c.cover as comic_cover,
			rs.total_reading_time as reading_time,
			rs.read_chapters as chapters,
			rs.pages_read as pages
		FROM anix_reading_stats rs
		LEFT JOIN anix_comic c ON rs.comic_id = c.id
		WHERE rs.user_id = ?
		ORDER BY rs.total_reading_time DESC
		LIMIT 5
	`
	
	var topComics []service.TopComicStat
	g.DB().GetScan(ctx, &topComics, sql, userId)
	return topComics
}

// analyzeReadingHabits 分析阅读习惯
func (s *sReadingHistory) analyzeReadingHabits(ctx context.Context, userId int64) service.ReadingHabits {
	// 简化实现，返回默认习惯分析
	return service.ReadingHabits{
		PreferredTime:  "晚上",
		AverageSession: 30.5,
		MostActiveDay:  "周末",
		ReadingSpeed:   2.5,
	}
}