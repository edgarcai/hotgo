package totp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateSecret 测试密钥生成
func TestGenerateSecret(t *testing.T) {
	service := NewService("TestApp")
	
	// 测试密钥生成
	key, err := service.GenerateSecret("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotEmpty(t, key.Secret())
	assert.Contains(t, key.URL(), "TestApp")
	assert.Contains(t, key.URL(), "testuser")
}

// TestValidateCode 测试验证码验证
func TestValidateCode(t *testing.T) {
	service := NewService("TestApp")
	
	// 生成测试密钥
	key, err := service.GenerateSecret("testuser")
	assert.NoError(t, err)
	
	// 获取当前验证码
	code, err := service.GetCurrentCode(key.Secret())
	assert.NoError(t, err)
	assert.Len(t, code, 6)
	
	// 验证正确的验证码
	valid := service.ValidateCode(key.Secret(), code)
	assert.True(t, valid)
	
	// 验证错误的验证码
	valid = service.ValidateCode(key.Secret(), "000000")
	assert.False(t, valid)
}

// TestGenerateBackupCodes 测试备用码生成
func TestGenerateBackupCodes(t *testing.T) {
	service := NewService("TestApp")
	
	// 生成10个备用码
	codes := service.GenerateBackupCodes(10)
	assert.Len(t, codes, 10)
	
	// 检查每个备用码的格式
	for _, code := range codes {
		assert.Len(t, code, 8) // 8位备用码
		assert.Regexp(t, "^[A-Z0-9]+$", code) // 只包含大写字母和数字
	}
	
	// 确保所有备用码都是唯一的
	uniqueMap := make(map[string]bool)
	for _, code := range codes {
		assert.False(t, uniqueMap[code], "备用码应该是唯一的")
		uniqueMap[code] = true
	}
}

// TestGetQRCodeURL 测试二维码URL生成
func TestGetQRCodeURL(t *testing.T) {
	service := NewService("TestApp")
	
	key, err := service.GenerateSecret("testuser")
	assert.NoError(t, err)
	
	url := service.GetQRCodeURL(key)
	assert.NotEmpty(t, url)
	assert.Equal(t, key.URL(), url)
}

// TestFormatSecret 测试密钥格式化
func TestFormatSecret(t *testing.T) {
	service := NewService("TestApp")
	
	// 测试正常密钥
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	formatted := service.FormatSecret(secret)
	assert.Equal(t, "GEZD GNBV GY3T QOJQ GEZD GNBV GY3T QOJQ", formatted)
	
	// 测试空密钥
	formatted = service.FormatSecret("")
	assert.Equal(t, "", formatted)
}

// TestValidateSecret 测试密钥验证
func TestValidateSecret(t *testing.T) {
	service := NewService("TestApp")
	
	// 测试有效密钥
	validSecret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	err := service.ValidateSecret(validSecret)
	assert.NoError(t, err)
	
	// 测试无效密钥（包含无效字符）
	invalidSecret := "INVALID123!"
	err = service.ValidateSecret(invalidSecret)
	assert.Error(t, err)
	
	// 测试空密钥
	err = service.ValidateSecret("")
	assert.Error(t, err)
}

// TestTimeWindow 测试时间窗口容错
func TestTimeWindow(t *testing.T) {
	service := NewService("TestApp")
	
	key, err := service.GenerateSecret("testuser")
	assert.NoError(t, err)
	
	// 获取当前时间的验证码
	code, err := service.GetCurrentCode(key.Secret())
	assert.NoError(t, err)
	
	// 当前验证码应该有效
	valid := service.ValidateCode(key.Secret(), code)
	assert.True(t, valid)
}

// BenchmarkGenerateSecret 性能测试：密钥生成
func BenchmarkGenerateSecret(b *testing.B) {
	service := NewService("TestApp")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateSecret("testuser")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValidateCode 性能测试：验证码验证
func BenchmarkValidateCode(b *testing.B) {
	service := NewService("TestApp")
	key, _ := service.GenerateSecret("testuser")
	code, _ := service.GetCurrentCode(key.Secret())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateCode(key.Secret(), code)
	}
}