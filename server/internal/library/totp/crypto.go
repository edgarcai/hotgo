// Package totp
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
//
package totp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// CryptoService 加密服务
type CryptoService struct {
	key []byte // 32字节AES-256密钥
}

// NewCryptoService 创建加密服务实例
func NewCryptoService(key string) (*CryptoService, error) {
	// 确保密钥长度为32字节（AES-256）
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("密钥长度必须为32字节，当前长度: %d", len(keyBytes))
	}
	
	return &CryptoService{
		key: keyBytes,
	}, nil
}

// EncryptSecret 加密TOTP密钥
// 使用AES-256-GCM加密算法
func (c *CryptoService) EncryptSecret(secret string) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("创建AES密码器失败: %v", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %v", err)
	}
	
	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("生成nonce失败: %v", err)
	}
	
	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(secret), nil)
	
	// 返回base64编码的结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret 解密TOTP密钥
func (c *CryptoService) DecryptSecret(encryptedSecret string) (string, error) {
	// 解码base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		return "", fmt.Errorf("base64解码失败: %v", err)
	}
	
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", fmt.Errorf("创建AES密码器失败: %v", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("创建GCM模式失败: %v", err)
	}
	
	// 检查密文长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("密文长度不足")
	}
	
	// 提取nonce和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}
	
	return string(plaintext), nil
}

// HashBackupCode 哈希备用恢复码
// 使用bcrypt算法进行哈希
func (c *CryptoService) HashBackupCode(code string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("生成哈希失败: %v", err)
	}
	return string(hash), nil
}

// VerifyBackupCode 验证备用恢复码
func (c *CryptoService) VerifyBackupCode(code, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(code))
	return err == nil
}

// GenerateEncryptionKey 生成32字节的加密密钥
// 用于初始化配置
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", fmt.Errorf("生成密钥失败: %v", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}