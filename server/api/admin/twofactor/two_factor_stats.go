// Package twofactor
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package twofactor

import (
	"hotgo/internal/model/input/sysin"

	"github.com/gogf/gf/v2/frame/g"
)

// TwoFactorStatsGetStatsReq 获取2FA统计数据请求
type TwoFactorStatsGetStatsReq struct {
	g.Meta `path:"/twofactor/stats/get" method:"get" tags:"2FA统计" summary:"获取2FA统计数据"`
	sysin.TwoFactorStatsInp
}

type TwoFactorStatsGetStatsRes struct {
	*sysin.TwoFactorStatsModel
}

// TwoFactorStatsGetTrendReq 获取2FA趋势数据请求
type TwoFactorStatsGetTrendReq struct {
	g.Meta `path:"/twofactor/stats/trend" method:"get" tags:"2FA统计" summary:"获取2FA趋势数据"`
	sysin.TwoFactorTrendInp
}

type TwoFactorStatsGetTrendRes struct {
	List []*sysin.TwoFactorTrendModel `json:"list" dc:"趋势数据列表"`
}
