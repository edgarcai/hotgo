package totp

import (
	"encoding/base32"
	"strings"
	"testing"
	"time"
)

// TestGenerateSecret 测试密钥生成
func TestGenerateSecret(t *testing.T) {
	config := DefaultConfig()
	config.SecretSize = 20

	secret, err := GenerateSecret(config)
	if err != nil {
		t.Fatalf("生成密钥失败: %v", err)
	}

	if secret == "" {
		t.Fatal("生成的密钥为空")
	}

	// 验证Base32格式
	if err := ValidateSecret(secret); err != nil {
		t.Fatalf("生成的密钥格式无效: %v", err)
	}

	t.Logf("生成的密钥: %s", secret)
}

// TestGenerateQRCodeURL 测试二维码URL生成
func TestGenerateQRCodeURL(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	issuer := "HotGo"
	accountName := "test@example.com"

	url := GenerateQRCodeURL(secret, issuer, accountName)

	if !strings.HasPrefix(url, "otpauth://totp/") {
		t.Fatal("二维码URL格式错误")
	}

	if !strings.Contains(url, secret) {
		t.Fatal("二维码URL中缺少密钥")
	}

	if !strings.Contains(url, issuer) {
		t.Fatal("二维码URL中缺少发行者")
	}

	t.Logf("生成的二维码URL: %s", url)
}

// TestValidateCode 测试验证码验证
func TestValidateCode(t *testing.T) {
	// 使用固定密钥进行测试
	secret := "JBSWY3DPEHPK3PXP"
	config := DefaultConfig()

	// 生成当前时间的验证码
	code, err := GenerateCode(secret, config)
	if err != nil {
		t.Fatalf("生成验证码失败: %v", err)
	}

	// 验证生成的验证码
	if !ValidateCode(secret, code, config) {
		t.Fatal("验证码验证失败")
	}

	// 测试错误的验证码
	if ValidateCode(secret, "000000", config) {
		t.Fatal("错误的验证码通过了验证")
	}

	t.Logf("当前验证码: %s", code)
}

// TestValidateSecret 测试密钥验证
func TestValidateSecret(t *testing.T) {
	// 测试有效密钥
	validSecret := "JBSWY3DPEHPK3PXP"
	if err := ValidateSecret(validSecret); err != nil {
		t.Fatalf("有效密钥验证失败: %v", err)
	}

	// 测试无效密钥
	invalidSecrets := []string{
		"",          // 空密钥
		"invalid",   // 无效Base32
		"123456789", // 无效字符
	}

	for _, secret := range invalidSecrets {
		if err := ValidateSecret(secret); err == nil {
			t.Fatalf("无效密钥 '%s' 通过了验证", secret)
		}
	}
}

// TestGetRemainingTime 测试剩余时间计算
func TestGetRemainingTime(t *testing.T) {
	remaining := GetRemainingTime(30)

	if remaining < 0 || remaining > 30 {
		t.Fatalf("剩余时间计算错误: %d", remaining)
	}

	t.Logf("当前验证码剩余时间: %d秒", remaining)
}

// TestTimeWindowTolerance 测试时间窗口容错
func TestTimeWindowTolerance(t *testing.T) {
	secret := "JBSWY3DPEHPK3PXP"
	config := DefaultConfig()

	// 模拟不同时间点的验证码
	now := time.Now().Unix()
	timeStep := now / int64(config.Period)

	// 测试当前时间步的验证码
	currentCode := generateCode(decodeSecret(secret), timeStep, config.Digits)
	if !ValidateCode(secret, currentCode, config) {
		t.Fatal("当前时间步验证码验证失败")
	}

	// 测试前一个时间步的验证码
	prevCode := generateCode(decodeSecret(secret), timeStep-1, config.Digits)
	if !ValidateCode(secret, prevCode, config) {
		t.Fatal("前一个时间步验证码验证失败")
	}

	// 测试后一个时间步的验证码
	nextCode := generateCode(decodeSecret(secret), timeStep+1, config.Digits)
	if !ValidateCode(secret, nextCode, config) {
		t.Fatal("后一个时间步验证码验证失败")
	}

	// 测试超出容错范围的验证码
	oldCode := generateCode(decodeSecret(secret), timeStep-2, config.Digits)
	if ValidateCode(secret, oldCode, config) {
		t.Fatal("超出容错范围的验证码通过了验证")
	}
}

// decodeSecret 辅助函数：解码Base32密钥
func decodeSecret(secret string) []byte {
	secretBytes, _ := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	return secretBytes
}

// BenchmarkValidateCode 验证码验证性能测试
func BenchmarkValidateCode(b *testing.B) {
	secret := "JBSWY3DPEHPK3PXP"
	config := DefaultConfig()
	code, _ := GenerateCode(secret, config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateCode(secret, code, config)
	}
}

// BenchmarkGenerateSecret 密钥生成性能测试
func BenchmarkGenerateSecret(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateSecret(config)
	}
}
