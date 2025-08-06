<template>
  <div class="two-factor-input">
    <!-- TOTP模式 -->
    <div v-if="mode === 'totp'" class="totp-mode">
      <div class="input-section">
        <n-form-item :show-label="false">
          <n-input
            ref="totpInputRef"
            v-model:value="code"
            data-testid="totp-input"
            placeholder="请输入6位验证码"
            maxlength="6"
            size="large"
            :status="error ? 'error' : undefined"
            @input="handleTotpInput"
            @keyup.enter="handleSubmit"
            :disabled="loading"
          >
            <template #prefix>
              <n-icon size="18" color="#808695">
                <SafetyCertificateOutlined />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
      </div>
      
      <!-- 倒计时显示 -->
      <div class="countdown-section" data-testid="countdown">
        <n-progress
          type="circle"
          :percentage="countdownPercentage"
          :size="60"
          :show-indicator="false"
          :color="countdownColor"
        >
          <template #default>
            <span class="countdown-text">{{ countdown }}</span>
          </template>
        </n-progress>
        <div class="countdown-label">剩余时间</div>
      </div>
      
      <!-- 错误信息 -->
      <div v-if="error" class="error-section" data-testid="error-message">
        <n-alert type="error" :show-icon="false">
          {{ error }}
        </n-alert>
      </div>
      
      <!-- 切换到备用码 -->
      <div class="switch-section">
        <n-button text type="primary" @click="handleSwitchMode">
          使用备用恢复码
        </n-button>
      </div>
    </div>
    
    <!-- 备用码模式 -->
    <div v-else-if="mode === 'backup'" class="backup-mode">
      <div class="input-section">
        <n-form-item :show-label="false">
          <n-input
            ref="backupInputRef"
            v-model:value="code"
            data-testid="backup-input"
            placeholder="请输入备用恢复码"
            size="large"
            :status="error ? 'error' : undefined"
            @input="handleBackupInput"
            @keyup.enter="handleSubmit"
            :disabled="loading"
          >
            <template #prefix>
              <n-icon size="18" color="#808695">
                <KeyOutlined />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
      </div>
      
      <!-- 说明文字 -->
      <div class="backup-hint">
        <n-text depth="3" style="font-size: 12px;">
          备用恢复码是在设置双因素认证时生成的一次性代码
        </n-text>
      </div>
      
      <!-- 错误信息 -->
      <div v-if="error" class="error-section" data-testid="error-message">
        <n-alert type="error" :show-icon="false">
          {{ error }}
        </n-alert>
      </div>
      
      <!-- 切换到TOTP -->
      <div class="switch-section">
        <n-button text type="primary" @click="handleSwitchMode">
          使用验证器应用
        </n-button>
      </div>
    </div>
    
    <!-- 提交按钮 -->
    <div class="submit-section">
      <n-button
        type="primary"
        size="large"
        block
        :loading="loading"
        :disabled="!isValidCode"
        @click="handleSubmit"
      >
        验证
      </n-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { SafetyCertificateOutlined, KeyOutlined } from '@vicons/antd'

// 组件属性定义
interface Props {
  mode: 'totp' | 'backup' // 验证模式
  loading?: boolean       // 加载状态
  error?: string         // 错误信息
}

// 组件事件定义
interface Emits {
  (e: 'submit', code: string): void
  (e: 'switch-mode'): void
  (e: 'update:error', error: string | undefined): void
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  error: undefined
})

const emit = defineEmits<Emits>()

// 响应式数据
const code = ref('')
const countdown = ref(30)
const totpInputRef = ref()
const backupInputRef = ref()
let countdownTimer: NodeJS.Timeout | null = null

// 计算属性
const isValidCode = computed(() => {
  if (props.mode === 'totp') {
    return /^\d{6}$/.test(code.value)
  } else {
    return code.value.length >= 8 // 备用码通常较长
  }
})

const countdownPercentage = computed(() => {
  return (countdown.value / 30) * 100
})

const countdownColor = computed(() => {
  if (countdown.value > 10) return '#18a058'
  if (countdown.value > 5) return '#f0a020'
  return '#d03050'
})

// 方法
const handleTotpInput = (value: string) => {
  // 只允许数字输入
  const numericValue = value.replace(/\D/g, '')
  code.value = numericValue.slice(0, 6)
  
  // 清除错误信息
  if (props.error) {
    emit('update:error', undefined)
  }
  
  // 自动提交
  if (numericValue.length === 6) {
    nextTick(() => {
      handleSubmit()
    })
  }
}

const handleBackupInput = (value: string) => {
  // 移除空格和特殊字符，保留字母数字
  code.value = value.replace(/[^a-zA-Z0-9]/g, '').toLowerCase()
  
  // 清除错误信息
  if (props.error) {
    emit('update:error', undefined)
  }
}

const handleSubmit = () => {
  if (!isValidCode.value || props.loading) {
    return
  }
  
  emit('submit', code.value)
}

const handleSwitchMode = () => {
  code.value = ''
  emit('update:error', undefined)
  emit('switch-mode')
}

const startCountdown = () => {
  countdown.value = 30
  
  if (countdownTimer) {
    clearInterval(countdownTimer)
  }
  
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      countdown.value = 30 // 重置倒计时
    }
  }, 1000)
}

const stopCountdown = () => {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  }
}

// 监听模式变化，自动聚焦
watch(() => props.mode, async (newMode) => {
  await nextTick()
  if (newMode === 'totp' && totpInputRef.value) {
    totpInputRef.value.focus()
  } else if (newMode === 'backup' && backupInputRef.value) {
    backupInputRef.value.focus()
  }
})

// 生命周期
onMounted(() => {
  if (props.mode === 'totp') {
    startCountdown()
    nextTick(() => {
      if (totpInputRef.value) {
        totpInputRef.value.focus()
      }
    })
  } else {
    nextTick(() => {
      if (backupInputRef.value) {
        backupInputRef.value.focus()
      }
    })
  }
})

onUnmounted(() => {
  stopCountdown()
})

// 暴露方法给父组件
defineExpose({
  focus: () => {
    if (props.mode === 'totp' && totpInputRef.value) {
      totpInputRef.value.focus()
    } else if (props.mode === 'backup' && backupInputRef.value) {
      backupInputRef.value.focus()
    }
  },
  clear: () => {
    code.value = ''
  }
})
</script>

<style lang="less" scoped>
.two-factor-input {
  width: 100%;
  
  .totp-mode,
  .backup-mode {
    .input-section {
      margin-bottom: 16px;
    }
    
    .countdown-section {
      display: flex;
      flex-direction: column;
      align-items: center;
      margin-bottom: 16px;
      
      .countdown-text {
        font-size: 14px;
        font-weight: 500;
        color: #333;
      }
      
      .countdown-label {
        margin-top: 8px;
        font-size: 12px;
        color: #666;
      }
    }
    
    .backup-hint {
      margin-bottom: 16px;
      text-align: center;
      padding: 8px;
      background-color: #f8f9fa;
      border-radius: 4px;
    }
    
    .error-section {
      margin-bottom: 16px;
    }
    
    .switch-section {
      text-align: center;
      margin-bottom: 24px;
    }
  }
  
  .submit-section {
    margin-top: 16px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .two-factor-input {
    .countdown-section {
      :deep(.n-progress) {
        width: 50px !important;
        height: 50px !important;
      }
      
      .countdown-text {
        font-size: 12px;
      }
    }
  }
}
</style>