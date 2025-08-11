// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysTwoFactorLog is the golang structure for table sys_two_factor_log.
type SysTwoFactorLog struct {
	Id        int64       `json:"id"        description:"主键ID"`
	MemberId  int64       `json:"memberId"  description:"用户ID"`
	Username  string      `json:"username"  description:"用户名"`
	Action    string      `json:"action"    description:"操作类型：enable,disable,verify,backup_used"`
	Method    string      `json:"method"    description:"2FA方法：totp,backup_code"`
	Result    string      `json:"result"    description:"操作结果：success,failed"`
	Ip        string      `json:"ip"        description:"IP地址"`
	UserAgent string      `json:"userAgent" description:"用户代理"`
	Remark    string      `json:"remark"    description:"备注"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"更新时间"`
}