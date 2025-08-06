// Package totp
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package totp

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// Service TOTP服务
type Service struct {
	issuer string // 发行者名称
}

// NewService 创建TOTP服务实例
func NewService(issuer string) *Service {
	return &Service{
		issuer: issuer,
	}
}

// GenerateSecret 生成TOTP密钥
// 返回密钥对象，包含密钥字符串和二维码URL
func (s *Service) GenerateSecret(username string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      s.issuer,
		AccountName: username,
		SecretSize:  32,                // 32字节密钥长度
		Digits:      otp.DigitsSix,     // 6位验证码
		Algorithm:   otp.AlgorithmSHA1, // SHA1算法
		Period:      30,                // 30秒时间窗口
	})
}

// ValidateCode 验证TOTP代码
// 支持时钟偏移容忍（前后各1个时间窗口）
func (s *Service) ValidateCode(secret, code string) bool {
	// 使用内置的验证方法，它已经包含了时钟偏移容忍
	return totp.Validate(code, secret)
}

// GenerateBackupCodes 生成备用恢复码
// 生成指定数量的8位随机恢复码
func (s *Service) GenerateBackupCodes(count int) []string {
	codes := make([]string, count)
	for i := range codes {
		codes[i] = s.generateRandomCode(8)
	}
	return codes
}

// generateRandomCode 生成随机恢复码
func (s *Service) generateRandomCode(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	// 使用base32编码并去除填充字符
	encoded := base32.StdEncoding.EncodeToString(bytes)
	encoded = strings.TrimRight(encoded, "=")

	// 截取指定长度
	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return strings.ToUpper(encoded)
}

// GetQRCodeURL 获取二维码URL
func (s *Service) GetQRCodeURL(key *otp.Key) string {
	return key.URL()
}

// GenerateQRCodeImage 生成二维码图片（base64编码）
// 返回 data:image/png;base64,... 格式的字符串，可直接用于 img 标签的 src 属性
func (s *Service) GenerateQRCodeImage(key *otp.Key) (string, error) {
	// 生成二维码图片字节数据
	qrBytes, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %v", err)
	}

	// 转换为base64编码
	base64Str := base64.StdEncoding.EncodeToString(qrBytes)

	// 返回完整的data URL格式
	return fmt.Sprintf("data:image/png;base64,%s", base64Str), nil
}

// FormatSecret 格式化密钥显示
// 将密钥按4个字符分组，便于用户手动输入
func (s *Service) FormatSecret(secret string) string {
	var formatted strings.Builder
	for i, char := range secret {
		if i > 0 && i%4 == 0 {
			formatted.WriteString(" ")
		}
		formatted.WriteRune(char)
	}
	return formatted.String()
}

// ValidateSecret 验证密钥格式
func (s *Service) ValidateSecret(secret string) error {
	// 移除空格
	secret = strings.ReplaceAll(secret, " ", "")

	// 检查长度（base32编码的32字节密钥应该是52个字符）
	if len(secret) < 16 {
		return fmt.Errorf("密钥长度不足")
	}

	// 检查是否为有效的base32编码
	if _, err := base32.StdEncoding.DecodeString(secret); err != nil {
		return fmt.Errorf("无效的密钥格式: %v", err)
	}

	return nil
}

// GetCurrentCode 获取当前时间的TOTP代码（用于测试）
func (s *Service) GetCurrentCode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
