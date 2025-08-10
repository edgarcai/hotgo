package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/library/totp"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/service"
)

type sAdmin2fa struct{}

func NewAdmin2fa() *sAdmin2fa {
	return &sAdmin2fa{}
}

func init() {
	service.RegisterAdmin2fa(NewAdmin2fa())
}

// Setup 设置2FA
func (s *sAdmin2fa) Setup(ctx context.Context, in *adminin.Admin2faSetupInp) (res *adminin.Admin2faSetupModel, err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 检查是否已经设置过2FA
	existingConfig, err := dao.Admin2faConfig.Ctx(ctx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).One()
	if err != nil {
		return nil, err
	}
	if !existingConfig.IsEmpty() {
		return nil, gerror.New("2FA已经设置，请先禁用后重新设置")
	}

	// 生成TOTP密钥
	config := totp.DefaultConfig()
	secretKey, err := totp.GenerateSecret(config)
	if err != nil {
		return nil, gerror.Wrap(err, "生成TOTP密钥失败")
	}

	// 获取用户信息用于生成二维码
	adminInfo, err := dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Id, adminId).One()
	if err != nil {
		return nil, err
	}
	if adminInfo.IsEmpty() {
		return nil, gerror.New("用户不存在")
	}

	// 生成二维码URL
	issuer := "HotGo"
	accountName := gconv.String(adminInfo["username"])
	qrCodeURL := totp.GenerateQRCodeURL(secretKey, issuer, accountName)

	// 生成备用恢复码
	backupCodes, err := s.generateBackupCodes(8)
	if err != nil {
		return nil, gerror.Wrap(err, "生成备用恢复码失败")
	}

	// 开启事务
	err = dao.Admin2faConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 保存2FA配置（未验证状态）
		_, err = dao.Admin2faConfig.Ctx(ctx).TX(tx).Insert(&entity.Admin2faConfig{
			AdminId:    adminId,
			SecretKey:  secretKey,
			IsEnabled:  consts.StatusDisable, // 初始状态为禁用
			IsVerified: consts.StatusDisable, // 初始状态为未验证
			CreatedAt:  gtime.Now(),
			UpdatedAt:  gtime.Now(),
		})
		if err != nil {
			return err
		}

		// 保存备用恢复码
		for _, code := range backupCodes {
			_, err = dao.Admin2faBackupCode.Ctx(ctx).TX(tx).Insert(&entity.Admin2faBackupCode{
				AdminId:   adminId,
				Code:      code,
				IsUsed:    consts.StatusDisable,
				CreatedAt: gtime.Now(),
			})
			if err != nil {
				return err
			}
		}

		// 记录操作日志
		return s.logAction(ctx, tx, adminId, "setup", "success", "设置2FA")
	})

	if err != nil {
		return nil, err
	}

	return &adminin.Admin2faSetupModel{
		SecretKey:   secretKey,
		QRCodeURL:   qrCodeURL,
		BackupCodes: backupCodes,
	}, nil
}

// VerifySetup 验证2FA设置
func (s *sAdmin2fa) VerifySetup(ctx context.Context, in *adminin.Admin2faVerifySetupInp) (err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return gerror.New("用户未登录")
	}

	// 获取2FA配置
	config, err := dao.Admin2faConfig.Ctx(ctx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).One()
	if err != nil {
		return err
	}
	if config.IsEmpty() {
		return gerror.New("请先设置2FA")
	}

	// 验证验证码
	secretKey := gconv.String(config["secret_key"])
	totpConfig := totp.DefaultConfig()
	if !totp.ValidateCode(secretKey, in.Code, totpConfig) {
		// 记录失败日志
		s.logAction(ctx, nil, adminId, "verify_setup", "failed", "验证码错误")
		return gerror.New("验证码错误")
	}

	// 启用2FA
	_, err = dao.Admin2faConfig.Ctx(ctx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).Update(g.Map{
		dao.Admin2faConfig.Columns().IsEnabled:  consts.StatusEnabled,
		dao.Admin2faConfig.Columns().IsVerified: consts.StatusEnabled,
		dao.Admin2faConfig.Columns().UpdatedAt:  gtime.Now(),
	})
	if err != nil {
		return err
	}

	// 记录成功日志
	return s.logAction(ctx, nil, adminId, "verify_setup", "success", "验证设置成功")
}

// Verify 验证2FA
func (s *sAdmin2fa) Verify(ctx context.Context, in *adminin.Admin2faVerifyInp) (err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return gerror.New("用户未登录")
	}

	// 获取2FA配置
	config, err := dao.Admin2faConfig.Ctx(ctx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).One()
	if err != nil {
		return err
	}
	if config.IsEmpty() {
		return gerror.New("用户未设置2FA")
	}

	isEnabled := gconv.Int(config["is_enabled"])
	isVerified := gconv.Int(config["is_verified"])
	if isEnabled != consts.StatusEnabled || isVerified != consts.StatusEnabled {
		return gerror.New("2FA未启用或未验证")
	}

	// 验证验证码
	secretKey := gconv.String(config["secret_key"])
	totpConfig := totp.DefaultConfig()
	if !totp.ValidateCode(secretKey, in.Code, totpConfig) {
		// 记录失败日志
		s.logAction(ctx, nil, adminId, "verify", "failed", "验证码错误")
		return gerror.New("验证码错误")
	}

	// 记录成功日志
	return s.logAction(ctx, nil, adminId, "verify", "success", "验证成功")
}

// Disable 禁用2FA
func (s *sAdmin2fa) Disable(ctx context.Context, in *adminin.Admin2faDisableInp) (err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return gerror.New("用户未登录")
	}

	// 验证当前验证码
	verifyInp := &adminin.Admin2faVerifyInp{
		Code: in.Code,
	}
	if err = s.Verify(ctx, verifyInp); err != nil {
		return err
	}

	// 开启事务删除所有2FA相关数据
	err = dao.Admin2faConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除2FA配置
		_, err = dao.Admin2faConfig.Ctx(ctx).TX(tx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).Delete()
		if err != nil {
			return err
		}

		// 删除备用恢复码
		_, err = dao.Admin2faBackupCode.Ctx(ctx).TX(tx).Where(dao.Admin2faBackupCode.Columns().AdminId, adminId).Delete()
		if err != nil {
			return err
		}

		// 记录操作日志
		return s.logAction(ctx, tx, adminId, "disable", "success", "禁用2FA")
	})

	return err
}

// Status 获取2FA状态
func (s *sAdmin2fa) Status(ctx context.Context, in *adminin.Admin2faStatusInp) (res *adminin.Admin2faStatusModel, err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 获取2FA配置
	config, err := dao.Admin2faConfig.Ctx(ctx).Where(dao.Admin2faConfig.Columns().AdminId, adminId).One()
	if err != nil {
		return nil, err
	}

	res = &adminin.Admin2faStatusModel{
		IsEnabled:      false,
		IsVerified:     false,
		HasBackupCodes: false,
	}

	if !config.IsEmpty() {
		isEnabled := gconv.Int(config["is_enabled"])
		isVerified := gconv.Int(config["is_verified"])
		res.IsEnabled = isEnabled == consts.StatusEnabled
		res.IsVerified = isVerified == consts.StatusEnabled

		// 检查是否有未使用的备用码
		backupCount, err := dao.Admin2faBackupCode.Ctx(ctx).Where(g.Map{
			dao.Admin2faBackupCode.Columns().AdminId: adminId,
			dao.Admin2faBackupCode.Columns().IsUsed:  consts.StatusDisable,
		}).Count()
		if err != nil {
			return nil, err
		}
		res.HasBackupCodes = backupCount > 0
	}

	return res, nil
}

// BackupCodes 获取备用码
func (s *sAdmin2fa) BackupCodes(ctx context.Context, in *adminin.Admin2faBackupCodesInp) (res *adminin.Admin2faBackupCodesModel, err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 获取未使用的备用码
	backupCodes, err := dao.Admin2faBackupCode.Ctx(ctx).Where(g.Map{
		dao.Admin2faBackupCode.Columns().AdminId: adminId,
		dao.Admin2faBackupCode.Columns().IsUsed:  consts.StatusDisable,
	}).OrderAsc(dao.Admin2faBackupCode.Columns().CreatedAt).All()
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(backupCodes))
	for _, code := range backupCodes {
		codes = append(codes, gconv.String(code["code"]))
	}

	return &adminin.Admin2faBackupCodesModel{
		BackupCodes: codes,
	}, nil
}

// RegenerateBackupCodes 重新生成备用码
func (s *sAdmin2fa) RegenerateBackupCodes(ctx context.Context, in *adminin.Admin2faRegenerateBackupCodesInp) (res *adminin.Admin2faBackupCodesModel, err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	// 验证当前验证码
	verifyInp := &adminin.Admin2faVerifyInp{
		Code: in.Code,
	}
	if err = s.Verify(ctx, verifyInp); err != nil {
		return nil, err
	}

	// 生成新的备用恢复码
	backupCodes, err := s.generateBackupCodes(8)
	if err != nil {
		return nil, gerror.Wrap(err, "生成备用恢复码失败")
	}

	// 开启事务
	err = dao.Admin2faBackupCode.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除旧的备用码
		_, err = dao.Admin2faBackupCode.Ctx(ctx).TX(tx).Where(dao.Admin2faBackupCode.Columns().AdminId, adminId).Delete()
		if err != nil {
			return err
		}

		// 保存新的备用恢复码
		for _, code := range backupCodes {
			_, err = dao.Admin2faBackupCode.Ctx(ctx).TX(tx).Insert(&entity.Admin2faBackupCode{
				AdminId:   adminId,
				Code:      code,
				IsUsed:    consts.StatusDisable,
				CreatedAt: gtime.Now(),
			})
			if err != nil {
				return err
			}
		}

		// 记录操作日志
		return s.logAction(ctx, tx, adminId, "regenerate_backup", "success", "重新生成备用码")
	})

	if err != nil {
		return nil, err
	}

	return &adminin.Admin2faBackupCodesModel{
		BackupCodes: backupCodes,
	}, nil
}

// UseBackupCode 使用备用码
func (s *sAdmin2fa) UseBackupCode(ctx context.Context, in *adminin.Admin2faUseBackupCodeInp) (err error) {
	adminId := in.AdminId
	if adminId <= 0 {
		return gerror.New("管理员ID无效")
	}

	// 查找备用码
	backupCode, err := dao.Admin2faBackupCode.Ctx(ctx).Where(g.Map{
		dao.Admin2faBackupCode.Columns().AdminId: adminId,
		dao.Admin2faBackupCode.Columns().Code:    in.Code,
		dao.Admin2faBackupCode.Columns().IsUsed:  consts.StatusDisable,
	}).One()
	if err != nil {
		return err
	}
	if backupCode.IsEmpty() {
		// 记录失败日志
		s.logAction(ctx, nil, adminId, "backup_used", "failed", "备用码无效")
		return gerror.New("备用码无效或已使用")
	}

	// 标记备用码为已使用
	_, err = dao.Admin2faBackupCode.Ctx(ctx).Where(dao.Admin2faBackupCode.Columns().Id, gconv.Int64(backupCode["id"])).Update(g.Map{
		dao.Admin2faBackupCode.Columns().IsUsed: consts.StatusEnabled,
		dao.Admin2faBackupCode.Columns().UsedAt: gtime.Now(),
	})
	if err != nil {
		return err
	}

	// 记录成功日志
	return s.logAction(ctx, nil, adminId, "backup_used", "success", "使用备用码验证成功")
}

// LogList 获取2FA日志列表
func (s *sAdmin2fa) LogList(ctx context.Context, in *adminin.Admin2faLogListInp) (res *adminin.Admin2faLogListModel, err error) {
	adminId := contexts.GetUserId(ctx)
	if adminId <= 0 {
		return nil, gerror.New("用户未登录")
	}

	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}

	// 构建查询条件
	where := g.Map{
		dao.Admin2faLog.Columns().AdminId: adminId,
	}
	if in.Action != "" {
		where[dao.Admin2faLog.Columns().Action] = in.Action
	}
	if in.Result != "" {
		where[dao.Admin2faLog.Columns().Result] = in.Result
	}

	// 获取总数
	total, err := dao.Admin2faLog.Ctx(ctx).Where(where).Count()
	if err != nil {
		return nil, err
	}

	// 获取列表
	logs, err := dao.Admin2faLog.Ctx(ctx).Where(where).
		OrderDesc(dao.Admin2faLog.Columns().CreatedAt).
		Page(in.Page, in.PageSize).
		All()
	if err != nil {
		return nil, err
	}

	list := make([]*entity.Admin2faLog, 0, len(logs))
	for _, log := range logs {
		var logEntity entity.Admin2faLog
		if err = gconv.Struct(log, &logEntity); err != nil {
			return nil, err
		}
		list = append(list, &logEntity)
	}

	return &adminin.Admin2faLogListModel{
		List:     list,
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    int64(total),
	}, nil
}

// IsEnabled 检查用户是否启用了2FA
func (s *sAdmin2fa) IsEnabled(ctx context.Context, adminId uint64) (bool, error) {
	if adminId <= 0 {
		return false, nil
	}

	config, err := dao.Admin2faConfig.Ctx(ctx).Where(g.Map{
		dao.Admin2faConfig.Columns().AdminId:    adminId,
		dao.Admin2faConfig.Columns().IsEnabled:  consts.StatusEnabled,
		dao.Admin2faConfig.Columns().IsVerified: consts.StatusEnabled,
	}).One()
	if err != nil {
		return false, err
	}

	return !config.IsEmpty(), nil
}

// generateBackupCodes 生成备用恢复码
func (s *sAdmin2fa) generateBackupCodes(count int) ([]string, error) {
	codes := make([]string, count)
	for i := 0; i < count; i++ {
		// 生成8位随机码
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		codes[i] = strings.ToUpper(hex.EncodeToString(bytes))
	}
	return codes, nil
}

// logAction 记录操作日志
func (s *sAdmin2fa) logAction(ctx context.Context, tx gdb.TX, adminId int64, action, result, remark string) error {
	// 获取请求信息
	ip := "unknown"
	userAgent := "unknown"

	// 尝试从上下文中获取HTTP请求信息
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
		userAgent = r.Header.Get("User-Agent")
	}

	logData := &entity.Admin2faLog{
		AdminId:   adminId,
		Action:    action,
		Result:    result,
		Ip:        ip,
		UserAgent: userAgent,
		Remark:    remark,
		CreatedAt: gtime.Now(),
	}

	if tx != nil {
		_, err := dao.Admin2faLog.Ctx(ctx).TX(tx).Insert(logData)
		return err
	} else {
		_, err := dao.Admin2faLog.Ctx(ctx).Insert(logData)
		return err
	}
}
