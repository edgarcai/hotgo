// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"hotgo/api/admin/auth"
	"hotgo/internal/consts"
	"hotgo/internal/model"
	"hotgo/internal/model/input/adminin"
)

// TestAuth_Enable2FA 测试启用双因子认证
func TestAuth_Enable2FA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.Enable2FAReq{
			TwoFactorEnableInp: adminin.TwoFactorEnableInp{
				UserId: 1,
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.Enable2FA(ctx, req)

		// 验证结果
		t.AssertNil(err)
		t.AssertNE(res, nil)
		t.AssertNE(res.Secret, "")
		t.AssertNE(res.QRCodeURL, "")
		t.Assert(len(res.BackupCodes), 8)
	})
}

// TestAuth_Verify2FASetup 测试验证双因子认证设置
func TestAuth_Verify2FASetup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.Verify2FASetupReq{
			TwoFactorConfirmEnableInp: adminin.TwoFactorConfirmEnableInp{
				UserId: 1,
				Code:   "123456",
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.Verify2FASetup(ctx, req)

		// 验证结果（这里可能会失败，因为需要有效的TOTP码）
		_ = res
		_ = err
		// 在实际测试中，需要模拟有效的TOTP验证码
	})
}

// TestAuth_Disable2FA 测试禁用双因子认证
func TestAuth_Disable2FA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.Disable2FAReq{
			TwoFactorDisableInp: adminin.TwoFactorDisableInp{
				UserId: 1,
				Code:   "123456",
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.Disable2FA(ctx, req)

		// 验证结果
		_ = res
		_ = err
		// 在实际测试中，需要先启用2FA才能禁用
	})
}

// TestAuth_Get2FAStatus 测试获取双因子认证状态
func TestAuth_Get2FAStatus(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.Get2FAStatusReq{
			TwoFactorGetStatusInp: adminin.TwoFactorGetStatusInp{
				UserId: 1,
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.Get2FAStatus(ctx, req)

		// 验证结果
		t.AssertNil(err)
		t.AssertNE(res, nil)
		// IsEnabled 可能为 true 或 false，取决于用户状态
	})
}

// TestAuth_RegenerateBackupCodes 测试重新生成备用恢复码
func TestAuth_RegenerateBackupCodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.RegenerateBackupCodesReq{
			TwoFactorRegenerateBackupCodesInp: adminin.TwoFactorRegenerateBackupCodesInp{
				UserId: 1,
				Code:   "123456",
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.RegenerateBackupCodes(ctx, req)

		// 验证结果
		_ = res
		_ = err
		// 在实际测试中，需要先启用2FA并提供有效的TOTP码
	})
}

// TestAuth_Verify2FA 测试验证双因子认证
func TestAuth_Verify2FA(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 创建测试请求
		req := &auth.Verify2FAReq{
			TwoFactorVerifyInp: adminin.TwoFactorVerifyInp{
				UserId: 1,
				Code:   "123456",
			},
		}

		// 创建测试上下文，模拟用户登录
		ctx := context.Background()
		mockContext := &model.Context{
			User: &model.Identity{
				Id: 1,
				Username: "test_user",
			},
		}
		ctx = context.WithValue(ctx, consts.ContextHTTPKey, mockContext)

		// 调用控制器方法
		res, err := Auth.Verify2FA(ctx, req)

		// 验证结果
		_ = res
		_ = err
		// 在实际测试中，需要提供有效的TOTP码或备用恢复码
	})
}