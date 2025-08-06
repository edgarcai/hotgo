// Package middleware
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package middleware

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"go.opentelemetry.io/otel/attribute"
	"hotgo/internal/consts"
	"hotgo/internal/library/addons"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/response"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/service"
	"hotgo/utility/simple"
	"hotgo/utility/validate"
	"net/http"
	"strings"
)

type sMiddleware struct {
	LoginUrl         string // 登录路由地址
	DemoWhiteList    g.Map  // 演示模式放行的路由白名单
	NotRecordRequest g.Map  // 不记录请求数据的路由（当前请求数据过大时会影响响应效率，可以将路径放到该选项中改善）
}

func init() {
	service.RegisterMiddleware(NewMiddleware())
}

func NewMiddleware() *sMiddleware {
	return &sMiddleware{
		LoginUrl: "/common",
		DemoWhiteList: g.Map{
			"/admin/site/accountLogin": struct{}{}, // 账号登录
			"/admin/site/mobileLogin":  struct{}{}, // 手机号登录
			"/admin/genCodes/preview":  struct{}{}, // 预览代码
		},
		NotRecordRequest: g.Map{
			"/admin/upload/file":       struct{}{}, // 上传文件
			"/admin/upload/uploadPart": struct{}{}, // 上传分片
		},
	}
}

// Ctx 初始化请求上下文
func (s *sMiddleware) Ctx(r *ghttp.Request) {
	if g.Cfg().MustGet(r.Context(), "jaeger.switch").Bool() {
		ctx, span := gtrace.NewSpan(r.Context(), "middleware.ctx")
		span.SetAttributes(attribute.KeyValue{
			Key:   "traceID",
			Value: attribute.StringValue(gctx.CtxId(ctx)),
		})
		span.End()
		r.SetCtx(ctx)
	}

	data := make(g.Map)
	if _, ok := s.NotRecordRequest[r.URL.Path]; ok {
		data["request.body"] = gjson.New(nil)
	} else {
		data["request.body"] = gjson.New(r.GetBodyString())
	}

	contexts.Init(r, &model.Context{
		Data:   data,
		Module: getModule(r.URL.Path),
	})

	if len(r.Cookie.GetSessionId()) == 0 {
		r.Cookie.SetSessionId(gctx.CtxId(r.Context()))
	}

	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

func getModule(path string) (module string) {
	slice := strings.Split(path, "/")
	if len(slice) < 2 {
		module = consts.AppDefault
		return
	}

	if slice[1] == "" {
		module = consts.AppDefault
		return
	}
	return slice[1]
}

// CORS allows Cross-origin resource sharing.
func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

// DemoLimit 演示系统操作限制
func (s *sMiddleware) DemoLimit(r *ghttp.Request) {
	if !simple.IsDemo(r.Context()) {
		r.Middleware.Next()
		return
	}

	if r.Method == http.MethodPost {
		if _, ok := s.DemoWhiteList[r.URL.Path]; ok {
			r.Middleware.Next()
			return
		}
		response.JsonExit(r, gcode.CodeNotSupported.Code(), "演示系统禁止操作！")
		return
	}

	r.Middleware.Next()
}

// Addon 插件中间件
func (s *sMiddleware) Addon(r *ghttp.Request) {
	var ctx = r.Context()

	if contexts.Get(ctx).Module == "" {
		g.Log().Warning(ctx, "application module is not initialized.")
		return
	}

	// 替换掉应用模块前缀
	path := gstr.Replace(r.URL.Path, "/"+contexts.Get(ctx).Module+"/", "", 1)
	ss := gstr.Explode("/", path)
	if len(ss) == 0 {
		g.Log().Warning(ctx, "addon was not recognized.")
		return
	}

	module := addons.GetModule(ss[0])
	if module == nil {
		g.Log().Warningf(ctx, "addon module = nil, name:%v", ss[0])
		return
	}

	sk := module.GetSkeleton()
	if sk == nil {
		g.Log().Warningf(ctx, "addon skeleton = nil, name:%v", ss[0])
		return
	}

	contexts.SetAddonName(ctx, sk.Name)
	r.Middleware.Next()
}

// DeliverUserContext 将用户信息传递到上下文中
func (s *sMiddleware) DeliverUserContext(r *ghttp.Request) (err error) {
	user, err := token.ParseLoginUser(r)
	if err != nil {
		return
	}

	switch user.App {
	case consts.AppAdmin:
		if err = service.AdminSite().BindUserContext(r.Context(), user); err != nil {
			return
		}
	default:
		contexts.SetUser(r.Context(), user)
	}
	return
}

// IsExceptAuth 是否是不需要验证权限的路由地址
func (s *sMiddleware) IsExceptAuth(ctx context.Context, appName, path string) bool {
	pathList := g.Cfg().MustGet(ctx, fmt.Sprintf("router.%v.exceptAuth", appName)).Strings()

	for i := 0; i < len(pathList); i++ {
		if validate.InSliceExistStr(pathList[i], path) {
			return true
		}
	}
	return false
}

// IsExceptLogin 是否是不需要登录的路由地址
func (s *sMiddleware) IsExceptLogin(ctx context.Context, appName, path string) bool {
	pathList := g.Cfg().MustGet(ctx, fmt.Sprintf("router.%v.exceptLogin", appName)).Strings()

	for i := 0; i < len(pathList); i++ {
		if validate.InSliceExistStr(pathList[i], path) {
			return true
		}
	}
	return false
}

// IsExcept2FA 是否是不需要2FA验证的路由地址
func (s *sMiddleware) IsExcept2FA(ctx context.Context, path string) bool {
	// 2FA相关路由不需要验证2FA状态
	except2FAPaths := []string{
		"/auth/verify-login-2fa", // 2FA验证接口
		"/auth/login",           // 登录接口
		"/auth/logout",          // 登出接口
		"/auth/2fa-status",      // 2FA状态查询接口
		"/auth/enable-2fa",      // 启用2FA接口
		"/auth/verify-2fa-setup", // 验证2FA设置接口
		"/auth/disable-2fa",     // 禁用2FA接口
		"/auth/regenerate-backup-codes", // 重新生成备用码接口
		"/two-factor",           // 2FA管理相关接口
	}

	for _, exceptPath := range except2FAPaths {
		if validate.InSliceExistStr(exceptPath, path) {
			return true
		}
	}
	return false
}

// Check2FAVerification 检查2FA验证状态
func (s *sMiddleware) Check2FAVerification(r *ghttp.Request) error {
	ctx := r.Context()
	user := contexts.GetUser(ctx)
	if user == nil {
		return gerror.New("用户信息获取失败")
	}

	// 检查用户是否启用了2FA
	twoFactorStatus, err := service.AdminTwoFactor().GetStatus(ctx, &adminin.TwoFactorGetStatusInp{
		UserId: user.Id,
	})
	if err != nil {
		return err
	}

	// 如果用户启用了2FA但未完成验证，拒绝访问
	if twoFactorStatus.IsEnabled && !user.TwoFactorVerified {
		return gerror.New("需要完成双因子认证才能访问此资源")
	}

	return nil
}
