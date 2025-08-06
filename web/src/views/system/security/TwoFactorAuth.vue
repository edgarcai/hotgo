<template>
  <div class="two-factor-auth">
    <n-card title="双因素认证" :bordered="false">
      <template #header-extra>
        <n-tag :type="twoFactorStatus?.enabled ? 'success' : 'warning'" size="small">
          {{ twoFactorStatus?.enabled ? '已启用' : '未启用' }}
        </n-tag>
      </template>
      
      <!-- 状态概览 -->
      <div class="status-overview">
        <n-space :size="24">
          <div class="status-item">
            <div class="status-icon">
              <n-icon size="32" :color="twoFactorStatus?.enabled ? '#18a058' : '#f0a020'">
                <ShieldCheckmarkOutline v-if="twoFactorStatus?.enabled" />
                <ShieldOutline v-else />
              </n-icon>
            </div>
            <div class="status-content">
              <div class="status-title">安全状态</div>
              <div class="status-desc">
                {{ twoFactorStatus?.enabled ? '您的账户已启用双因素认证保护' : '建议启用双因素认证以增强账户安全性' }}
              </div>
            </div>
          </div>
          
          <div v-if="twoFactorStatus?.enabled" class="status-item">
            <div class="status-icon">
              <n-icon size="32" color="#2080f0">
                <KeyOutline />
              </n-icon>
            </div>
            <div class="status-content">
              <div class="status-title">备用恢复码</div>
              <div class="status-desc">
                剩余 {{ twoFactorStatus?.backupCodesCount || 0 }} 个备用码
              </div>
            </div>
          </div>
        </n-space>
      </div>
      
      <!-- 操作区域 -->
      <div class="action-section">
        <n-space :size="16">
          <!-- 未启用时显示启用按钮 -->
          <template v-if="!twoFactorStatus?.enabled">
            <n-button type="primary" @click="handleEnable2FA" :loading="loading">
              <template #icon>
                <n-icon><ShieldCheckmarkOutline /></n-icon>
              </template>
              启用双因素认证
            </n-button>
          </template>
          
          <!-- 已启用时显示管理选项 -->
          <template v-else>
            <n-button @click="handleRegenerateBackupCodes" :loading="regenerateLoading">
              <template #icon>
                <n-icon><RefreshOutline /></n-icon>
              </template>
              重新生成备用码
            </n-button>
            
            <n-button type="error" @click="handleDisable2FA" :loading="disableLoading">
              <template #icon>
                <n-icon><ShieldOutline /></n-icon>
              </template>
              禁用双因素认证
            </n-button>
          </template>
        </n-space>
      </div>
      
      <!-- 使用说明 -->
      <div class="help-section">
        <n-collapse>
          <n-collapse-item title="什么是双因素认证？" name="what">
            <p>双因素认证（2FA）是一种安全措施，要求您在登录时提供两种不同的身份验证方式：</p>
            <ul>
              <li>您知道的信息（密码）</li>
              <li>您拥有的设备（手机上的验证器应用）</li>
            </ul>
            <p>这大大提高了您账户的安全性，即使密码被泄露，攻击者也无法访问您的账户。</p>
          </n-collapse-item>
          
          <n-collapse-item title="如何使用验证器应用？" name="how">
            <p>推荐使用以下验证器应用：</p>
            <ul>
              <li><strong>Google Authenticator</strong> - 适用于 iOS 和 Android</li>
              <li><strong>Microsoft Authenticator</strong> - 适用于 iOS 和 Android</li>
              <li><strong>Authy</strong> - 支持多设备同步</li>
            </ul>
            <p>设置步骤：</p>
            <ol>
              <li>在手机上下载并安装验证器应用</li>
              <li>点击"启用双因素认证"按钮</li>
              <li>使用验证器应用扫描二维码</li>
              <li>输入验证器显示的6位数字完成设置</li>
            </ol>
          </n-collapse-item>
          
          <n-collapse-item title="备用恢复码的作用？" name="backup">
            <p>备用恢复码是一次性使用的代码，用于在以下情况下恢复账户访问：</p>
            <ul>
              <li>手机丢失或损坏</li>
              <li>验证器应用无法使用</li>
              <li>更换新设备</li>
            </ul>
            <p><strong>重要提醒：</strong></p>
            <ul>
              <li>请将备用码保存在安全的地方</li>
              <li>每个备用码只能使用一次</li>
              <li>建议定期重新生成备用码</li>
            </ul>
          </n-collapse-item>
        </n-collapse>
      </div>
    </n-card>
    
    <!-- 启用2FA对话框 -->
    <n-modal v-model:show="showEnableModal" preset="card" title="启用双因素认证" style="width: 500px;">
      <Enable2FAWizard @success="handleEnableSuccess" @cancel="showEnableModal = false" />
    </n-modal>
    
    <!-- 禁用2FA对话框 -->
    <n-modal v-model:show="showDisableModal" preset="dialog" title="禁用双因素认证">
      <template #default>
        <Disable2FAForm @success="handleDisableSuccess" @cancel="showDisableModal = false" />
      </template>
    </n-modal>
    
    <!-- 备用码显示对话框 -->
    <n-modal v-model:show="showBackupCodesModal" preset="card" title="备用恢复码" style="width: 500px;">
      <BackupCodesDisplay :backup-codes="backupCodes" @close="showBackupCodesModal = false" />
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import { ShieldCheckmarkOutline, ShieldOutline, KeyOutline, RefreshOutline } from '@vicons/ionicons5'
import { getTwoFactorStatus, regenerateBackupCodes } from '@/api/auth/twoFactor'
import type { TwoFactorStatusResponse } from '@/api/auth/twoFactor'
import Enable2FAWizard from './components/Enable2FAWizard.vue'
import Disable2FAForm from './components/Disable2FAForm.vue'
import BackupCodesDisplay from './components/BackupCodesDisplay.vue'

// 响应式数据
const message = useMessage()
const loading = ref(false)
const regenerateLoading = ref(false)
const disableLoading = ref(false)
const twoFactorStatus = ref<TwoFactorStatusResponse | null>(null)
const showEnableModal = ref(false)
const showDisableModal = ref(false)
const showBackupCodesModal = ref(false)
const backupCodes = ref<string[]>([])

// 方法
const loadTwoFactorStatus = async () => {
  try {
    const response = await getTwoFactorStatus()
    if (response.code === 200) {
      twoFactorStatus.value = response.data
    }
  } catch (error) {
    console.error('获取2FA状态失败:', error)
  }
}

const handleEnable2FA = () => {
  showEnableModal.value = true
}

const handleEnableSuccess = () => {
  showEnableModal.value = false
  message.success('双因素认证启用成功')
  loadTwoFactorStatus()
}

const handleDisable2FA = () => {
  showDisableModal.value = true
}

const handleDisableSuccess = () => {
  showDisableModal.value = false
  message.success('双因素认证已禁用')
  loadTwoFactorStatus()
}

const handleRegenerateBackupCodes = async () => {
  // 这里应该先要求用户输入密码确认
  const password = prompt('请输入当前密码以确认操作：')
  if (!password) {
    return
  }
  
  regenerateLoading.value = true
  try {
    const response = await regenerateBackupCodes(password)
    if (response.code === 200) {
      backupCodes.value = response.data.backupCodes
      showBackupCodesModal.value = true
      message.success('备用恢复码已重新生成')
      loadTwoFactorStatus()
    } else {
      message.error(response.message || '重新生成失败')
    }
  } catch (error: any) {
    message.error(error.message || '重新生成失败')
  } finally {
    regenerateLoading.value = false
  }
}

// 生命周期
onMounted(() => {
  loadTwoFactorStatus()
})
</script>

<style lang="less" scoped>
.two-factor-auth {
  .status-overview {
    margin-bottom: 32px;
    padding: 24px;
    background-color: #fafafa;
    border-radius: 8px;
    
    .status-item {
      display: flex;
      align-items: center;
      gap: 16px;
      
      .status-icon {
        flex-shrink: 0;
      }
      
      .status-content {
        .status-title {
          font-size: 16px;
          font-weight: 500;
          color: #333;
          margin-bottom: 4px;
        }
        
        .status-desc {
          font-size: 14px;
          color: #666;
          line-height: 1.4;
        }
      }
    }
  }
  
  .action-section {
    margin-bottom: 32px;
    padding-bottom: 24px;
    border-bottom: 1px solid #f0f0f0;
  }
  
  .help-section {
    :deep(.n-collapse-item__content-wrapper) {
      padding: 16px 0;
    }
    
    p {
      margin: 0 0 12px 0;
      line-height: 1.6;
      color: #333;
    }
    
    ul, ol {
      margin: 0 0 12px 0;
      padding-left: 20px;
      
      li {
        margin-bottom: 6px;
        line-height: 1.5;
        color: #666;
        
        strong {
          color: #333;
        }
      }
    }
  }
}
</style>