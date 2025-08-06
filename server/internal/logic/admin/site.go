// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/token"
	"hotgo/internal/model"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/simple"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

type sAdminSite struct{}

func NewAdminSite() *sAdminSite {
	return &sAdminSite{}
}

func init() {
	service.RegisterAdminSite(NewAdminSite())
}

// Register 账号注册
func (s *sAdminSite) Register(ctx context.Context, in *adminin.RegisterInp) (err error) {
	config, err := service.SysConfig().GetLogin(ctx)
	if err != nil {
		return
	}

	if config.ForceInvite == 1 && in.InviteCode == "" {
		err = gerror.New("请填写邀请码")
		return
	}

	var data adminin.MemberAddInp

	// 默认上级
	data.Pid = 1

	// 存在邀请人
	if in.InviteCode != "" {
		pmb, err := service.AdminMember().GetIdByCode(ctx, &adminin.GetIdByCodeInp{Code: in.InviteCode})
		if err != nil {
			return err
		}

		if pmb == nil {
			err = gerror.New("邀请人信息不存在")
			return err
		}

		data.Pid = pmb.Id
	}

	if config.RegisterSwitch != 1 {
		err = gerror.New("管理员未开放注册")
		return
	}

	if config.RoleId < 1 {
		err = gerror.New("管理员未配置默认角色")
		return
	}

	if config.DeptId < 1 {
		err = gerror.New("管理员未配置默认部门")
		return
	}

	// 验证唯一性
	err = service.AdminMember().VerifyUnique(ctx, &adminin.VerifyUniqueInp{
		Where: g.Map{
			dao.AdminMember.Columns().Username: in.Username,
			dao.AdminMember.Columns().Mobile:   in.Mobile,
		},
	})
	if err != nil {
		return
	}

	// 验证短信验证码
	err = service.SysSmsLog().VerifyCode(ctx, &sysin.VerifyCodeInp{
		Event:  consts.SmsTemplateRegister,
		Mobile: in.Mobile,
		Code:   in.Code,
	})
	if err != nil {
		return
	}

	data.MemberEditInp = &adminin.MemberEditInp{
		Id:       0,
		RoleId:   config.RoleId,
		PostIds:  config.PostIds,
		DeptId:   config.DeptId,
		Username: in.Username,
		Password: in.Password,
		RealName: "",
		Avatar:   config.Avatar,
		Sex:      3, // 保密
		Mobile:   in.Mobile,
		Status:   consts.StatusEnabled,
	}
	data.Salt = grand.S(6)
	data.InviteCode = grand.S(12)
	data.PasswordHash = gmd5.MustEncryptString(data.Password + data.Salt)
	data.Level, data.Tree, err = service.AdminMember().GenTree(ctx, data.Pid)
	if err != nil {
		return
	}

	// 提交注册信息
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		id, err := dao.AdminMember.Ctx(ctx).Data(data).OmitEmptyData().InsertAndGetId()
		if err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
			return
		}

		// 更新岗位
		if err = service.AdminMemberPost().UpdatePostIds(ctx, id, config.PostIds); err != nil {
			err = gerror.Wrap(err, consts.ErrorORM)
		}
		return
	})
}

// AccountLogin 账号登录
func (s *sAdminSite) AccountLogin(ctx context.Context, in *adminin.AccountLoginInp) (res *adminin.LoginModel, err error) {
	var twoFactorEnabled bool
	var loginStep string = "password_verification"
	
	defer func() {
		// 记录登录日志，包含2FA相关信息
		logInp := &sysin.LoginLogPushInp{
			Response:          res,
			Err:              err,
			TwoFactorEnabled:  twoFactorEnabled,
			TwoFactorMethod:   "", // 在密码验证阶段还未使用2FA方法
			TwoFactorSuccess:  false, // 密码验证阶段2FA还未成功
			LoginStep:        loginStep,
		}
		service.SysLoginLog().Push(ctx, logInp)
	}()

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where("username", in.Username).Scan(&mb); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if mb == nil {
		err = gerror.New("用户名或密码错误")
		return
	}

	res = new(adminin.LoginModel)
	res.Id = mb.Id
	res.Username = mb.Username
	if mb.Salt == "" {
		err = gerror.New("用户信息错误")
		return
	}

	if err = simple.CheckPassword(in.Password, mb.Salt, mb.PasswordHash); err != nil {
		return
	}

	if mb.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用")
		return
	}

	// 检查用户是否启用了2FA
	twoFactorStatus, err := service.AdminTwoFactor().GetStatus(ctx, &adminin.TwoFactorGetStatusInp{
		UserId: mb.Id,
	})
	if err != nil {
		return
	}

	twoFactorEnabled = twoFactorStatus.IsEnabled

	// 如果启用了2FA，返回需要2FA验证的响应
	if twoFactorStatus.IsEnabled {
		// 生成临时token用于2FA验证
		var tempToken string
		tempToken, err = s.generateTempToken(ctx, mb.Id)
		if err != nil {
			return
		}
		
		// 返回需要2FA验证的响应
		res.RequiresTwoFactor = true
		res.TempToken = tempToken
		loginStep = "awaiting_2fa"
		return
	}

	// 未启用2FA，直接完成登录
	res, err = s.handleLoginWithTwoFactor(ctx, mb, false)
	if err == nil {
		loginStep = "login_complete"
	}
	return
}

// MobileLogin 手机号登录
func (s *sAdminSite) MobileLogin(ctx context.Context, in *adminin.MobileLoginInp) (res *adminin.LoginModel, err error) {
	defer func() {
		service.SysLoginLog().Push(ctx, &sysin.LoginLogPushInp{Response: res, Err: err})
	}()

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where("mobile ", in.Mobile).Scan(&mb); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if mb == nil {
		err = gerror.New("账号不存在")
		return
	}

	res = new(adminin.LoginModel)
	res.Id = mb.Id
	res.Username = mb.Username

	err = service.SysSmsLog().VerifyCode(ctx, &sysin.VerifyCodeInp{
		Event:  consts.SmsTemplateLogin,
		Mobile: in.Mobile,
		Code:   in.Code,
	})

	if err != nil {
		return
	}

	if mb.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用")
		return
	}

	res, err = s.handleLogin(ctx, mb)
	return
}

// handleLogin .
// generateTempToken 生成临时token用于2FA验证
func (s *sAdminSite) generateTempToken(ctx context.Context, userId int64) (string, error) {
	// 生成临时token，有效期5分钟
	identity := &model.Identity{
		Id:       userId,
		RoleKey:  "temp_2fa", // 临时角色，仅用于2FA验证
		DeptId:   0,
		Username: "",
		RealName: "",
		Avatar:   "",
		Email:    "",
		Mobile:   "",
		App:      consts.AppAdmin,
		LoginAt:  gtime.Now(),
	}

	// 生成临时token，有效期5分钟
	tempToken, _, err := token.Login(ctx, identity)
	if err != nil {
		return "", err
	}

	return tempToken, nil
}

// HandleLogin 公共登录处理方法
func (s *sAdminSite) HandleLogin(ctx context.Context, member *entity.AdminMember) (res *adminin.LoginModel, err error) {
	return s.handleLogin(ctx, member)
}

// handleLogin 处理登录逻辑，根据调用场景设置2FA验证状态
func (s *sAdminSite) handleLogin(ctx context.Context, member *entity.AdminMember) (res *adminin.LoginModel, err error) {
	return s.handleLoginWithTwoFactor(ctx, member, true)
}

// handleLoginWithTwoFactor 处理登录逻辑，可指定2FA验证状态
func (s *sAdminSite) handleLoginWithTwoFactor(ctx context.Context, member *entity.AdminMember, twoFactorVerified bool) (res *adminin.LoginModel, err error) {
	role, dept, err := s.getLoginRoleAndDept(ctx, member.RoleId, member.DeptId)
	if err != nil {
		return nil, err
	}

	user := &model.Identity{
		Id:                member.Id,
		Pid:               member.Pid,
		DeptId:            dept.Id,
		DeptType:          dept.Type,
		RoleId:            role.Id,
		RoleKey:           role.Key,
		Username:          member.Username,
		RealName:          member.RealName,
		Avatar:            member.Avatar,
		Email:             member.Email,
		Mobile:            member.Mobile,
		App:               consts.AppAdmin,
		LoginAt:           gtime.Now(),
		TwoFactorVerified: twoFactorVerified,
	}

	lt, expires, err := token.Login(ctx, user)
	if err != nil {
		return nil, err
	}

	res = &adminin.LoginModel{
		Username: user.Username,
		Id:       user.Id,
		Token:    lt,
		Expires:  expires,
	}
	return
}

// getLoginRoleAndDept 获取登录的角色和部门信息
func (s *sAdminSite) getLoginRoleAndDept(ctx context.Context, roleId, deptId int64) (role *entity.AdminRole, dept *entity.AdminDept, err error) {
	if err = dao.AdminRole.Ctx(ctx).Fields(dao.AdminRole.Columns().Id, dao.AdminRole.Columns().Key, dao.AdminRole.Columns().Status).WherePri(roleId).Scan(&role); err != nil {
		err = gerror.Wrap(err, consts.ErrorORM)
		return
	}

	if role == nil {
		err = gerror.New("角色不存在或已被删除")
		return
	}

	if role.Status != consts.StatusEnabled {
		err = gerror.New("角色已被禁用，如有疑问请联系管理员")
		return
	}

	if err = dao.AdminDept.Ctx(ctx).Fields(dao.AdminDept.Columns().Id, dao.AdminDept.Columns().Type, dao.AdminDept.Columns().Status).WherePri(deptId).Scan(&dept); err != nil {
		err = gerror.Wrap(err, "获取部门信息失败，请稍后重试！")
		return
	}

	if dept == nil {
		err = gerror.New("部门不存在或已被删除")
		return
	}

	if dept.Status != consts.StatusEnabled {
		err = gerror.New("部门已被禁用，如有疑问请联系管理员")
		return
	}
	return
}

// BindUserContext 绑定用户上下文
func (s *sAdminSite) BindUserContext(ctx context.Context, claims *model.Identity) (err error) {
	//// 如果不想每次访问都重新加载用户信息，可以放开注释。但在本次登录未失效前，用户信息不会刷新
	// contexts.SetUser(ctx, claims)
	// return

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(claims.Id).Scan(&mb); err != nil {
		err = gerror.Wrap(err, "获取用户信息失败，请稍后重试！")
		return
	}

	if mb == nil {
		err = gerror.Wrap(err, "账号不存在或已被删除！")
		return
	}

	if mb.Status != consts.StatusEnabled {
		err = gerror.New("账号已被禁用，如有疑问请联系管理员")
		return
	}

	role, dept, err := s.getLoginRoleAndDept(ctx, mb.RoleId, mb.DeptId)
	if err != nil {
		return err
	}

	user := &model.Identity{
		Id:                mb.Id,
		Pid:               mb.Pid,
		DeptId:            dept.Id,
		DeptType:          dept.Type,
		RoleId:            mb.RoleId,
		RoleKey:           role.Key,
		Username:          mb.Username,
		RealName:          mb.RealName,
		Avatar:            mb.Avatar,
		Email:             mb.Email,
		Mobile:            mb.Mobile,
		App:               claims.App,
		LoginAt:           claims.LoginAt,
		TwoFactorVerified: claims.TwoFactorVerified, // 保持原有的2FA验证状态
	}

	contexts.SetUser(ctx, user)
	return
}
