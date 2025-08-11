// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TwoFactorStatsInp 获取2FA统计数据
type TwoFactorStatsInp struct {
	DateRange []*gtime.Time `json:"dateRange" dc:"统计时间范围"`
}

type TwoFactorStatsModel struct {
	TotalUsers       int64   `json:"totalUsers"        dc:"总用户数"`
	EnabledUsers     int64   `json:"enabledUsers"      dc:"已启用2FA用户数"`
	DisabledUsers    int64   `json:"disabledUsers"     dc:"未启用2FA用户数"`
	EnableRate       float64 `json:"enableRate"        dc:"启用率（百分比）"`
	TodayOperations  int64   `json:"todayOperations"   dc:"今日2FA操作次数"`
	WeekOperations   int64   `json:"weekOperations"    dc:"本周2FA操作次数"`
	MonthOperations  int64   `json:"monthOperations"   dc:"本月2FA操作次数"`
	SuccessfulLogins int64   `json:"successfulLogins"  dc:"成功登录次数"`
	FailedLogins     int64   `json:"failedLogins"      dc:"失败登录次数"`
	SuccessRate      float64 `json:"successRate"       dc:"成功率（百分比）"`
}

// TwoFactorTrendInp 获取2FA趋势数据
type TwoFactorTrendInp struct {
	DateRange []*gtime.Time `json:"dateRange" dc:"统计时间范围"`
	Type      string        `json:"type"      dc:"趋势类型：daily,weekly,monthly"`
}

type TwoFactorTrendModel struct {
	Date       string `json:"date"        dc:"日期"`
	Operations int64  `json:"operations"  dc:"操作次数"`
	Successful int64  `json:"successful"  dc:"成功次数"`
	Failed     int64  `json:"failed"      dc:"失败次数"`
}
