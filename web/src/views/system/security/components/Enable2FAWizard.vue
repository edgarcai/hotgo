<template>
  <div class="enable-2fa-wizard">
    <n-steps :current="currentStep" :status="stepStatus">
      <n-step title="验证密码" description="确认您的身份" />
      <n-step title="扫描二维码" description="使用验证器应用扫描" />
      <n-step title="验证设置" description="输入验证码完成设置" />
      <n-step title="保存备用码" description="保存恢复码以备不时之需" />
    </n-steps>
    
    <div class="step-content">
      <!-- 步骤1: 验证密码 -->
      <div v-if="currentStep === 1" class="step-panel">
        <div class="step-title">验证您的密码</div>
        <div class="step-desc">为了安全起见，请输入您的当前密码以继续设置双因素认证。</div>
        
        <n-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules">
          <n-form-item path="password" label="当前密码">
            <n-input
              v-model:value="passwordForm.password"
              type="password"
              placeholder="请输入当前密码"
              show-password-on="click"
              @keyup.enter="handlePasswordVerify"
            />
          </n-form-item>
        </n-form>
        
        <div class="step-actions">
          <n-space>
            <n-button @click="$emit('cancel')">取消</n-button>
            <n-button type="primary" @click="handlePasswordVerify" :loading="verifyingPassword">
              验证密码
            </n-button>
          </n-space>
        </div>
      </div>
      
      <!-- 步骤2: 扫描二维码 -->
      <div v-if="currentStep === 2" class="step-panel">
        <div class="step-title">扫描二维码</div>
        <div class="step-desc">使用您的验证器应用扫描下方二维码，然后输入应用显示的6位验证码。</div>
        
        <div class="qr-section">
          <div class="qr-code">
            <img v-if="qrCodeUrl" :src="qrCodeUrl" alt="2FA QR Code" />
            <n-spin v-else size="large" />
          </div>
          
          <div class="qr-info">
            <div class="secret-key">
              <div class="label">手动输入密钥（如果无法扫描）：</div>
              <div class="key-value">
                <n-input readonly :value="secretKey" />
                <n-button text @click="copySecretKey">
                  <template #icon>
                    <n-icon><CopyOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
            
            <div class="app-recommendations">
              <div class="label">推荐的验证器应用：</div>
              <n-space>
                <n-tag>Google Authenticator</n-tag>
                <n-tag>Microsoft Authenticator</n-tag>
                <n-tag>Authy</n-tag>
              </n-space>
            </div>
          </div>
        </div>
        
        <div class="step-actions">
          <n-space>
            <n-button @click="currentStep = 1">上一步</n-button>
            <n-button type="primary" @click="currentStep = 3">
              我已扫描二维码
            </n-button>
          </n-space>
        </div>
      </div>
      
      <!-- 步骤3: 验证设置 -->
      <div v-if="currentStep === 3" class="step-panel">
        <div class="step-title">验证设置</div>
        <div class="step-desc">请输入验证器应用显示的6位验证码以完成设置。</div>
        
        <div class="verify-section">
          <TwoFactorInput
            v-model:code="verificationCode"
            :mode="'totp'"
            :error="verificationError"
            @submit="handleVerifySetup"
            @input="verificationError = ''"
          />
        </div>
        
        <div class="step-actions">
          <n-space>
            <n-button @click="currentStep = 2">上一步</n-button>
            <n-button type="primary" @click="handleVerifySetup" :loading="verifyingSetup">
              验证并启用
            </n-button>
          </n-space>
        </div>
      </div>
      
      <!-- 步骤4: 保存备用码 -->
      <div v-if="currentStep === 4" class="step-panel">
        <div class="step-title">保存备用恢复码</div>
        <div class="step-desc">请将这些备用恢复码保存在安全的地方。当您无法使用验证器应用时，可以使用这些代码登录。</div>
        
        <div class="backup-codes-section">
          <n-alert type="warning" title="重要提醒" style="margin-bottom: 16px;">
            <ul>
              <li>每个备用码只能使用一次</li>
              <li>请将这些代码保存在安全的地方</li>
              <li>不要与他人分享这些代码</li>
              <li>如果丢失，您可以稍后重新生成</li>
            </ul>
          </n-alert>
          
          <div class="codes-grid">
            <div v-for="(code, index) in backupCodes" :key="index" class="backup-code">
              {{ code }}
            </div>
          </div>
          
          <div class="codes-actions">
            <n-space>
              <n-button @click="downloadBackupCodes">
                <template #icon>
                  <n-icon><DownloadOutline /></n-icon>
                </template>
                下载备用码
              </n-button>
              
              <n-button @click="copyBackupCodes">
                <template #icon>
                  <n-icon><CopyOutline /></n-icon>
                </template>
                复制备用码
              </n-button>
            </n-space>
          </div>
        </div>
        
        <div class="step-actions">
          <n-checkbox v-model:checked="confirmSaved">
            我已安全保存这些备用恢复码
          </n-checkbox>
          
          <n-space style="margin-top: 16px;">
            <n-button type="primary" @click="handleComplete" :disabled="!confirmSaved">
              完成设置
            </n-button>
          </n-space>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import { CopyOutline, DownloadOutline } from '@vicons/ionicons5'
import { enable2FA, verifySetup2FA, get2FAQRCode } from '@/api/auth/twoFactor'
import TwoFactorInput from '@/components/TwoFactor/TwoFactorInput.vue'

// 事件定义
const emit = defineEmits<{
  success: []
  cancel: []
}>()

// 响应式数据
const message = useMessage()
const currentStep = ref(1)
const verifyingPassword = ref(false)
const verifyingSetup = ref(false)
const qrCodeUrl = ref('')
const secretKey = ref('')
const verificationCode = ref('')
const verificationError = ref('')
const backupCodes = ref<string[]>([])
const confirmSaved = ref(false)

// 表单数据
const passwordForm = ref({
  password: ''
})

// 表单验证规则
const passwordRules = {
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur'
  }
}

// 计算属性
const stepStatus = computed(() => {
  if (currentStep.value === 4) return 'finish'
  return 'process'
})

// 方法
const handlePasswordVerify = async () => {
  if (!passwordForm.value.password) {
    message.error('请输入密码')
    return
  }
  
  verifyingPassword.value = true
  try {
    // 验证密码并获取二维码
    const response = await get2FAQRCode()
    if (response.code === 200) {
      qrCodeUrl.value = response.data.qrCodeUrl
      secretKey.value = response.data.secretKey
      currentStep.value = 2
    } else {
      message.error(response.message || '密码验证失败')
    }
  } catch (error: any) {
    message.error(error.message || '密码验证失败')
  } finally {
    verifyingPassword.value = false
  }
}

const handleVerifySetup = async () => {
  if (!verificationCode.value || verificationCode.value.length !== 6) {
    verificationError.value = '请输入6位验证码'
    return
  }
  
  verifyingSetup.value = true
  try {
    const response = await verifySetup2FA({
      code: verificationCode.value,
      secretKey: secretKey.value
    })
    
    if (response.code === 200) {
      backupCodes.value = response.data.backupCodes
      currentStep.value = 4
    } else {
      verificationError.value = response.message || '验证码错误'
    }
  } catch (error: any) {
    verificationError.value = error.message || '验证失败'
  } finally {
    verifyingSetup.value = false
  }
}

const copySecretKey = async () => {
  try {
    await navigator.clipboard.writeText(secretKey.value)
    message.success('密钥已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

const copyBackupCodes = async () => {
  try {
    const codesText = backupCodes.value.join('\n')
    await navigator.clipboard.writeText(codesText)
    message.success('备用码已复制到剪贴板')
  } catch (error) {
    message.error('复制失败')
  }
}

const downloadBackupCodes = () => {
  const codesText = backupCodes.value.join('\n')
  const blob = new Blob([codesText], { type: 'text/plain' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = '2FA-backup-codes.txt'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  message.success('备用码已下载')
}

const handleComplete = () => {
  emit('success')
}
</script>

<style lang="less" scoped>
.enable-2fa-wizard {
  .step-content {
    margin-top: 32px;
    
    .step-panel {
      .step-title {
        font-size: 18px;
        font-weight: 500;
        color: #333;
        margin-bottom: 8px;
      }
      
      .step-desc {
        font-size: 14px;
        color: #666;
        line-height: 1.5;
        margin-bottom: 24px;
      }
      
      .step-actions {
        margin-top: 32px;
        padding-top: 24px;
        border-top: 1px solid #f0f0f0;
      }
    }
  }
  
  .qr-section {
    .qr-code {
      text-align: center;
      margin-bottom: 24px;
      
      img {
        width: 200px;
        height: 200px;
        border: 1px solid #e0e0e0;
        border-radius: 8px;
      }
    }
    
    .qr-info {
      .secret-key {
        margin-bottom: 16px;
        
        .label {
          font-size: 14px;
          color: #666;
          margin-bottom: 8px;
        }
        
        .key-value {
          display: flex;
          gap: 8px;
          align-items: center;
        }
      }
      
      .app-recommendations {
        .label {
          font-size: 14px;
          color: #666;
          margin-bottom: 8px;
        }
      }
    }
  }
  
  .verify-section {
    display: flex;
    justify-content: center;
    margin: 32px 0;
  }
  
  .backup-codes-section {
    .codes-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;
      margin: 16px 0;
      
      .backup-code {
        padding: 12px;
        background-color: #f8f9fa;
        border: 1px solid #e9ecef;
        border-radius: 6px;
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 14px;
        text-align: center;
        color: #333;
      }
    }
    
    .codes-actions {
      margin: 16px 0;
      text-align: center;
    }
  }
}
</style>