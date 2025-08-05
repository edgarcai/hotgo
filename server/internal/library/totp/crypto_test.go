package totp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncryptDecryptSecret 测试密钥加密和解密
func TestEncryptDecryptSecret(t *testing.T) {
	// 使用32字节测试密钥
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	// 测试数据
	originalSecret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	// 加密
	encrypted, err := service.EncryptSecret(originalSecret)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	assert.NotEqual(t, originalSecret, encrypted)

	// 解密
	decrypted, err := service.DecryptSecret(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, originalSecret, decrypted)
}

// TestEncryptSecretDifferentResults 测试相同输入产生不同加密结果
func TestEncryptSecretDifferentResults(t *testing.T) {
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	// 多次加密相同内容应该产生不同结果（由于随机nonce）
	encrypted1, err := service.EncryptSecret(secret)
	assert.NoError(t, err)

	encrypted2, err := service.EncryptSecret(secret)
	assert.NoError(t, err)

	assert.NotEqual(t, encrypted1, encrypted2)

	// 但解密后应该得到相同结果
	decrypted1, err := service.DecryptSecret(encrypted1)
	assert.NoError(t, err)

	decrypted2, err := service.DecryptSecret(encrypted2)
	assert.NoError(t, err)

	assert.Equal(t, secret, decrypted1)
	assert.Equal(t, secret, decrypted2)
}

// TestHashVerifyBackupCode 测试备用码哈希和验证
func TestHashVerifyBackupCode(t *testing.T) {
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	// 测试备用码
	backupCode := "12345678"

	// 哈希
	hashed, err := service.HashBackupCode(backupCode)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, backupCode, hashed)

	// 验证正确的备用码
	valid := service.VerifyBackupCode(backupCode, hashed)
	assert.True(t, valid)

	// 验证错误的备用码
	valid = service.VerifyBackupCode("87654321", hashed)
	assert.False(t, valid)

	// 验证空备用码
	valid = service.VerifyBackupCode("", hashed)
	assert.False(t, valid)
}

// TestHashBackupCodeDifferentResults 测试相同备用码产生不同哈希
func TestHashBackupCodeDifferentResults(t *testing.T) {
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	backupCode := "12345678"

	// 多次哈希相同内容应该产生不同结果（由于随机salt）
	hashed1, err := service.HashBackupCode(backupCode)
	assert.NoError(t, err)

	hashed2, err := service.HashBackupCode(backupCode)
	assert.NoError(t, err)

	assert.NotEqual(t, hashed1, hashed2)

	// 但都应该能验证原始备用码
	valid1 := service.VerifyBackupCode(backupCode, hashed1)
	assert.True(t, valid1)

	valid2 := service.VerifyBackupCode(backupCode, hashed2)
	assert.True(t, valid2)
}

// TestInvalidEncryptedData 测试无效加密数据的处理
func TestInvalidEncryptedData(t *testing.T) {
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	// 测试解密无效数据
	_, err = service.DecryptSecret("invalid-encrypted-data")
	assert.Error(t, err)

	// 测试解密空数据
	_, err = service.DecryptSecret("")
	assert.Error(t, err)
}

// TestEmptyInputs 测试空输入的处理
func TestEmptyInputs(t *testing.T) {
	testKey := "12345678901234567890123456789012"
	service, err := NewCryptoService(testKey)
	assert.NoError(t, err)

	// 测试加密空字符串
	encrypted, err := service.EncryptSecret("")
	assert.NoError(t, err) // 应该允许加密空字符串

	decrypted, err := service.DecryptSecret(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "", decrypted)

	// 测试哈希空备用码
	hashed, err := service.HashBackupCode("")
	assert.NoError(t, err) // 应该允许哈希空字符串

	valid := service.VerifyBackupCode("", hashed)
	assert.True(t, valid)
}

// BenchmarkEncryptSecret 性能测试：密钥加密
func BenchmarkEncryptSecret(b *testing.B) {
	testKey := "12345678901234567890123456789012"
	service, _ := NewCryptoService(testKey)
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.EncryptSecret(secret)
	}
}

// BenchmarkDecryptSecret 性能测试：密钥解密
func BenchmarkDecryptSecret(b *testing.B) {
	testKey := "12345678901234567890123456789012"
	service, _ := NewCryptoService(testKey)
	secret := "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ"
	encrypted, _ := service.EncryptSecret(secret)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.DecryptSecret(encrypted)
	}
}

// BenchmarkHashBackupCode 性能测试：备用码哈希
func BenchmarkHashBackupCode(b *testing.B) {
	testKey := "12345678901234567890123456789012"
	service, _ := NewCryptoService(testKey)
	backupCode := "12345678"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.HashBackupCode(backupCode)
	}
}

// BenchmarkVerifyBackupCode 性能测试：备用码验证
func BenchmarkVerifyBackupCode(b *testing.B) {
	testKey := "12345678901234567890123456789012"
	service, _ := NewCryptoService(testKey)
	backupCode := "12345678"
	hashed, _ := service.HashBackupCode(backupCode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.VerifyBackupCode(backupCode, hashed)
	}
}
