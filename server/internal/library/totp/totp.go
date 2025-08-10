// Package totp 实现基于时间的一次性密码算法(TOTP)
// 符合RFC 6238标准，支持Google Authenticator等主流应用
package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// TOTPConfig TOTP配置结构
type TOTPConfig struct {
	Issuer      string // 发行者名称
	AccountName string // 账户名称
	SecretSize  int    // 密钥长度(字节)
	Period      int    // 时间步长(秒)
	Digits      int    // 验证码位数
	Algorithm   string // 哈希算法
}

// DefaultConfig 返回默认TOTP配置
func DefaultConfig() *TOTPConfig {
	return &TOTPConfig{
		Issuer:      "HotGo",
		AccountName: "",
		SecretSize:  20,
		Period:      30,
		Digits:      6,
		Algorithm:   "SHA1",
	}
}

// GenerateSecret 生成TOTP密钥
// 返回Base32编码的密钥字符串
func GenerateSecret(config *TOTPConfig) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	if config.SecretSize <= 0 {
		config.SecretSize = 20
	}
	
	secret := make([]byte, config.SecretSize)
	_, err := rand.Read(secret)
	if err != nil {
		return "", gerror.Wrap(err, "生成随机密钥失败")
	}
	
	return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateQRCodeURL 生成二维码URL
// 用于身份验证应用扫描绑定
func GenerateQRCodeURL(secret, issuer, accountName string) string {
	if issuer == "" {
		issuer = "HotGo"
	}
	
	// 构建otpauth URL
	v := url.Values{}
	v.Set("secret", secret)
	v.Set("issuer", issuer)
	v.Set("algorithm", "SHA1")
	v.Set("digits", "6")
	v.Set("period", "30")
	
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?%s",
		url.QueryEscape(issuer),
		url.QueryEscape(accountName),
		v.Encode(),
	)
}

// ValidateCode 验证TOTP验证码
// 支持时间窗口容错，允许前后一个时间步长的验证码
func ValidateCode(secret, code string, config *TOTPConfig) bool {
	if config == nil {
		config = DefaultConfig()
	}
	
	// 解码Base32密钥
	secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return false
	}
	
	// 获取当前时间步
	now := time.Now().Unix()
	timeStep := now / int64(config.Period)
	
	// 允许前后一个时间窗口的验证码(容错机制)
	for i := -1; i <= 1; i++ {
		expectedCode := generateCode(secretBytes, timeStep+int64(i), config.Digits)
		if expectedCode == code {
			return true
		}
	}
	
	return false
}

// GenerateCode 生成当前时间的TOTP验证码
// 主要用于测试和调试
func GenerateCode(secret string, config *TOTPConfig) (string, error) {
	if config == nil {
		config = DefaultConfig()
	}
	
	secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", gerror.Wrap(err, "解码密钥失败")
	}
	
	now := time.Now().Unix()
	timeStep := now / int64(config.Period)
	
	return generateCode(secretBytes, timeStep, config.Digits), nil
}

// generateCode 生成指定时间步的验证码
// 实现TOTP算法核心逻辑
func generateCode(secret []byte, timeStep int64, digits int) string {
	// 将时间步转换为8字节大端序
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeStep))
	
	// 使用HMAC-SHA1计算哈希
	mac := hmac.New(sha1.New, secret)
	mac.Write(buf)
	hash := mac.Sum(nil)
	
	// 动态截取(Dynamic Truncation)
	offset := hash[len(hash)-1] & 0x0F
	code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF
	
	// 取模得到指定位数的验证码
	code = code % uint32(math.Pow10(digits))
	
	// 格式化为指定位数的字符串，不足位数前补0
	return fmt.Sprintf("%0*d", digits, code)
}

// GetRemainingTime 获取当前验证码剩余有效时间(秒)
func GetRemainingTime(period int) int {
	if period <= 0 {
		period = 30
	}
	
	now := time.Now().Unix()
	return period - int(now%int64(period))
}

// ValidateSecret 验证密钥格式是否正确
func ValidateSecret(secret string) error {
	if secret == "" {
		return gerror.New("密钥不能为空")
	}
	
	// 检查Base32格式
	_, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return gerror.Wrap(err, "密钥格式错误，必须是有效的Base32字符串")
	}
	
	return nil
}