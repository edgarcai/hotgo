// Package middleware
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package middleware

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/consts"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/response"
	"hotgo/internal/service"
	"hotgo/utility/simple"
)

// Admin2FA 后台2FA验证中间件
func (s *sMiddleware) Admin2FA(r *ghttp.Request) {
	var (
		ctx  = r.Context()
		path = gstr.Replace(r.URL.Path, simple.RouterPrefix(ctx, consts.AppAdmin), "", 1)
	)

	// 不需要验证登录的路由地址
	if s.IsExceptLogin(ctx, consts.AppAdmin, path) {
		r.Middleware.Next()
		return
	}

	// 2FA相关的路由不需要验证2FA
	if s.IsExcept2FA(ctx, path) {
		r.Middleware.Next()
		return
	}

	// 将用户信息传递到上下文中
	if err := s.DeliverUserContext(r); err != nil {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), err.Error())
		return
	}

	// 检查系统是否启用2FA
	if !s.Is2FAGlobalEnabled(ctx) {
		r.Middleware.Next()
		return
	}

	// 获取当前用户ID
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		response.JsonExit(r, gcode.CodeNotAuthorized.Code(), "用户未登录")
		return
	}

	// 检查用户是否启用了2FA
	isEnabled, err := service.Admin2fa().IsEnabled(ctx, uint64(adminId))
	if err != nil {
		g.Log().Errorf(ctx, "检查用户2FA状态失败: %v", err)
		response.JsonExit(r, gcode.CodeInternalError.Code(), "系统错误")
		return
	}

	// 如果用户未启用2FA，检查是否强制启用
	if !isEnabled {
		if s.Is2FAEnforced(ctx) {
			// 强制启用2FA，返回需要设置2FA的错误
			response.JsonExit(r, gcode.CodeSecurityReason.Code(), "系统要求启用2FA，请先设置2FA")
			return
		}
		// 非强制模式，允许通过
		r.Middleware.Next()
		return
	}

	// 检查会话中是否已经通过2FA验证
	if s.Is2FAVerified(r) {
		r.Middleware.Next()
		return
	}

	// 需要2FA验证
	response.JsonExit(r, gcode.CodeSecurityReason.Code(), "需要2FA验证")
}

// IsExcept2FA 是否是不需要验证2FA的路由地址
func (s *sMiddleware) IsExcept2FA(ctx context.Context, path string) bool {
	// 2FA相关的API路由不需要验证2FA
	except2FAList := []string{
		"/admin2fa/setup",
		"/admin2fa/verify-setup",
		"/admin2fa/verify",
		"/admin2fa/disable",
		"/admin2fa/status",
		"/admin2fa/backup-codes",
		"/admin2fa/regenerate-backup-codes",
		"/admin2fa/use-backup-code",
		"/admin2fa/log/list",
		"/site/accountLogin",
		"/site/mobileLogin",
		"/site/logout",
	}

	for _, exceptPath := range except2FAList {
		if gstr.Contains(path, exceptPath) {
			return true
		}
	}

	return false
}

// Is2FAGlobalEnabled 检查系统是否全局启用2FA
func (s *sMiddleware) Is2FAGlobalEnabled(ctx context.Context) bool {
	// 从系统配置中获取基础配置
	config, err := service.SysConfig().GetBasic(ctx)
	if err != nil || config == nil {
		return false
	}

	// 检查2FA全局开关
	return config.TwoFASwitch == 1
}

// Is2FAEnforced 检查系统是否强制启用2FA
func (s *sMiddleware) Is2FAEnforced(ctx context.Context) bool {
	// 从系统配置中获取基础配置
	config, err := service.SysConfig().GetBasic(ctx)
	if err != nil || config == nil {
		return false
	}

	// 检查2FA强制开关
	return config.TwoFAForce == 1
}

// Is2FAVerified 检查会话中是否已经通过2FA验证
func (s *sMiddleware) Is2FAVerified(r *ghttp.Request) bool {
	// 从会话中获取2FA验证状态
	session := r.Session
	if session == nil {
		return false
	}

	// 检查会话中的2FA验证标志
	verified := session.MustGet("2fa_verified").Bool()
	adminId := session.MustGet("admin_id").Int64()
	currentAdminId := contexts.GetUserId(r.Context())

	// 确保会话中的用户ID与当前用户ID一致
	return verified && adminId == currentAdminId
}

// Set2FAVerified 设置会话中的2FA验证状态
func (s *sMiddleware) Set2FAVerified(r *ghttp.Request, adminId int64) error {
	session := r.Session
	if session == nil {
		return gerror.New("会话不存在")
	}

	// 设置2FA验证标志
	if err := session.Set("2fa_verified", true); err != nil {
		return err
	}

	// 设置用户ID
	if err := session.Set("admin_id", adminId); err != nil {
		return err
	}

	return nil
}

// Clear2FAVerified 清除会话中的2FA验证状态
func (s *sMiddleware) Clear2FAVerified(r *ghttp.Request) error {
	session := r.Session
	if session == nil {
		return nil
	}

	// 清除2FA验证标志
	if err := session.Remove("2fa_verified"); err != nil {
		return err
	}

	// 清除用户ID
	if err := session.Remove("admin_id"); err != nil {
		return err
	}

	return nil
}