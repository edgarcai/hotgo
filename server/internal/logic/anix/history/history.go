package history

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"hotgo/api/anix/history/v1"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/service"
)

type sHistory struct{}

func init() {
	service.RegisterHistory(New())
}

func New() *sHistory {
	return &sHistory{}
}

// SaveReadingProgress 保存阅读进度
func (s *sHistory) SaveReadingProgress(ctx context.Context, userId int64, req *v1.SaveReadingProgressReq) error {
	// 开启事务
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 检查是否已存在记录
		existingRecord := &entity.AnixReadingHistory{}
		err := dao.AnixReadingHistory.Ctx(ctx).TX(tx).Where(dao.AnixReadingHistory.Columns().UserId, userId).
			Where(dao.AnixReadingHistory.Columns().ComicId, req.ComicId).
			Where(dao.AnixReadingHistory.Columns().ChapterId, req.ChapterId).
			Scan(existingRecord)

		if err != nil && err != sql.ErrNoRows {
			return err
		}

		now := gtime.Now()
		if existingRecord.Id > 0 {
			// 更新现有记录
			updateData := g.Map{
				dao.AnixReadingHistory.Columns().PageIndex:   req.PageIndex,
				dao.AnixReadingHistory.Columns().Progress:    req.Progress,
				dao.AnixReadingHistory.Columns().ReadingTime: existingRecord.ReadingTime + req.ReadingTime,
				dao.AnixReadingHistory.Columns().LastReadAt:  now,
				dao.AnixReadingHistory.Columns().UpdatedAt:   now,
			}
			_, err = dao.AnixReadingHistory.Ctx(ctx).TX(tx).Where(dao.AnixReadingHistory.Columns().Id, existingRecord.Id).Update(updateData)
		} else {
			// 插入新记录
			insertData := &entity.AnixReadingHistory{
				UserId:      userId,
				ComicId:     req.ComicId,
				ChapterId:   req.ChapterId,
				PageIndex:   req.PageIndex,
				Progress:    req.Progress,
				ReadingTime: req.ReadingTime,
				LastReadAt:  now,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			_, err = dao.AnixReadingHistory.Ctx(ctx).TX(tx).Insert(insertData)
		}

		if err != nil {
			return err
		}

		// 更新阅读统计
		return s.updateReadingStats(ctx, tx, userId, req.ComicId, req.ReadingTime)
	})
}

// GetReadingHistory 获取阅读历史
func (s *sHistory) GetReadingHistory(ctx context.Context, userId int64, req *v1.GetReadingHistoryReq) (*v1.GetReadingHistoryRes, error) {
	// 计算偏移量
	offset := (req.Page - 1) * req.Limit

	// 查询阅读历史记录
	var records []entity.AnixReadingHistory
	err := dao.AnixReadingHistory.Ctx(ctx).Where(dao.AnixReadingHistory.Columns().UserId, userId).
		OrderDesc(dao.AnixReadingHistory.Columns().LastReadAt).
		Limit(req.Limit).
		Offset(offset).
		Scan(&records)

	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := dao.AnixReadingHistory.Ctx(ctx).Where(dao.AnixReadingHistory.Columns().UserId, userId).Count()
	if err != nil {
		return nil, err
	}

	// 构建返回数据
	list := make([]v1.ReadingHistoryItem, 0, len(records))
	for _, record := range records {
		// 获取漫画信息
		comic := &entity.AnixComic{}
		err = dao.AnixComic.Ctx(ctx).Where(dao.AnixComic.Columns().Id, record.ComicId).Scan(comic)
		if err != nil {
			continue
		}

		// 获取章节信息
		chapter := &entity.AnixChapter{}
		err = dao.AnixChapter.Ctx(ctx).Where(dao.AnixChapter.Columns().Id, record.ChapterId).Scan(chapter)
		if err != nil {
			continue
		}

		item := v1.ReadingHistoryItem{
			Id:           record.Id,
			ComicId:      record.ComicId,
			ComicTitle:   comic.Title,
			ComicCover:   comic.Cover,
			ChapterId:    record.ChapterId,
			ChapterTitle: chapter.Title,
			PageIndex:    record.PageIndex,
			Progress:     record.Progress,
			ReadingTime:  record.ReadingTime,
			LastReadAt:   record.LastReadAt,
		}
		list = append(list, item)
	}

	return &v1.GetReadingHistoryRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// GetReadingProgress 获取阅读进度
func (s *sHistory) GetReadingProgress(ctx context.Context, userId int64, comicId, chapterId int64) (*v1.GetReadingProgressRes, error) {
	record := &entity.AnixReadingHistory{}
	err := dao.AnixReadingHistory.Ctx(ctx).Where(dao.AnixReadingHistory.Columns().UserId, userId).
		Where(dao.AnixReadingHistory.Columns().ComicId, comicId).
		Where(dao.AnixReadingHistory.Columns().ChapterId, chapterId).
		Scan(record)

	if err != nil {
		if err == sql.ErrNoRows {
			// 没有阅读记录，返回默认值
			return &v1.GetReadingProgressRes{
				PageIndex:   0,
				Progress:    0,
				ReadingTime: 0,
				LastReadAt:  nil,
			}, nil
		}
		return nil, err
	}

	return &v1.GetReadingProgressRes{
		PageIndex:   record.PageIndex,
		Progress:    record.Progress,
		ReadingTime: record.ReadingTime,
		LastReadAt:  record.LastReadAt,
	}, nil
}

// AddBookmark 添加书签
func (s *sHistory) AddBookmark(ctx context.Context, userId int64, req *v1.AddBookmarkReq) (int64, error) {
	now := gtime.Now()
	bookmark := &entity.AnixBookmarks{
		UserId:    userId,
		ComicId:   req.ComicId,
		ChapterId: req.ChapterId,
		PageIndex: req.PageIndex,
		Note:      req.Note,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := dao.AnixBookmarks.Ctx(ctx).Insert(bookmark)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetBookmarks 获取书签列表
func (s *sHistory) GetBookmarks(ctx context.Context, userId int64, req *v1.GetBookmarksReq) (*v1.GetBookmarksRes, error) {
	// 计算偏移量
	offset := (req.Page - 1) * req.Limit

	// 构建查询条件
	query := dao.AnixBookmarks.Ctx(ctx).Where(dao.AnixBookmarks.Columns().UserId, userId)
	if req.ComicId > 0 {
		query = query.Where(dao.AnixBookmarks.Columns().ComicId, req.ComicId)
	}

	// 查询书签记录
	var bookmarks []entity.AnixBookmarks
	err := query.OrderDesc(dao.AnixBookmarks.Columns().CreatedAt).
		Limit(req.Limit).
		Offset(offset).
		Scan(&bookmarks)

	if err != nil {
		return nil, err
	}

	// 获取总数
	total, err := query.Count()
	if err != nil {
		return nil, err
	}

	// 构建返回数据
	list := make([]v1.BookmarkItem, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		// 获取漫画信息
		comic := &entity.AnixComic{}
		err = dao.AnixComic.Ctx(ctx).Where(dao.AnixComic.Columns().Id, bookmark.ComicId).Scan(comic)
		if err != nil {
			continue
		}

		// 获取章节信息
		chapter := &entity.AnixChapter{}
		err = dao.AnixChapter.Ctx(ctx).Where(dao.AnixChapter.Columns().Id, bookmark.ChapterId).Scan(chapter)
		if err != nil {
			continue
		}

		item := v1.BookmarkItem{
			Id:           bookmark.Id,
			ComicId:      bookmark.ComicId,
			ComicTitle:   comic.Title,
			ComicCover:   comic.Cover,
			ChapterId:    bookmark.ChapterId,
			ChapterTitle: chapter.Title,
			PageIndex:    bookmark.PageIndex,
			Note:         bookmark.Note,
			CreatedAt:    bookmark.CreatedAt,
		}
		list = append(list, item)
	}

	return &v1.GetBookmarksRes{
		List:  list,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

// DeleteBookmark 删除书签
func (s *sHistory) DeleteBookmark(ctx context.Context, userId int64, bookmarkId int64) error {
	_, err := dao.AnixBookmarks.Ctx(ctx).Where(dao.AnixBookmarks.Columns().Id, bookmarkId).
		Where(dao.AnixBookmarks.Columns().UserId, userId).
		Delete()
	return err
}

// GetReadingStats 获取阅读统计
func (s *sHistory) GetReadingStats(ctx context.Context, userId int64, comicId int64) (*v1.GetReadingStatsRes, error) {
	if comicId > 0 {
		// 获取指定漫画的统计
		stats := &entity.AnixReadingStats{}
		err := dao.AnixReadingStats.Ctx(ctx).Where(dao.AnixReadingStats.Columns().UserId, userId).
			Where(dao.AnixReadingStats.Columns().ComicId, comicId).
			Scan(stats)

		if err != nil {
			if err == sql.ErrNoRows {
				return &v1.GetReadingStatsRes{
					Stats:            []v1.ReadingStatsItem{},
					TotalReadingTime: 0,
					TotalComics:      0,
				}, nil
			}
			return nil, err
		}

		// 获取漫画信息
		comic := &entity.AnixComic{}
		err = dao.AnixComic.Ctx(ctx).Where(dao.AnixComic.Columns().Id, stats.ComicId).Scan(comic)
		if err != nil {
			return nil, err
		}

		item := v1.ReadingStatsItem{
			ComicId:          stats.ComicId,
			ComicTitle:       comic.Title,
			ComicCover:       comic.Cover,
			TotalChapters:    stats.TotalChapters,
			ReadChapters:     stats.ReadChapters,
			TotalReadingTime: stats.TotalReadingTime,
			FirstReadAt:      stats.FirstReadAt,
			LastReadAt:       stats.LastReadAt,
		}

		return &v1.GetReadingStatsRes{
			Stats:            []v1.ReadingStatsItem{item},
			TotalReadingTime: stats.TotalReadingTime,
			TotalComics:      1,
		}, nil
	} else {
		// 获取所有漫画的统计
		var statsList []entity.AnixReadingStats
		err := dao.AnixReadingStats.Ctx(ctx).Where(dao.AnixReadingStats.Columns().UserId, userId).
			OrderDesc(dao.AnixReadingStats.Columns().LastReadAt).
			Scan(&statsList)

		if err != nil {
			return nil, err
		}

		list := make([]v1.ReadingStatsItem, 0, len(statsList))
		var totalReadingTime int64
		for _, stats := range statsList {
			// 获取漫画信息
			comic := &entity.AnixComic{}
			err = dao.AnixComic.Ctx(ctx).Where(dao.AnixComic.Columns().Id, stats.ComicId).Scan(comic)
			if err != nil {
				continue
			}

			item := v1.ReadingStatsItem{
				ComicId:          stats.ComicId,
				ComicTitle:       comic.Title,
				ComicCover:       comic.Cover,
				TotalChapters:    stats.TotalChapters,
				ReadChapters:     stats.ReadChapters,
				TotalReadingTime: stats.TotalReadingTime,
				FirstReadAt:      stats.FirstReadAt,
				LastReadAt:       stats.LastReadAt,
			}
			list = append(list, item)
			totalReadingTime += stats.TotalReadingTime
		}

		return &v1.GetReadingStatsRes{
			Stats:            list,
			TotalReadingTime: totalReadingTime,
			TotalComics:      len(list),
		}, nil
	}
}

// updateReadingStats 更新阅读统计
func (s *sHistory) updateReadingStats(ctx context.Context, tx gdb.TX, userId, comicId int64, readingTime int64) error {
	// 检查是否已存在统计记录
	existingStats := &entity.AnixReadingStats{}
	err := dao.AnixReadingStats.Ctx(ctx).TX(tx).Where(dao.AnixReadingStats.Columns().UserId, userId).
		Where(dao.AnixReadingStats.Columns().ComicId, comicId).
		Scan(existingStats)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// 获取漫画总章节数
	totalChapters, err := dao.AnixChapter.Ctx(ctx).TX(tx).Where(dao.AnixChapter.Columns().ComicId, comicId).Count()
	if err != nil {
		return err
	}

	// 获取用户已读章节数
	readChapters, err := dao.AnixReadingHistory.Ctx(ctx).TX(tx).Where(dao.AnixReadingHistory.Columns().UserId, userId).
		Where(dao.AnixReadingHistory.Columns().ComicId, comicId).
		Count()
	if err != nil {
		return err
	}

	now := gtime.Now()
	if existingStats.Id > 0 {
		// 更新现有统计
		updateData := g.Map{
			dao.AnixReadingStats.Columns().TotalChapters:    totalChapters,
			dao.AnixReadingStats.Columns().ReadChapters:     readChapters,
			dao.AnixReadingStats.Columns().TotalReadingTime: existingStats.TotalReadingTime + readingTime,
			dao.AnixReadingStats.Columns().LastReadAt:       now,
			dao.AnixReadingStats.Columns().UpdatedAt:        now,
		}
		_, err = dao.AnixReadingStats.Ctx(ctx).TX(tx).Where(dao.AnixReadingStats.Columns().Id, existingStats.Id).Update(updateData)
	} else {
		// 插入新统计记录
		insertData := &entity.AnixReadingStats{
			UserId:           userId,
			ComicId:          comicId,
			TotalChapters:    int(totalChapters),
			ReadChapters:     int(readChapters),
			TotalReadingTime: readingTime,
			FirstReadAt:      now,
			LastReadAt:       now,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		_, err = dao.AnixReadingStats.Ctx(ctx).TX(tx).Insert(insertData)
	}

	return err
}