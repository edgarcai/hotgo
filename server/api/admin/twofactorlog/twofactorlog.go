// Package twofactorlog
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package twofactorlog

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/internal/model/input/form"
	"hotgo/internal/model/input/sysin"
)

// ListReq 查询2FA操作日志列表
type ListReq struct {
	g.Meta `path:"/two-factor-log/list" method:"get" tags:"2FA操作日志" summary:"获取2FA操作日志列表"`
	sysin.TwoFactorLogListInp
}

type ListRes struct {
	form.PageRes
	List []*sysin.TwoFactorLogListModel `json:"list"   dc:"数据列表"`
}

// ViewReq 获取2FA操作日志详情
type ViewReq struct {
	g.Meta `path:"/two-factor-log/view" method:"get" tags:"2FA操作日志" summary:"获取2FA操作日志详情"`
	sysin.TwoFactorLogViewInp
}

type ViewRes struct {
	*sysin.TwoFactorLogViewModel
}

// DeleteReq 删除2FA操作日志
type DeleteReq struct {
	g.Meta `path:"/two-factor-log/delete" method:"post" tags:"2FA操作日志" summary:"删除2FA操作日志"`
	sysin.TwoFactorLogDeleteInp
}

type DeleteRes struct{}