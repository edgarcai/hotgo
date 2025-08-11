// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
)

// TwoFactorLogDeleteInp 删除2FA日志
type TwoFactorLogDeleteInp struct {
	Id interface{} `json:"id" v:"required#2FA日志ID不能为空" dc:"2FA日志ID"`
}

func (in *TwoFactorLogDeleteInp) Filter(ctx context.Context) (err error) {
	return
}

type TwoFactorLogDeleteModel struct{}

// TwoFactorLogViewInp 获取2FA日志详情
type TwoFactorLogViewInp struct {
	Id int64 `json:"id" v:"required#2FA日志ID不能为空" dc:"2FA日志ID"`
}

func (in *TwoFactorLogViewInp) Filter(ctx context.Context) (err error) {
	if in.Id <= 0 {
		err = gerror.New("2FA日志ID不能为空")
		return
	}
	return
}

type TwoFactorLogViewModel struct {
	entity.SysTwoFactorLog
}

// TwoFactorLogListInp 获取2FA日志列表
type TwoFactorLogListInp struct {
	form.PageReq
	MemberId     int64         `json:"memberId"     dc:"用户ID"`
	Username     string        `json:"username"     dc:"用户名"`
	OperationType string        `json:"operationType" dc:"操作类型"`
	Method       string        `json:"method"       dc:"验证方式"`
	Status       int           `json:"status"       dc:"状态"`
	CreatedAt    []*gtime.Time `json:"createdAt"    dc:"创建时间"`
	Ip           string        `json:"ip"           dc:"IP地址"`
}

func (in *TwoFactorLogListInp) Filter(ctx context.Context) (err error) {
	return
}

type TwoFactorLogListModel struct {
	entity.SysTwoFactorLog
	Username string `json:"username" dc:"用户名"`
}

// TwoFactorLogPushInp 推送2FA日志
type TwoFactorLogPushInp struct {
	MemberId      int64  `json:"memberId"      dc:"用户ID"`
	Username      string `json:"username"      dc:"用户名"`
	OperationType string `json:"operationType" dc:"操作类型：enable,disable,reset,verify"`
	Method        string `json:"method"        dc:"验证方式：totp,backup_code"`
	Status        int    `json:"status"        dc:"状态：1成功,2失败"`
	Remark        string `json:"remark"        dc:"备注"`
	Ip            string `json:"ip"            dc:"IP地址"`
	UserAgent     string `json:"userAgent"     dc:"用户代理"`
}

func (in *TwoFactorLogPushInp) Filter(ctx context.Context) (err error) {
	if in.MemberId <= 0 {
		err = gerror.New("用户ID不能为空")
		return
	}
	if in.OperationType == "" {
		err = gerror.New("操作类型不能为空")
		return
	}
	return
}

type TwoFactorLogPushModel struct{}