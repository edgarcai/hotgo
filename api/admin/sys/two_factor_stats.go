// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

// TwoFactorStatsGetStatsReq 获取2FA统计数据请求
type TwoFactorStatsGetStatsReq struct {
	DateRange []string `json:"dateRange" dc:"日期范围"`
}

// TwoFactorStatsModel 2FA统计数据模型
type TwoFactorStatsModel struct {
	TotalUsers       int64   `json:"totalUsers" dc:"总用户数"`
	EnabledUsers     int64   `json:"enabledUsers" dc:"已启用2FA用户数"`
	DisabledUsers    int64   `json:"disabledUsers" dc:"未启用2FA用户数"`
	EnableRate       float64 `json:"enableRate" dc:"启用率(%)"`
	TodayOperations  int64   `json:"todayOperations" dc:"今日操作次数"`
	WeekOperations   int64   `json:"weekOperations" dc:"本周操作次数"`
	MonthOperations  int64   `json:"monthOperations" dc:"本月操作次数"`
	SuccessfulLogins int64   `json:"successfulLogins" dc:"成功登录次数"`
	FailedLogins     int64   `json:"failedLogins" dc:"失败登录次数"`
	SuccessRate      float64 `json:"successRate" dc:"成功率(%)"`
}

// TwoFactorStatsGetStatsRes 获取2FA统计数据响应
type TwoFactorStatsGetStatsRes struct {
	Data *TwoFactorStatsModel `json:"data" dc:"统计数据"`
}

// TwoFactorStatsGetTrendReq 获取2FA趋势数据请求
type TwoFactorStatsGetTrendReq struct {
	Type      string   `json:"type" dc:"趋势类型:daily=日,weekly=周,monthly=月"`
	DateRange []string `json:"dateRange" dc:"日期范围"`
}

// TwoFactorTrendModel 2FA趋势数据模型
type TwoFactorTrendModel struct {
	Date       string `json:"date" dc:"日期"`
	Operations int64  `json:"operations" dc:"操作次数"`
	Successful int64  `json:"successful" dc:"成功次数"`
	Failed     int64  `json:"failed" dc:"失败次数"`
}

// TwoFactorStatsGetTrendRes 获取2FA趋势数据响应
type TwoFactorStatsGetTrendRes struct {
	Data []*TwoFactorTrendModel `json:"data" dc:"趋势数据"`
}
