# HotGo项目2FA功能实施指南

## 1. 实施概述

本指南详细说明了在HotGo项目中实施二次因子验证(2FA)功能的具体步骤，包括后端服务实现、前端组件开发、数据库迁移、安全配置和测试验证等各个环节。

## 2. 后端实现方案

### 2.1 TOTP算法库实现

```go
// internal/library/totp/totp.go
package totp

import (
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base32"
    "encoding/binary"
    "fmt"
    "math"
    "strings"
    "time"
)

// TOTPConfig TOTP配置
type TOTPConfig struct {
    Issuer      string // 发行者名称
    AccountName string // 账户名称
    SecretSize  int    // 密钥长度
    Period      int    // 时间步长（秒）
    Digits      int    // 验证码位数
}

// GenerateSecret 生成TOTP密钥
func GenerateSecret(config *TOTPConfig) (string, error) {
    secret := make([]byte, config.SecretSize)
    _, err := rand.Read(secret)
    if err != nil {
        return "", err
    }
    return base32.StdEncoding.EncodeToString(secret), nil
}

// GenerateQRCode 生成二维码URL
func GenerateQRCode(secret, issuer, accountName string) string {
    return fmt.Sprintf(
        "otpauth://totp/%s:%s?secret=%s&issuer=%s",
        issuer, accountName, secret, issuer,
    )
}

// ValidateCode 验证TOTP码
func ValidateCode(secret, code string, config *TOTPConfig) bool {
    secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
    if err != nil {
        return false
    }
    
    now := time.Now().Unix()
    timeStep := now / int64(config.Period)
    
    // 允许前后一个时间窗口的验证码
    for i := -1; i <= 1; i++ {
        if generateCode(secretBytes, timeStep+int64(i), config.Digits) == code {
            return true
        }
    }
    return false
}

// generateCode 生成指定时间步的验证码
func generateCode(secret []byte, timeStep int64, digits int) string {
    buf := make([]byte, 8)
    binary.BigEndian.PutUint64(buf, uint64(timeStep))
    
    mac := hmac.New(sha1.New, secret)
    mac.Write(buf)
    hash := mac.Sum(nil)
    
    offset := hash[len(hash)-1] & 0x0F
    code := binary.BigEndian.Uint32(hash[offset:offset+4]) & 0x7FFFFFFF
    code = code % uint32(math.Pow10(digits))
    
    return fmt.Sprintf("%0*d", digits, code)
}
```

### 2.2 2FA服务层实现

```go
// internal/service/admin_2fa.go
package service

import (
    "context"
    "crypto/rand"
    "fmt"
    "hotgo/internal/dao"
    "hotgo/internal/library/totp"
    "hotgo/internal/model/entity"
    "hotgo/internal/model/input/adminin"
    "time"
)

type sAdmin2FA struct{}

func NewAdmin2FA() *sAdmin2FA {
    return &sAdmin2FA{}
}

// GenerateSecret 生成2FA密钥和二维码
func (s *sAdmin2FA) GenerateSecret(ctx context.Context, memberId int64) (*adminin.Generate2FASecretModel, error) {
    config := &totp.TOTPConfig{
        Issuer:      "HotGo",
        AccountName: fmt.Sprintf("user_%d", memberId),
        SecretSize:  20,
        Period:      30,
        Digits:      6,
    }
    
    secret, err := totp.GenerateSecret(config)
    if err != nil {
        return nil, err
    }
    
    qrCodeURL := totp.GenerateQRCode(secret, config.Issuer, config.AccountName)
    qrCodeImage, err := s.generateQRCodeImage(qrCodeURL)
    if err != nil {
        return nil, err
    }
    
    return &adminin.Generate2FASecretModel{
        Secret:    secret,
        QRCode:    qrCodeImage,
        BackupURL: qrCodeURL,
    }, nil
}

// VerifyAndEnable 验证并启用2FA
func (s *sAdmin2FA) VerifyAndEnable(ctx context.Context, in *adminin.VerifyAndEnable2FAInp) (*adminin.VerifyAndEnable2FAModel, error) {
    config := &totp.TOTPConfig{
        Period: 30,
        Digits: 6,
    }
    
    // 验证TOTP码
    if !totp.ValidateCode(in.Secret, in.Code, config) {
        return nil, gerror.New("验证码错误")
    }
    
    // 保存2FA配置
    _, err := dao.Admin2FAConfig.Ctx(ctx).Insert(&entity.Admin2FAConfig{
        MemberId:  in.MemberId,
        SecretKey: in.Secret,
        Enabled:   true,
        EnabledAt: gtime.Now(),
    })
    if err != nil {
        return nil, err
    }
    
    // 生成备用恢复码
    backupCodes, err := s.generateBackupCodes(ctx, in.MemberId)
    if err != nil {
        return nil, err
    }
    
    // 记录操作日志
    s.log2FAAction(ctx, in.MemberId, "enable", true, "")
    
    return &adminin.VerifyAndEnable2FAModel{
        Success:     true,
        BackupCodes: backupCodes,
    }, nil
}

// Verify2FA 验证2FA码
func (s *sAdmin2FA) Verify2FA(ctx context.Context, in *adminin.Verify2FAInp) (*adminin.Verify2FAModel, error) {
    // 获取用户2FA配置
    var config *entity.Admin2FAConfig
    err := dao.Admin2FAConfig.Ctx(ctx).Where("member_id", in.MemberId).Scan(&config)
    if err != nil || config == nil {
        return nil, gerror.New("用户未启用2FA")
    }
    
    var success bool
    var isBackupCode bool
    
    // 检查是否为备用码
    if len(in.Code) > 6 {
        success, err = s.verifyBackupCode(ctx, in.MemberId, in.Code)
        isBackupCode = true
    } else {
        // 验证TOTP码
        totpConfig := &totp.TOTPConfig{
            Period: 30,
            Digits: 6,
        }
        success = totp.ValidateCode(config.SecretKey, in.Code, totpConfig)
    }
    
    // 记录验证日志
    action := "verify_totp"
    if isBackupCode {
        action = "verify_backup_code"
    }
    s.log2FAAction(ctx, in.MemberId, action, success, "")
    
    if !success {
        return nil, gerror.New("验证码错误")
    }
    
    return &adminin.Verify2FAModel{
        Success: true,
    }, nil
}

// generateBackupCodes 生成备用恢复码
func (s *sAdmin2FA) generateBackupCodes(ctx context.Context, memberId int64) ([]string, error) {
    codes := make([]string, 10)
    
    for i := 0; i < 10; i++ {
        code := s.generateRandomCode(8)
        codes[i] = code
        
        // 保存到数据库
        _, err := dao.Admin2FABackupCode.Ctx(ctx).Insert(&entity.Admin2FABackupCode{
            MemberId: memberId,
            Code:     code,
            Used:     false,
        })
        if err != nil {
            return nil, err
        }
    }
    
    return codes, nil
}

// generateRandomCode 生成随机码
func (s *sAdmin2FA) generateRandomCode(length int) string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    rand.Read(b)
    for i := range b {
        b[i] = charset[b[i]%byte(len(charset))]
    }
    return string(b)
}
```

### 2.3 认证中间件扩展

```go
// internal/middleware/auth_2fa.go
package middleware

import (
    "github.com/gogf/gf/v2/net/ghttp"
    "hotgo/internal/service"
)

// Auth2FA 2FA认证中间件
func Auth2FA(r *ghttp.Request) {
    user := contexts.GetUser(r.Context())
    if user == nil {
        response.JsonExit(r, 401, "请先登录")
    }
    
    // 检查用户是否启用了2FA
    has2FA, err := service.Admin2FA().HasEnabled2FA(r.Context(), user.Id)
    if err != nil {
        response.JsonExit(r, 500, "系统错误")
    }
    
    if has2FA {
        // 检查当前会话是否已通过2FA验证
        verified := r.Session.GetBool("2fa_verified")
        if !verified {
            response.JsonExit(r, 403, "需要2FA验证")
        }
    }
    
    r.Middleware.Next()
}
```

## 3. 前端实现方案

### 3.1 2FA设置组件

```vue
<!-- src/components/2FA/Setup2FA.vue -->
<template>
  <n-card title="二次验证设置" :bordered="false">
    <div v-if="!enabled">
      <!-- 未启用状态 -->
      <n-alert type="info" :show-icon="true" style="margin-bottom: 16px">
        启用二次验证可以大大提高您账户的安全性
      </n-alert>
      
      <n-steps :current="currentStep" :status="stepStatus">
        <n-step title="生成密钥" description="获取您的专属密钥" />
        <n-step title="扫描二维码" description="使用身份验证应用扫描" />
        <n-step title="验证绑定" description="输入验证码完成绑定" />
      </n-steps>
      
      <div style="margin-top: 24px;">
        <!-- 步骤1：生成密钥 -->
        <div v-if="currentStep === 1">
          <n-button type="primary" @click="generateSecret" :loading="loading">
            生成2FA密钥
          </n-button>
        </div>
        
        <!-- 步骤2：显示二维码 -->
        <div v-if="currentStep === 2" class="qr-code-section">
          <div class="qr-code-container">
            <img :src="qrCodeData" alt="2FA二维码" class="qr-code" />
          </div>
          <n-divider>或手动输入密钥</n-divider>
          <n-input-group>
            <n-input :value="secretKey" readonly />
            <n-button @click="copySecret">复制</n-button>
          </n-input-group>
          <n-button type="primary" @click="nextStep" style="margin-top: 16px;">
            我已扫描二维码
          </n-button>
        </div>
        
        <!-- 步骤3：验证绑定 -->
        <div v-if="currentStep === 3">
          <n-form-item label="验证码">
            <n-input
              v-model:value="verificationCode"
              placeholder="请输入6位验证码"
              maxlength="6"
              :input-props="{ inputmode: 'numeric' }"
            />
          </n-form-item>
          <n-space>
            <n-button @click="prevStep">上一步</n-button>
            <n-button type="primary" @click="verifyAndEnable" :loading="verifying">
              验证并启用
            </n-button>
          </n-space>
        </div>
      </div>
    </div>
    
    <!-- 已启用状态 -->
    <div v-else>
      <n-alert type="success" :show-icon="true" style="margin-bottom: 16px">
        您的账户已启用二次验证保护
      </n-alert>
      
      <n-descriptions :column="1" bordered>
        <n-descriptions-item label="启用时间">
          {{ enabledAt }}
        </n-descriptions-item>
        <n-descriptions-item label="备用恢复码">
          <n-button text @click="showBackupCodes">查看备用码</n-button>
        </n-descriptions-item>
      </n-descriptions>
      
      <n-button type="error" @click="disable2FA" style="margin-top: 16px;">
        禁用二次验证
      </n-button>
    </div>
    
    <!-- 备用码弹窗 -->
    <n-modal v-model:show="showBackupModal" preset="card" title="备用恢复码">
      <n-alert type="warning" :show-icon="true" style="margin-bottom: 16px;">
        请将这些备用码保存在安全的地方，每个备用码只能使用一次
      </n-alert>
      <n-list>
        <n-list-item v-for="(code, index) in backupCodes" :key="index">
          <n-thing>
            <template #header>
              <n-text code>{{ code }}</n-text>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
      <template #footer>
        <n-button @click="downloadBackupCodes">下载备用码</n-button>
      </template>
    </n-modal>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { generate2FASecret, verifyAndEnable2FA, get2FAStatus } from '@/api/admin/2fa'

const message = useMessage()

// 响应式数据
const enabled = ref(false)
const currentStep = ref(1)
const stepStatus = ref('process')
const loading = ref(false)
const verifying = ref(false)
const qrCodeData = ref('')
const secretKey = ref('')
const verificationCode = ref('')
const enabledAt = ref('')
const showBackupModal = ref(false)
const backupCodes = ref<string[]>([])

// 生成2FA密钥
const generateSecret = async () => {
  loading.value = true
  try {
    const result = await generate2FASecret()
    qrCodeData.value = result.qrCode
    secretKey.value = result.secret
    currentStep.value = 2
  } catch (error) {
    message.error('生成密钥失败')
  } finally {
    loading.value = false
  }
}

// 验证并启用2FA
const verifyAndEnable = async () => {
  if (!verificationCode.value || verificationCode.value.length !== 6) {
    message.error('请输入6位验证码')
    return
  }
  
  verifying.value = true
  try {
    const result = await verifyAndEnable2FA({
      secret: secretKey.value,
      code: verificationCode.value
    })
    
    if (result.success) {
      backupCodes.value = result.backupCodes
      enabled.value = true
      stepStatus.value = 'finish'
      message.success('2FA启用成功')
      showBackupModal.value = true
    }
  } catch (error) {
    message.error('验证失败，请检查验证码')
  } finally {
    verifying.value = false
  }
}

// 复制密钥
const copySecret = () => {
  navigator.clipboard.writeText(secretKey.value)
  message.success('密钥已复制到剪贴板')
}

// 下载备用码
const downloadBackupCodes = () => {
  const content = backupCodes.value.join('\n')
  const blob = new Blob([content], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'hotgo-2fa-backup-codes.txt'
  a.click()
  URL.revokeObjectURL(url)
}

// 步骤控制
const nextStep = () => {
  currentStep.value++
}

const prevStep = () => {
  currentStep.value--
}

// 初始化
onMounted(async () => {
  try {
    const status = await get2FAStatus()
    enabled.value = status.enabled
    enabledAt.value = status.enabledAt
  } catch (error) {
    console.error('获取2FA状态失败', error)
  }
})
</script>

<style scoped>
.qr-code-section {
  text-align: center;
}

.qr-code-container {
  display: inline-block;
  padding: 16px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.qr-code {
  width: 200px;
  height: 200px;
}

@media (max-width: 768px) {
  .qr-code {
    width: 150px;
    height: 150px;
  }
}
</style>
```

### 3.2 2FA验证组件

```vue
<!-- src/components/2FA/Verify2FA.vue -->
<template>
  <div class="verify-2fa-container">
    <n-card title="二次验证" :bordered="false" class="verify-card">
      <div class="verify-content">
        <n-icon size="48" color="#1890ff" class="security-icon">
          <ShieldCheckmarkOutline />
        </n-icon>
        
        <h3>请输入验证码</h3>
        <p>打开您的身份验证应用，输入6位验证码</p>
        
        <n-form @submit.prevent="verify">
          <n-form-item>
            <n-input
              v-model:value="verificationCode"
              placeholder="请输入6位验证码"
              size="large"
              maxlength="6"
              :input-props="{ inputmode: 'numeric', autocomplete: 'one-time-code' }"
              class="code-input"
            >
              <template #suffix>
                <n-countdown
                  v-if="showCountdown"
                  :duration="30 * 1000"
                  :active="true"
                  @finish="onCountdownFinish"
                />
              </template>
            </n-input>
          </n-form-item>
          
          <n-form-item>
            <n-button
              type="primary"
              size="large"
              block
              :loading="verifying"
              @click="verify"
            >
              验证
            </n-button>
          </n-form-item>
        </n-form>
        
        <n-divider>遇到问题？</n-divider>
        
        <n-collapse>
          <n-collapse-item title="使用备用恢复码" name="backup">
            <n-input
              v-model:value="backupCode"
              placeholder="请输入备用恢复码"
              size="large"
            />
            <n-button
              type="default"
              size="large"
              block
              :loading="verifying"
              @click="verifyBackupCode"
              style="margin-top: 8px;"
            >
              使用备用码验证
            </n-button>
          </n-collapse-item>
        </n-collapse>
        
        <div class="help-links">
          <n-button text @click="showHelp">需要帮助？</n-button>
          <n-button text @click="contactAdmin">联系管理员</n-button>
        </div>
      </div>
    </n-card>
    
    <!-- 帮助弹窗 -->
    <n-modal v-model:show="showHelpModal" preset="card" title="2FA验证帮助">
      <div class="help-content">
        <h4>如何获取验证码？</h4>
        <ol>
          <li>打开您的身份验证应用（如Google Authenticator）</li>
          <li>找到HotGo账户对应的条目</li>
          <li>输入显示的6位数字验证码</li>
        </ol>
        
        <h4>验证码无效？</h4>
        <ul>
          <li>确保设备时间准确</li>
          <li>验证码每30秒更新一次</li>
          <li>如果仍有问题，请使用备用恢复码</li>
        </ul>
        
        <h4>丢失了身份验证设备？</h4>
        <p>请使用备用恢复码登录，然后重新配置2FA设备</p>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { ShieldCheckmarkOutline } from '@vicons/ionicons5'
import { verify2FA } from '@/api/admin/2fa'

const router = useRouter()
const message = useMessage()

// 响应式数据
const verificationCode = ref('')
const backupCode = ref('')
const verifying = ref(false)
const showCountdown = ref(true)
const showHelpModal = ref(false)

// 验证2FA码
const verify = async () => {
  if (!verificationCode.value || verificationCode.value.length !== 6) {
    message.error('请输入6位验证码')
    return
  }
  
  verifying.value = true
  try {
    const result = await verify2FA({
      code: verificationCode.value
    })
    
    if (result.success) {
      message.success('验证成功')
      router.push('/dashboard')
    }
  } catch (error: any) {
    message.error(error.message || '验证失败')
    verificationCode.value = ''
  } finally {
    verifying.value = false
  }
}

// 验证备用码
const verifyBackupCode = async () => {
  if (!backupCode.value) {
    message.error('请输入备用恢复码')
    return
  }
  
  verifying.value = true
  try {
    const result = await verify2FA({
      code: backupCode.value
    })
    
    if (result.success) {
      message.success('验证成功')
      message.warning('建议您重新生成备用恢复码')
      router.push('/dashboard')
    }
  } catch (error: any) {
    message.error(error.message || '备用码无效')
    backupCode.value = ''
  } finally {
    verifying.value = false
  }
}

// 显示帮助
const showHelp = () => {
  showHelpModal.value = true
}

// 联系管理员
const contactAdmin = () => {
  message.info('请联系系统管理员获取帮助')
}

// 倒计时结束
const onCountdownFinish = () => {
  showCountdown.value = false
  setTimeout(() => {
    showCountdown.value = true
  }, 1000)
}
</script>

<style scoped>
.verify-2fa-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.verify-card {
  width: 100%;
  max-width: 400px;
  margin: 20px;
}

.verify-content {
  text-align: center;
  padding: 20px;
}

.security-icon {
  margin-bottom: 16px;
}

.code-input {
  font-size: 18px;
  text-align: center;
  letter-spacing: 4px;
}

.help-links {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
}

.help-content h4 {
  color: #1890ff;
  margin-top: 16px;
  margin-bottom: 8px;
}

.help-content ol,
.help-content ul {
  text-align: left;
  padding-left: 20px;
}

@media (max-width: 768px) {
  .verify-card {
    margin: 10px;
  }
  
  .verify-content {
    padding: 16px;
  }
}
</style>
```

## 4. 安全考虑

### 4.1 密钥安全
- TOTP密钥使用AES-256加密存储
- 密钥生成使用加密安全的随机数生成器
- 备用恢复码使用bcrypt哈希存储
- 定期轮换加密密钥

### 4.2 防暴力破解
- 限制验证码尝试次数（默认5次）
- 失败后账户临时锁定（默认30分钟）
- 使用Redis记录失败次数和锁定状态
- IP级别的访问频率限制

### 4.3 会话安全
- 2FA验证成功后设置会话标记
- 敏感操作需要重新验证2FA
- 会话超时后需要重新进行2FA验证
- 支持强制注销所有会话

### 4.4 审计日志
- 记录所有2FA相关操作
- 包含IP地址、用户代理、时间戳
- 异常行为自动告警
- 日志定期归档和备份

## 5. 测试方案

### 5.1 单元测试

```go
// internal/library/totp/totp_test.go
package totp

import (
    "testing"
    "time"
)

func TestGenerateSecret(t *testing.T) {
    config := &TOTPConfig{
        SecretSize: 20,
    }
    
    secret, err := GenerateSecret(config)
    if err != nil {
        t.Fatalf("生成密钥失败: %v", err)
    }
    
    if len(secret) == 0 {
        t.Fatal("生成的密钥为空")
    }
    
    // 验证Base32编码
    _, err = base32.StdEncoding.DecodeString(secret)
    if err != nil {
        t.Fatalf("密钥不是有效的Base32编码: %v", err)
    }
}

func TestValidateCode(t *testing.T) {
    config := &TOTPConfig{
        Period: 30,
        Digits: 6,
    }
    
    secret := "JBSWY3DPEHPK3PXP"
    
    // 生成当前时间的验证码
    now := time.Now().Unix()
    timeStep := now / 30
    code := generateCode(base32.StdEncoding.DecodeString(secret), timeStep, 6)
    
    // 验证生成的验证码
    if !ValidateCode(secret, code, config) {
        t.Fatal("验证码验证失败")
    }
    
    // 验证错误的验证码
    if ValidateCode(secret, "000000", config) {
        t.Fatal("错误的验证码通过了验证")
    }
}
```

### 5.2 集成测试

```go
// test/integration/2fa_test.go
package integration

import (
    "testing"
    "hotgo/internal/service"
)

func Test2FAWorkflow(t *testing.T) {
    // 1. 生成密钥
    result, err := service.Admin2FA().GenerateSecret(ctx, testMemberId)
    if err != nil {
        t.Fatalf("生成密钥失败: %v", err)
    }
    
    // 2. 模拟验证码验证
    // 这里需要使用固定的时间来生成可预测的验证码
    
    // 3. 启用2FA
    enableResult, err := service.Admin2FA().VerifyAndEnable(ctx, &adminin.VerifyAndEnable2FAInp{
        MemberId: testMemberId,
        Secret:   result.Secret,
        Code:     expectedCode,
    })
    if err != nil {
        t.Fatalf("启用2FA失败: %v", err)
    }
    
    // 4. 验证备用码
    for _, backupCode := range enableResult.BackupCodes {
        verifyResult, err := service.Admin2FA().Verify2FA(ctx, &adminin.Verify2FAInp{
            MemberId: testMemberId,
            Code:     backupCode,
        })
        if err != nil {
            t.Fatalf("备用码验证失败: %v", err)
        }
        if !verifyResult.Success {
            t.Fatal("备用码验证未成功")
        }
        break // 只测试一个备用码
    }
}
```

### 5.3 端到端测试

```javascript
// e2e/2fa.spec.js
const { test, expect } = require('@playwright/test');

test('2FA完整流程测试', async ({ page }) => {
  // 1. 登录管理员账户
  await page.goto('/login');
  await page.fill('[data-testid="username"]', 'admin');
  await page.fill('[data-testid="password"]', 'password');
  await page.click('[data-testid="login-button"]');
  
  // 2. 进入2FA设置页面
  await page.goto('/account/security');
  await page.click('[data-testid="enable-2fa-button"]');
  
  // 3. 生成密钥
  await page.click('[data-testid="generate-secret-button"]');
  await expect(page.locator('[data-testid="qr-code"]')).toBeVisible();
  
  // 4. 获取密钥并生成验证码（需要使用TOTP库）
  const secret = await page.textContent('[data-testid="secret-key"]');
  const code = generateTOTPCode(secret); // 需要实现这个函数
  
  // 5. 输入验证码并启用
  await page.click('[data-testid="next-step-button"]');
  await page.fill('[data-testid="verification-code"]', code);
  await page.click('[data-testid="verify-enable-button"]');
  
  // 6. 验证启用成功
  await expect(page.locator('[data-testid="2fa-enabled-status"]')).toBeVisible();
  
  // 7. 测试登录时的2FA验证
  await page.click('[data-testid="logout-button"]');
  await page.fill('[data-testid="username"]', 'admin');
  await page.fill('[data-testid="password"]', 'password');
  await page.click('[data-testid="login-button"]');
  
  // 8. 应该跳转到2FA验证页面
  await expect(page.locator('[data-testid="2fa-verify-page"]')).toBeVisible();
  
  // 9. 输入验证码完成登录
  const newCode = generateTOTPCode(secret);
  await page.fill('[data-testid="2fa-code-input"]', newCode);
  await page.click('[data-testid="2fa-verify-button"]');
  
  // 10. 验证登录成功
  await expect(page.locator('[data-testid="dashboard"]')).toBeVisible();
});

test('2FA备用码测试', async ({ page }) => {
  // 测试备用码的使用流程
  // ...
});

test('2FA错误处理测试', async ({ page }) => {
  // 测试各种错误情况的处理
  // ...
});
```

## 6. 部署配置

### 6.1 环境变量配置

```yaml
# config/config.yaml
system:
  # 2FA相关配置
  security:
    2fa:
      enabled: true
      issuer: "HotGo"
      secretSize: 20
      period: 30
      digits: 6
      maxAttempts: 5
      lockoutDuration: 1800 # 30分钟
      
  # 加密配置
  encryption:
    key: "your-32-byte-encryption-key-here"
    algorithm: "AES-256-GCM"
```

### 6.2 数据库迁移脚本

```sql
-- migration/20240101_add_2fa_tables.sql
-- 执行数据库迁移
SOURCE /path/to/2fa_tables.sql;

-- 插入默认配置
INSERT INTO `hg_sys_config` (`group`, `key`, `value`, `type`, `title`, `tip`, `rule`, `extend`, `remark`, `sort`, `status`, `created_at`, `updated_at`) VALUES
('security', '2faEnabled', '1', 'switch', '启用二次验证', '是否启用系统二次验证功能', '', '', '控制整个系统的2FA功能开关', 100, 1, NOW(), NOW());
```

### 6.3 Nginx配置

```nginx
# 2FA相关的安全头
location /api/admin/2fa/ {
    # 限制请求频率
    limit_req zone=api burst=10 nodelay;
    
    # 安全头
    add_header X-Content-Type-Options nosniff;
    add_header X-Frame-Options DENY;
    add_header X-XSS-Protection "1; mode=block";
    
    # 代理到后端
    proxy_pass http://backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

### 6.4 监控和告警

```yaml
# monitoring/2fa-alerts.yml
groups:
  - name: 2fa-security
    rules:
      - alert: High2FAFailureRate
        expr: rate(hotgo_2fa_verify_failures_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "2FA验证失败率过高"
          description: "过去5分钟内2FA验证失败率超过10%"
          
      - alert: 2FABruteForceAttack
        expr: rate(hotgo_2fa_lockout_total[1m]) > 0.05
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "检测到2FA暴力破解攻击"
          description: "过去1分钟内账户锁定率异常"
```

## 7. 用户文档

### 7.1 管理员使用指南

#### 如何启用2FA
1. 登录管理后台
2. 进入"个人设置" -> "安全设置"
3. 点击"启用二次验证"
4. 使用身份验证应用扫描二维码
5. 输入验证码完成绑定
6. 保存备用恢复码

#### 推荐的身份验证应用
- Google Authenticator
- Microsoft Authenticator
- Authy
- 1Password

#### 备用恢复码使用
- 每个备用码只能使用一次
- 请将备用码保存在安全的地方
- 建议定期重新生成备用码

### 7.2 系统管理员配置指南

#### 全局2FA配置
1. 进入"系统管理" -> "安全配置"
2. 配置2FA全局开关
3. 设置强制启用策略
4. 配置安全参数（尝试次数、锁定时间等）

#### 安全策略建议
- 建议为所有管理员启用2FA
- 定期检查安全日志
- 设置合理的锁定策略
- 建立应急处理流程

## 8. 故障排除

### 8.1 常见问题

**Q: 验证码总是提示错误？**
A: 检查设备时间是否准确，TOTP基于时间同步

**Q: 丢失了身份验证设备怎么办？**
A: 使用备用恢复码登录，然后重新配置2FA

**Q: 备用码也丢失了怎么办？**
A: 联系系统管理员重置2FA配置

**Q: 2FA功能无法启用？**
A: 检查系统配置中的2FA开关是否启用

### 8.2 错误代码说明

| 错误代码 | 说明 | 解决方案 |
|----------|------|----------|
| 2FA_001 | 验证码错误 | 检查验证码是否正确，注意时间同步 |
| 2FA_002 | 验证码已过期 | 使用最新的验证码 |
| 2FA_003 | 尝试次数过多 | 等待锁定时间结束后重试 |
| 2FA_004 | 2FA未启用 | 先启用2FA功能 |
| 2FA_005 | 备用码无效 | 检查备用码是否正确且未使用过 |

这份实施指南提供了完整的2FA功能实现方案，包括后端服务、前端组件、安全考虑、测试方案和部署配置等各个方面，确保功能的安全性、可用性和可维护性。