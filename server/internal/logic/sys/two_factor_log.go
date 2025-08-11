// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

type sSysTwoFactorLog struct{}

func NewSysTwoFactorLog() *sSysTwoFactorLog {
	return &sSysTwoFactorLog{}
}

func init() {
	service.RegisterSysTwoFactorLog(NewSysTwoFactorLog())
}

// Delete 删除2FA操作日志
func (s *sSysTwoFactorLog) Delete(ctx context.Context, in *sysin.TwoFactorLogDeleteInp) (err error) {
	if in.Id == nil {
		err = gerror.New("删除ID不能为空")
		return
	}

	// 使用系统日志表存储2FA操作日志，通过operation_type字段区分
	_, err = dao.SysLog.Ctx(ctx).Where(dao.SysLog.Columns().Id, in.Id).Delete()
	return
}

// View 获取2FA操作日志详情
func (s *sSysTwoFactorLog) View(ctx context.Context, in *sysin.TwoFactorLogViewInp) (res *sysin.TwoFactorLogViewModel, err error) {
	if in.Id <= 0 {
		err = gerror.New("ID不能为空")
		return
	}

	var models *entity.SysTwoFactorLog
	if err = dao.SysTwoFactorLog.Ctx(ctx).Where(dao.SysTwoFactorLog.Columns().Id, in.Id).Scan(&models); err != nil {
		return
	}

	if models == nil {
		err = gerror.New("数据不存在")
		return
	}

	res = new(sysin.TwoFactorLogViewModel)
	res.SysTwoFactorLog = *models

	return
}

// List 获取2FA操作日志列表
func (s *sSysTwoFactorLog) List(ctx context.Context, in *sysin.TwoFactorLogListInp) (list []*sysin.TwoFactorLogListModel, totalCount int, err error) {
	mod := dao.SysTwoFactorLog.Ctx(ctx)

	// 条件过滤
	if in.MemberId > 0 {
		mod = mod.Where(dao.SysTwoFactorLog.Columns().MemberId, in.MemberId)
	}

	if in.Username != "" {
		mod = mod.WhereLike(dao.SysTwoFactorLog.Columns().Username, "%"+in.Username+"%")
	}

	if in.Method != "" {
		mod = mod.Where(dao.SysTwoFactorLog.Columns().Method, in.Method)
	}

	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(dao.SysTwoFactorLog.Columns().CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	if in.Ip != "" {
		mod = mod.WhereLike(dao.SysTwoFactorLog.Columns().Ip, "%"+in.Ip+"%")
	}

	// 获取总数
	totalCount, err = mod.Count()
	if err != nil {
		err = gerror.Wrap(err, "获取2FA操作日志总数失败，请稍后重试！")
		return
	}

	if totalCount == 0 {
		return
	}

	// 分页查询
	if err = mod.Fields(sysin.TwoFactorLogListModel{}).Page(in.Page, in.PerPage).OrderDesc(dao.SysTwoFactorLog.Columns().Id).Scan(&list); err != nil {
		err = gerror.Wrap(err, "获取2FA操作日志列表失败，请稍后重试！")
		return
	}

	return
}

// Push 推送2FA操作日志
func (s *sSysTwoFactorLog) Push(ctx context.Context, in *sysin.TwoFactorLogPushInp) {
	if in.MemberId <= 0 {
		return
	}

	// 构造日志数据
	logData := &entity.SysTwoFactorLog{
		MemberId:  in.MemberId,
		Username:  in.Username,
		Action:    in.OperationType,
		Method:    in.Method,
		Result:    "success",
		Ip:        in.Ip,
		UserAgent: in.UserAgent,
		Remark:    in.Remark,
		CreatedAt: gtime.Now(),
		UpdatedAt: gtime.Now(),
	}

	// 根据状态设置结果
	if in.Status == 2 {
		logData.Result = "failed"
	}

	// 异步写入日志
	go func() {
		ctx := gctx.New()
		if _, err := dao.SysTwoFactorLog.Ctx(ctx).Insert(logData); err != nil {
			g.Log().Error(ctx, "写入2FA操作日志失败:", err)
		}
	}()
}
