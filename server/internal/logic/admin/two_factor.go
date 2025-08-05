// Package admin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package admin

import (
	"context"
	"fmt"

	"hotgo/internal/dao"
	"hotgo/internal/library/totp"
	"hotgo/internal/model/do"
	"hotgo/internal/model/input/adminin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sTwoFactor struct {
	totpService   *totp.Service
	cryptoService *totp.CryptoService
}

func NewTwoFactor() *sTwoFactor {
	// 从配置中获取加密密钥
	encryptionKey := g.Cfg().MustGet(context.Background(), "totp.encryptionKey", "your-32-byte-encryption-key-here").String()

	totpService := totp.NewService("HotGo")
	cryptoService, err := totp.NewCryptoService(encryptionKey)
	if err != nil {
		g.Log().Fatal(context.Background(), "Failed to create crypto service:", err)
	}

	return &sTwoFactor{
		totpService:   totpService,
		cryptoService: cryptoService,
	}
}

func init() {
	service.RegisterAdminTwoFactor(NewTwoFactor())
}

// Enable 启用双因子认证
func (s *sTwoFactor) Enable(ctx context.Context, in *adminin.TwoFactorEnableInp) (res *adminin.TwoFactorEnableModel, err error) {
	// 生成TOTP密钥
	accountName := fmt.Sprintf("admin_%d", in.UserId)
	key, err := s.totpService.GenerateSecret(accountName)
	if err != nil {
		return nil, err
	}

	// 加密密钥
	encryptedSecret, err := s.cryptoService.EncryptSecret(key.Secret())
	if err != nil {
		return nil, err
	}

	// 生成二维码URL
	qrCodeURL := s.totpService.GetQRCodeURL(key)

	// 格式化密钥用于显示
	formattedSecret := s.totpService.FormatSecret(key.Secret())

	// 生成备用恢复码
	backupCodes := s.totpService.GenerateBackupCodes(8)

	// 哈希备用恢复码
	hashedBackupCodes := make([]string, len(backupCodes))
	for i, code := range backupCodes {
		hashedCode, err := s.cryptoService.HashBackupCode(code)
		if err != nil {
			return nil, err
		}
		hashedBackupCodes[i] = hashedCode
	}

	// 保存到数据库（临时状态，未启用）
	_, err = dao.AdminTwoFactor.Ctx(ctx).Data(g.Map{
		"member_id":    in.UserId,
		"secret_key":   encryptedSecret,
		"backup_codes": gjson.New(hashedBackupCodes),
		"is_enabled":   0, // 未启用
		"created_at":   gtime.Now(),
		"updated_at":   gtime.Now(),
	}).OnDuplicate("secret_key", "backup_codes", "updated_at").Save()
	if err != nil {
		return nil, err
	}

	res = &adminin.TwoFactorEnableModel{
		Secret:      formattedSecret,
		QRCodeURL:   qrCodeURL,
		BackupCodes: backupCodes,
	}
	return
}

// ConfirmEnable 确认启用双因子认证
func (s *sTwoFactor) ConfirmEnable(ctx context.Context, in *adminin.TwoFactorConfirmEnableInp) (err error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).One()
	if err != nil {
		return err
	}
	if record.IsEmpty() {
		return gerror.New("双因子认证未初始化")
	}

	// 检查是否已启用
	if record["is_enabled"].Bool() {
		return gerror.New("双因子认证已启用")
	}

	// 解密密钥
	encryptedSecret := record["secret_key"].String()
	secret, err := s.cryptoService.DecryptSecret(encryptedSecret)
	if err != nil {
		return err
	}

	// 验证TOTP代码
	valid := s.totpService.ValidateCode(secret, in.Code)
	if !valid {
		return gerror.New("验证码错误")
	}

	// 启用2FA
	_, err = dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).Data(do.AdminTwoFactor{
		IsEnabled: 1,
		UpdatedAt: gtime.Now(),
	}).Update()

	return err
}

// Disable 禁用双因子认证
func (s *sTwoFactor) Disable(ctx context.Context, in *adminin.TwoFactorDisableInp) (err error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).One()
	if err != nil {
		return err
	}
	if record.IsEmpty() {
		return gerror.New("双因子认证未启用")
	}

	// 解密密钥
	encryptedSecret := record["secret_key"].String()
	secret, err := s.cryptoService.DecryptSecret(encryptedSecret)
	if err != nil {
		return err
	}

	// 验证TOTP代码或备用恢复码
	var isValid bool
	if len(in.Code) == 6 {
		// TOTP代码
		isValid = s.totpService.ValidateCode(secret, in.Code)
	} else {
		// 备用恢复码
		isValid, err = s.verifyBackupCode(ctx, in.UserId, in.Code)
		if err != nil {
			return err
		}
	}

	if !isValid {
		return gerror.New("验证码错误")
	}

	// 删除2FA记录
	_, err = dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).Delete()

	return err
}

// Verify 验证双因子认证
func (s *sTwoFactor) Verify(ctx context.Context, in *adminin.TwoFactorVerifyInp) (err error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).Where("is_enabled", 1).One()
	if err != nil {
		return err
	}
	if record.IsEmpty() {
		return gerror.New("双因子认证未启用")
	}

	// 解密密钥
	encryptedSecret := record["secret_key"].String()
	secret, err := s.cryptoService.DecryptSecret(encryptedSecret)
	if err != nil {
		return err
	}

	// 验证TOTP代码或备用恢复码
	var isValid bool
	if len(in.Code) == 6 {
		// TOTP代码
		isValid = s.totpService.ValidateCode(secret, in.Code)
	} else {
		// 备用恢复码
		isValid, err = s.verifyAndUseBackupCode(ctx, in.UserId, in.Code)
		if err != nil {
			return err
		}
	}

	if !isValid {
		return gerror.New("验证码错误")
	}

	return nil
}

// GetStatus 获取双因子认证状态
func (s *sTwoFactor) GetStatus(ctx context.Context, in *adminin.TwoFactorGetStatusInp) (res *adminin.TwoFactorGetStatusModel, err error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).One()
	if err != nil {
		return nil, err
	}

	// 获取剩余备用恢复码数量
	backupCodesCount := 0
	if !record.IsEmpty() {
		// 从backup_codes字段解析剩余数量
		backupCodesCount = len(record["backup_codes"].Strings())
	}

	res = &adminin.TwoFactorGetStatusModel{
		IsEnabled:        !record.IsEmpty() && record["is_enabled"].Bool(),
		BackupCodesCount: backupCodesCount,
	}

	return res, nil
}

// RegenerateBackupCodes 重新生成备用恢复码
func (s *sTwoFactor) RegenerateBackupCodes(ctx context.Context, in *adminin.TwoFactorRegenerateBackupCodesInp) (res *adminin.TwoFactorRegenerateBackupCodesModel, err error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).Where("is_enabled", 1).One()
	if err != nil {
		return nil, err
	}
	if record.IsEmpty() {
		return nil, gerror.New("双因子认证未启用")
	}

	// 解密密钥
	encryptedSecret := record["secret_key"].String()
	secret, err := s.cryptoService.DecryptSecret(encryptedSecret)
	if err != nil {
		return nil, err
	}

	// 验证TOTP代码
	valid := s.totpService.ValidateCode(secret, in.Code)
	if !valid {
		return nil, gerror.New("验证码错误")
	}

	// 生成新的备用恢复码
	backupCodes := s.totpService.GenerateBackupCodes(8)

	// 哈希备用恢复码
	hashedCodes := make([]string, len(backupCodes))
	for i, code := range backupCodes {
		hashedCodes[i], err = s.cryptoService.HashBackupCode(code)
		if err != nil {
			return nil, err
		}
	}

	// 更新备用恢复码
	_, err = dao.AdminTwoFactor.Ctx(ctx).Where("member_id", in.UserId).Data(do.AdminTwoFactor{
		BackupCodes: gjson.New(hashedCodes),
		UpdatedAt:   gtime.Now(),
	}).Update()

	if err != nil {
		return nil, err
	}

	res = &adminin.TwoFactorRegenerateBackupCodesModel{
		BackupCodes: backupCodes,
	}

	return res, nil
}

// verifyBackupCode 验证备用恢复码（不标记为已使用）
func (s *sTwoFactor) verifyBackupCode(ctx context.Context, userId int64, code string) (bool, error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", userId).One()
	if err != nil {
		return false, err
	}
	if record.IsEmpty() {
		return false, nil
	}

	// 获取备用恢复码
	backupCodes := record["backup_codes"].Strings()

	// 验证备用恢复码
	for _, hash := range backupCodes {
		if s.cryptoService.VerifyBackupCode(code, hash) {
			return true, nil
		}
	}

	return false, nil
}

// verifyAndUseBackupCode 验证并使用备用恢复码
func (s *sTwoFactor) verifyAndUseBackupCode(ctx context.Context, userId int64, code string) (bool, error) {
	// 获取用户的2FA记录
	record, err := dao.AdminTwoFactor.Ctx(ctx).Where("member_id", userId).One()
	if err != nil {
		return false, err
	}
	if record.IsEmpty() {
		return false, nil
	}

	// 获取备用恢复码
	backupCodes := record["backup_codes"].Strings()
	updatedCodes := make([]string, 0)

	// 验证备用恢复码
	for _, hash := range backupCodes {
		if s.cryptoService.VerifyBackupCode(code, hash) {
			// 从列表中移除已使用的码
			_, err = dao.AdminTwoFactor.Ctx(ctx).Where("member_id", userId).Data(do.AdminTwoFactor{
				BackupCodes: gjson.New(updatedCodes),
				UpdatedAt:   gtime.Now(),
			}).Update()
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			updatedCodes = append(updatedCodes, hash)
		}
	}

	return false, nil
}
