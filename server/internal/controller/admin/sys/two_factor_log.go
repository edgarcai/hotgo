// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/api/admin/twofactorlog"
	"hotgo/internal/service"
)

var (
	TwoFactorLog = cTwoFactorLog{}
)

// cTwoFactorLog 2FA操作日志控制器
type cTwoFactorLog struct{}

// Delete 删除2FA操作日志
func (c *cTwoFactorLog) Delete(ctx context.Context, req *twofactorlog.DeleteReq) (res *twofactorlog.DeleteRes, err error) {
	err = service.SysTwoFactorLog().Delete(ctx, &req.TwoFactorLogDeleteInp)
	return
}

// View 获取2FA操作日志详情
func (c *cTwoFactorLog) View(ctx context.Context, req *twofactorlog.ViewReq) (res *twofactorlog.ViewRes, err error) {
	data, err := service.SysTwoFactorLog().View(ctx, &req.TwoFactorLogViewInp)
	if err != nil {
		return
	}

	res = new(twofactorlog.ViewRes)
	res.TwoFactorLogViewModel = data
	return
}

// List 获取2FA操作日志列表
func (c *cTwoFactorLog) List(ctx context.Context, req *twofactorlog.ListReq) (res *twofactorlog.ListRes, err error) {
	list, totalCount, err := service.SysTwoFactorLog().List(ctx, &req.TwoFactorLogListInp)
	if err != nil {
		return
	}

	res = new(twofactorlog.ListRes)
	res.List = list
	res.PageRes.Pack(req, totalCount)
	return
}
