// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Admin2faLog is the golang structure for table admin_2fa_log.
type Admin2faLog struct {
	Id        int64       `json:"id"        orm:"id"         description:"主键ID"`
	AdminId   int64       `json:"adminId"   orm:"admin_id"   description:"管理员ID"`
	Action    string      `json:"action"    orm:"action"     description:"操作类型(setup:设置 verify:验证 disable:禁用 backup_used:使用备用码)"`
	Result    string      `json:"result"    orm:"result"     description:"操作结果(success:成功 failed:失败)"`
	Ip        string      `json:"ip"        orm:"ip"         description:"IP地址"`
	UserAgent string      `json:"userAgent" orm:"user_agent" description:"用户代理"`
	Remark    string      `json:"remark"    orm:"remark"     description:"备注信息"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}