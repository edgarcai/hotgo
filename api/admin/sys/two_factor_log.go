// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

// TwoFactorLogDeleteReq 删除2FA操作日志请求
type TwoFactorLogDeleteReq struct {
	Id int64 `json:"id" dc:"日志ID"`
}

// TwoFactorLogDeleteRes 删除2FA操作日志响应
type TwoFactorLogDeleteRes struct {
	Code    int    `json:"code" dc:"状态码"`
	Message string `json:"message" dc:"响应消息"`
}

// TwoFactorLogViewReq 获取2FA操作日志详情请求
type TwoFactorLogViewReq struct {
	Id int64 `json:"id" dc:"日志ID"`
}

// TwoFactorLogViewModel 2FA操作日志详情模型
type TwoFactorLogViewModel struct {
	Id            int64  `json:"id" dc:"日志ID"`
	MemberId      int64  `json:"memberId" dc:"用户ID"`
	Username      string `json:"username" dc:"用户名"`
	OperationType string `json:"operationType" dc:"操作类型"`
	OperationDesc string `json:"operationDesc" dc:"操作描述"`
	Ip            string `json:"ip" dc:"IP地址"`
	UserAgent     string `json:"userAgent" dc:"用户代理"`
	Details       string `json:"details" dc:"详细信息"`
	CreatedAt     string `json:"createdAt" dc:"创建时间"`
}

// TwoFactorLogViewRes 查看2FA操作日志响应
type TwoFactorLogViewRes struct {
	Data *TwoFactorLogViewModel `json:"data" dc:"日志详情"`
}

// TwoFactorLogListReq 获取2FA操作日志列表请求
type TwoFactorLogListReq struct {
	MemberId      int64  `json:"memberId" dc:"用户ID"`
	Username      string `json:"username" dc:"用户名"`
	OperationType string `json:"operationType" dc:"操作类型"`
	Ip            string `json:"ip" dc:"IP地址"`
	CreatedAt     string `json:"createdAt" dc:"创建时间"`
	Page          int    `json:"page" dc:"页码"`
	PerPage       int    `json:"perPage" dc:"每页数量"`
}

// TwoFactorLogListModel 2FA操作日志列表模型
type TwoFactorLogListModel struct {
	Id            int64  `json:"id" dc:"日志ID"`
	MemberId      int64  `json:"memberId" dc:"用户ID"`
	Username      string `json:"username" dc:"用户名"`
	OperationType string `json:"operationType" dc:"操作类型"`
	OperationDesc string `json:"operationDesc" dc:"操作描述"`
	Ip            string `json:"ip" dc:"IP地址"`
	CreatedAt     string `json:"createdAt" dc:"创建时间"`
}

// TwoFactorLogListRes 2FA操作日志列表响应
type TwoFactorLogListRes struct {
	List  []*TwoFactorLogListModel `json:"list" dc:"日志列表"`
	Page  int                      `json:"page" dc:"页码"`
	Count int64                    `json:"count" dc:"总数"`
}
