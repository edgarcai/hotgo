// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

var (
	TwoFactorStats = cTwoFactorStats{}
)

// cTwoFactorStats 2FA统计控制器
type cTwoFactorStats struct{}

// GetStats 获取2FA统计数据
func (c *cTwoFactorStats) GetStats(ctx context.Context, req *sysin.TwoFactorStatsInp) (res g.Map, err error) {
	stats, err := service.SysTwoFactorStats().GetStats(ctx, req)
	if err != nil {
		return nil, err
	}

	return g.Map{
		"code":    0,
		"message": "获取统计数据成功",
		"data":    stats,
	}, nil
}

// GetTrend 获取2FA趋势数据
func (c *cTwoFactorStats) GetTrend(ctx context.Context, req *sysin.TwoFactorTrendInp) (res g.Map, err error) {
	trend, err := service.SysTwoFactorStats().GetTrend(ctx, req)
	if err != nil {
		return nil, err
	}

	return g.Map{
		"code":    0,
		"message": "获取趋势数据成功",
		"data":    trend,
	}, nil
}
