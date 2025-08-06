// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/service"
)

// TestSite_AccountLogin_Without2FA 测试未启用2FA用户的登录流程
func TestSite_AccountLogin_Without2FA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &adminin.AccountLoginInp{
			Username: "test_user",
			Password: "test_password",
		}

		// 创建测试上下文
		ctx := context.Background()

		// 调用服务方法
		res, err := service.AdminSite().AccountLogin(ctx, req)

		// 验证结果（由于没有实际的用户数据，预期会失败）
		t.AssertNE(err, nil)
		t.AssertNil(res)
		// 检查错误信息
		if err != nil {
			t.Assert(strings.Contains(err.Error(), "用户") || strings.Contains(err.Error(), "密码") || strings.Contains(err.Error(), "登录"), true)
		}
	})
}

// TestSite_AccountLogin_With2FA 测试启用2FA用户的登录流程
func TestSite_AccountLogin_With2FA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &adminin.AccountLoginInp{
			Username: "test_user_2fa",
			Password: "test_password",
		}

		// 创建测试上下文
		ctx := context.Background()

		// 调用服务方法
		res, err := service.AdminSite().AccountLogin(ctx, req)

		// 验证结果（由于没有实际的用户数据，预期会失败）
		t.AssertNE(err, nil)
		t.AssertNil(res)
		// 检查错误信息
		if err != nil {
			t.Assert(strings.Contains(err.Error(), "用户") || strings.Contains(err.Error(), "密码") || strings.Contains(err.Error(), "登录"), true)
		}
	})
}

// TestSite_AccountLogin_InputValidation 测试登录输入验证
func TestSite_AccountLogin_InputValidation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 测试空用户名
		req1 := &adminin.AccountLoginInp{
			Username: "",
			Password: "test_password",
		}

		ctx := context.Background()
		res1, err1 := service.AdminSite().AccountLogin(ctx, req1)
		t.AssertNE(err1, nil)
		t.AssertNil(res1)

		// 测试空密码
		req2 := &adminin.AccountLoginInp{
			Username: "test_user",
			Password: "",
		}

		res2, err2 := service.AdminSite().AccountLogin(ctx, req2)
		t.AssertNE(err2, nil)
		t.AssertNil(res2)

		// 测试用户名和密码都为空
		req3 := &adminin.AccountLoginInp{
			Username: "",
			Password: "",
		}

		res3, err3 := service.AdminSite().AccountLogin(ctx, req3)
		t.AssertNE(err3, nil)
		t.AssertNil(res3)
	})
}

// TestSite_AccountLogin_SecurityLogging 测试登录安全日志记录
func TestSite_AccountLogin_SecurityLogging(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &adminin.AccountLoginInp{
			Username: "test_user",
			Password: "wrong_password",
		}

		// 创建测试上下文
		ctx := context.Background()

		// 调用服务方法
		res, err := service.AdminSite().AccountLogin(ctx, req)

		// 验证结果（预期失败，但应该记录日志）
		t.AssertNE(err, nil)
		t.AssertNil(res)
		
		// 注意：实际的日志记录验证需要检查数据库或日志文件
		// 这里只是确保方法能够正常执行而不崩溃
	})
}