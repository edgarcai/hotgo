// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/internal/dao"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSysTwoFactorStats struct{}

func NewSysTwoFactorStats() *sSysTwoFactorStats {
	return &sSysTwoFactorStats{}
}

func init() {
	service.RegisterSysTwoFactorStats(NewSysTwoFactorStats())
}

// GetStats 获取2FA统计数据
func (s *sSysTwoFactorStats) GetStats(ctx context.Context, in *sysin.TwoFactorStatsInp) (res *sysin.TwoFactorStatsModel, err error) {
	res = &sysin.TwoFactorStatsModel{}

	// 获取总用户数
	totalUsers, err := dao.AdminMember.Ctx(ctx).Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取总用户数失败")
	}
	res.TotalUsers = int64(totalUsers)

	// 获取已启用2FA的用户数
	enabledUsers, err := dao.AdminTwoFactor.Ctx(ctx).Where("is_enabled = 1").Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取已启用2FA用户数失败")
	}
	res.EnabledUsers = int64(enabledUsers)

	// 计算未启用2FA的用户数
	res.DisabledUsers = res.TotalUsers - res.EnabledUsers

	// 计算启用率
	if res.TotalUsers > 0 {
		res.EnableRate = float64(res.EnabledUsers) / float64(res.TotalUsers) * 100
	}

	// 获取今日操作次数
	todayStart := gtime.Now().StartOfDay()
	todayEnd := gtime.Now().EndOfDay()
	todayOps, err := dao.SysLog.Ctx(ctx).
		Where("JSON_EXTRACT(post_data, '$.operation_type') IN ('2fa_enable', '2fa_disable', '2fa_reset', '2fa_verify')").
		WhereBetween(dao.SysLog.Columns().CreatedAt, todayStart, todayEnd).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取今日操作次数失败")
	}
	res.TodayOperations = int64(todayOps)

	// 获取本周操作次数
	weekStart := gtime.Now().StartOfWeek()
	weekEnd := gtime.Now().EndOfWeek()
	weekOps, err := dao.SysLog.Ctx(ctx).
		Where("JSON_EXTRACT(post_data, '$.operation_type') IN ('2fa_enable', '2fa_disable', '2fa_reset', '2fa_verify')").
		WhereBetween(dao.SysLog.Columns().CreatedAt, weekStart, weekEnd).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取本周操作次数失败")
	}
	res.WeekOperations = int64(weekOps)

	// 获取本月操作次数
	monthStart := gtime.Now().StartOfMonth()
	monthEnd := gtime.Now().EndOfMonth()
	monthOps, err := dao.SysLog.Ctx(ctx).
		Where("JSON_EXTRACT(post_data, '$.operation_type') IN ('2fa_enable', '2fa_disable', '2fa_reset', '2fa_verify')").
		WhereBetween(dao.SysLog.Columns().CreatedAt, monthStart, monthEnd).
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取本月操作次数失败")
	}
	res.MonthOperations = int64(monthOps)

	// 获取成功登录次数（2FA验证成功）
	successfulLogins, err := dao.SysLog.Ctx(ctx).
		Where("JSON_EXTRACT(post_data, '$.operation_type') = '2fa_verify_success'").
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取成功登录次数失败")
	}
	res.SuccessfulLogins = int64(successfulLogins)

	// 获取失败登录次数（2FA验证失败）
	failedLogins, err := dao.SysLog.Ctx(ctx).
		Where("JSON_EXTRACT(post_data, '$.operation_type') = '2fa_verify_failed'").
		Count()
	if err != nil {
		return nil, gerror.Wrap(err, "获取失败登录次数失败")
	}
	res.FailedLogins = int64(failedLogins)

	// 计算成功率
	totalAttempts := res.SuccessfulLogins + res.FailedLogins
	if totalAttempts > 0 {
		res.SuccessRate = float64(res.SuccessfulLogins) / float64(totalAttempts) * 100
	}

	return res, nil
}

// GetTrend 获取2FA趋势数据
func (s *sSysTwoFactorStats) GetTrend(ctx context.Context, in *sysin.TwoFactorTrendInp) (list []*sysin.TwoFactorTrendModel, err error) {
	list = make([]*sysin.TwoFactorTrendModel, 0)

	// 根据类型确定时间范围和格式
	var (
		startTime *gtime.Time
		endTime   *gtime.Time
		dayCount  int
	)

	switch in.Type {
	case "daily":
		startTime = gtime.Now().AddDate(0, 0, -30) // 最近30天
		endTime = gtime.Now()
		dayCount = 30
	case "weekly":
		startTime = gtime.Now().AddDate(0, 0, -84) // 最近12周
		endTime = gtime.Now()
		dayCount = 12
	case "monthly":
		startTime = gtime.Now().AddDate(0, -12, 0) // 最近12个月
		endTime = gtime.Now()
		dayCount = 12
	default:
		return nil, gerror.New("不支持的趋势类型")
	}

	// 构建SQL查询格式
	var sqlFormat string
	switch in.Type {
	case "daily":
		sqlFormat = "DATE_FORMAT(created_at, '%Y-%m-%d')"
	case "weekly":
		sqlFormat = "DATE_FORMAT(created_at, '%Y-W%u')"
	case "monthly":
		sqlFormat = "DATE_FORMAT(created_at, '%Y-%m')"
	}

	// 查询操作统计
	operationQuery := dao.SysLog.Ctx(ctx).
		Fields(g.Map{
			"date":       sqlFormat,
			"operations": "COUNT(*)",
		}).
		Where("JSON_EXTRACT(post_data, '$.operation_type') IN ('2fa_enable', '2fa_disable', '2fa_reset', '2fa_verify')").
		WhereBetween(dao.SysLog.Columns().CreatedAt, startTime, endTime).
		Group("date").
		Order("date ASC")

	var operationRecords []struct {
		Date       string `json:"date"`
		Operations int64  `json:"operations"`
	}

	err = operationQuery.Scan(&operationRecords)
	if err != nil {
		return nil, gerror.Wrap(err, "查询操作统计失败")
	}

	// 查询成功统计
	successQuery := dao.SysLog.Ctx(ctx).
		Fields(g.Map{
			"date":    sqlFormat,
			"success": "COUNT(*)",
		}).
		Where("JSON_EXTRACT(post_data, '$.operation_type') = '2fa_verify_success'").
		WhereBetween(dao.SysLog.Columns().CreatedAt, startTime, endTime).
		Group("date").
		Order("date ASC")

	var successRecords []struct {
		Date    string `json:"date"`
		Success int64  `json:"success"`
	}

	err = successQuery.Scan(&successRecords)
	if err != nil {
		return nil, gerror.Wrap(err, "查询成功统计失败")
	}

	// 查询失败统计
	failedQuery := dao.SysLog.Ctx(ctx).
		Fields(g.Map{
			"date":   sqlFormat,
			"failed": "COUNT(*)",
		}).
		Where("JSON_EXTRACT(post_data, '$.operation_type') = '2fa_verify_failed'").
		WhereBetween(dao.SysLog.Columns().CreatedAt, startTime, endTime).
		Group("date").
		Order("date ASC")

	var failedRecords []struct {
		Date   string `json:"date"`
		Failed int64  `json:"failed"`
	}

	err = failedQuery.Scan(&failedRecords)
	if err != nil {
		return nil, gerror.Wrap(err, "查询失败统计失败")
	}

	// 创建数据映射
	operationMap := make(map[string]int64)
	successMap := make(map[string]int64)
	failedMap := make(map[string]int64)

	for _, record := range operationRecords {
		operationMap[record.Date] = record.Operations
	}
	for _, record := range successRecords {
		successMap[record.Date] = record.Success
	}
	for _, record := range failedRecords {
		failedMap[record.Date] = record.Failed
	}

	// 生成完整的时间序列数据
	currentTime := startTime
	for i := 0; i < dayCount; i++ {
		var dateKey string
		switch in.Type {
		case "daily":
			dateKey = currentTime.Format("2006-01-02")
			currentTime = currentTime.AddDate(0, 0, 1)
		case "weekly":
			_, week := currentTime.ISOWeek()
			if week < 10 {
				dateKey = currentTime.Format("2006-W0") + g.NewVar(week).String()
			} else {
				dateKey = currentTime.Format("2006-W") + g.NewVar(week).String()
			}
			currentTime = currentTime.AddDate(0, 0, 7)
		case "monthly":
			dateKey = currentTime.Format("2006-01")
			currentTime = currentTime.AddDate(0, 1, 0)
		}

		model := &sysin.TwoFactorTrendModel{
			Date:       dateKey,
			Operations: operationMap[dateKey],
			Successful: successMap[dateKey],
			Failed:     failedMap[dateKey],
		}

		list = append(list, model)
	}

	return list, nil
}