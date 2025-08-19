// Package middleware AniX认证中间件
package middleware

import (
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
)

// AnixAuth AniX用户认证中间件
func (s *sMiddleware) AnixAuth(r *ghttp.Request) {
	// 从Header获取Token
	token := r.Header.Get("Authorization")
	if token == "" {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "未提供访问令牌",
			"data":    nil,
		})
		return
	}

	// 移除Bearer前缀
	if gstr.HasPrefix(token, "Bearer ") {
		token = gstr.SubStr(token, 7)
	}

	// 验证Token
	userId, err := service.AnixUser().VerifyToken(r.Context(), token)
	if err != nil {
		g.Log().Warning(r.Context(), "Token验证失败:", err)
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "访问令牌无效",
			"data":    nil,
		})
		return
	}

	// 获取用户信息
	user, err := service.AnixUser().GetProfile(r.Context(), userId)
	if err != nil {
		g.Log().Warning(r.Context(), "获取用户信息失败:", err)
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "用户不存在或已被禁用",
			"data":    nil,
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		r.Response.WriteJsonExit(g.Map{
			"code":    403,
			"message": "用户已被禁用",
			"data":    nil,
		})
		return
	}

	// 将用户ID存储到上下文中
	r.SetCtxVar("user_id", userId)
	r.SetCtxVar("user_info", user)

	// 继续处理请求
	r.Middleware.Next()
}

// AnixOptionalAuth AniX可选认证中间件（不强制要求登录）
func (s *sMiddleware) AnixOptionalAuth(r *ghttp.Request) {
	// 从Header获取Token
	token := r.Header.Get("Authorization")
	if token == "" {
		// 没有Token，继续处理请求
		r.Middleware.Next()
		return
	}

	// 移除Bearer前缀
	if gstr.HasPrefix(token, "Bearer ") {
		token = gstr.SubStr(token, 7)
	}

	// 验证Token
	userId, err := service.AnixUser().VerifyToken(r.Context(), token)
	if err != nil {
		// Token无效，但不阻止请求
		g.Log().Debug(r.Context(), "可选认证Token验证失败:", err)
		r.Middleware.Next()
		return
	}

	// 获取用户信息
	user, err := service.AnixUser().GetProfile(r.Context(), userId)
	if err != nil {
		// 用户不存在，但不阻止请求
		g.Log().Debug(r.Context(), "可选认证获取用户信息失败:", err)
		r.Middleware.Next()
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		// 用户被禁用，但不阻止请求
		g.Log().Debug(r.Context(), "可选认证用户已被禁用")
		r.Middleware.Next()
		return
	}

	// 将用户ID存储到上下文中
	r.SetCtxVar("user_id", userId)
	r.SetCtxVar("user_info", user)

	// 继续处理请求
	r.Middleware.Next()
}