// Package router
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package router

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	"hotgo/internal/controller/anix/chapter"
	"hotgo/internal/controller/anix/comic"
	"hotgo/internal/controller/anix/history"
	"hotgo/internal/controller/anix/image"
	"hotgo/internal/controller/anix/user"
	"hotgo/internal/service"
	"hotgo/utility/simple"
)

// Anix AniX漫画应用路由
func Anix(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group(simple.RouterPrefix(ctx, "anix"), func(group *ghttp.RouterGroup) {
		// 不需要认证的路由
		group.Bind(
			comic.NewV1(),   // 漫画相关（列表、详情、搜索等）
			chapter.NewV1(), // 章节相关（详情、页面等）
			image.NewV1(),   // 图片服务（处理、缓存、优化等）
		)

		// 用户注册和登录路由（不需要认证）
		userController := user.NewV1()
		group.POST("/register", userController.Register)     // 用户注册
		group.POST("/login", userController.Login)           // 用户登录

		// 需要认证的用户路由
		group.Middleware(service.Middleware().AnixAuth)
		group.GET("/profile", userController.Profile)           // 获取用户信息
		group.PUT("/profile", userController.UpdateProfile)     // 更新用户信息
		group.POST("/change-password", userController.ChangePassword) // 修改密码
		group.POST("/logout", userController.Logout)           // 用户登出

		// 阅读历史相关路由（需要认证）
		group.Bind(
			history.NewV1(),         // 阅读历史、书签、统计等
		)
	})
}