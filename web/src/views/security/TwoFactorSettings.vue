<template>
  <div class="two-factor-settings">
    <!-- 页面头部 -->
    <div class="page-header">
      <n-page-header @back="handleBack">
        <template #title>
          <n-space align="center" :size="12">
            <n-icon size="24" color="#18a058">
              <ShieldCheckmarkOutline />
            </n-icon>
            <span>双因素认证设置</span>
          </n-space>
        </template>
        <template #subtitle> 增强您的账户安全性 </template>
        <template #extra>
          <n-tag :type="twoFactorStatus?.enabled ? 'success' : 'warning'" size="medium">
            <template #icon>
              <n-icon>
                <ShieldCheckmarkOutline v-if="twoFactorStatus?.enabled" />
                <ShieldOutline v-else />
              </n-icon>
            </template>
            {{ twoFactorStatus?.enabled ? '已启用' : '未启用' }}
          </n-tag>
        </template>
      </n-page-header>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <n-grid :cols="24" :x-gap="24" :y-gap="24">
        <!-- 左侧：状态概览和操作 -->
        <n-grid-item :span="isMobile ? 24 : 16">
          <n-card title="安全状态" :bordered="false" class="status-card">
            <!-- 状态概览 -->
            <div class="status-overview">
              <n-space :size="24" :vertical="isMobile">
                <!-- 主要状态 -->
                <div class="status-item primary">
                  <div class="status-icon">
                    <n-icon size="48" :color="twoFactorStatus?.enabled ? '#18a058' : '#f0a020'">
                      <ShieldCheckmarkOutline v-if="twoFactorStatus?.enabled" />
                      <ShieldOutline v-else />
                    </n-icon>
                  </div>
                  <div class="status-content">
                    <div class="status-title">
                      {{ twoFactorStatus?.enabled ? '双因素认证已启用' : '双因素认证未启用' }}
                    </div>
                    <div class="status-desc">
                      {{
                        twoFactorStatus?.enabled
                          ? '您的账户受到双因素认证保护，安全性良好'
                          : '建议启用双因素认证以增强账户安全性'
                      }}
                    </div>
                    <div
                      v-if="twoFactorStatus?.enabled && twoFactorStatus.enabledAt"
                      class="status-meta"
                    >
                      启用时间：{{ formatDate(twoFactorStatus.enabledAt) }}
                    </div>
                  </div>
                </div>

                <!-- 备用码状态 -->
                <div v-if="twoFactorStatus?.enabled" class="status-item secondary">
                  <div class="status-icon">
                    <n-icon size="32" :color="getBackupCodesColor()">
                      <KeyOutline />
                    </n-icon>
                  </div>
                  <div class="status-content">
                    <div class="status-title">备用恢复码</div>
                    <div class="status-desc">
                      剩余 {{ twoFactorStatus?.backupCodesCount || 0 }} 个备用码
                    </div>
                    <div v-if="twoFactorStatus?.lastUsedAt" class="status-meta">
                      最后使用：{{ formatDate(twoFactorStatus.lastUsedAt) }}
                    </div>
                  </div>
                </div>
              </n-space>
            </div>

            <!-- 操作按钮区域 -->
            <div class="action-section">
              <n-space :size="16" :vertical="isMobile">
                <!-- 未启用时的操作 -->
                <template v-if="!twoFactorStatus?.enabled">
                  <n-button
                    type="primary"
                    size="large"
                    @click="handleEnable2FA"
                    :loading="loading"
                    class="primary-action"
                  >
                    <template #icon>
                      <n-icon><ShieldCheckmarkOutline /></n-icon>
                    </template>
                    启用双因素认证
                  </n-button>
                </template>

                <!-- 已启用时的操作 -->
                <template v-else>
                  <n-button
                    @click="handleShowBackupCodes"
                    :loading="showingBackupCodes"
                    size="large"
                  >
                    <template #icon>
                      <n-icon><EyeOutline /></n-icon>
                    </template>
                    查看备用码
                  </n-button>

                  <n-button
                    @click="handleRegenerateBackupCodes"
                    :loading="regenerateLoading"
                    size="large"
                  >
                    <template #icon>
                      <n-icon><RefreshOutline /></n-icon>
                    </template>
                    重新生成备用码
                  </n-button>

                  <n-button
                    type="error"
                    @click="handleDisable2FA"
                    :loading="disableLoading"
                    size="large"
                  >
                    <template #icon>
                      <n-icon><ShieldOutline /></n-icon>
                    </template>
                    禁用双因素认证
                  </n-button>
                </template>
              </n-space>
            </div>
          </n-card>
        </n-grid-item>

        <!-- 右侧：帮助信息 -->
        <n-grid-item :span="isMobile ? 24 : 8">
          <n-card title="使用指南" :bordered="false" class="help-card">
            <n-collapse>
              <n-collapse-item title="什么是双因素认证？" name="what">
                <div class="help-content">
                  <p>双因素认证（2FA）是一种安全措施，要求您在登录时提供两种不同的身份验证方式：</p>
                  <ul>
                    <li>您知道的信息（密码）</li>
                    <li>您拥有的设备（手机上的验证器应用）</li>
                  </ul>
                  <p>这大大提高了您账户的安全性，即使密码被泄露，攻击者也无法访问您的账户。</p>
                </div>
              </n-collapse-item>

              <n-collapse-item title="推荐的验证器应用" name="apps">
                <div class="help-content">
                  <div class="app-list">
                    <div class="app-item">
                      <n-icon size="20" color="#4285f4"><PhonePortraitOutline /></n-icon>
                      <div>
                        <div class="app-name">Google Authenticator</div>
                        <div class="app-desc">适用于 iOS 和 Android</div>
                      </div>
                    </div>
                    <div class="app-item">
                      <n-icon size="20" color="#00a1f1"><PhonePortraitOutline /></n-icon>
                      <div>
                        <div class="app-name">Microsoft Authenticator</div>
                        <div class="app-desc">支持云同步</div>
                      </div>
                    </div>
                    <div class="app-item">
                      <n-icon size="20" color="#ec1c24"><PhonePortraitOutline /></n-icon>
                      <div>
                        <div class="app-name">Authy</div>
                        <div class="app-desc">支持多设备同步</div>
                      </div>
                    </div>
                  </div>
                </div>
              </n-collapse-item>

              <n-collapse-item title="设置步骤" name="steps">
                <div class="help-content">
                  <ol>
                    <li>在手机上下载并安装验证器应用</li>
                    <li>点击"启用双因素认证"按钮</li>
                    <li>输入当前密码进行身份验证</li>
                    <li>使用验证器应用扫描二维码</li>
                    <li>输入验证器显示的6位数字完成设置</li>
                    <li>保存备用恢复码到安全位置</li>
                  </ol>
                </div>
              </n-collapse-item>

              <n-collapse-item title="安全提示" name="security">
                <div class="help-content">
                  <n-alert type="warning" :show-icon="false">
                    <ul class="security-tips">
                      <li>请将备用恢复码保存在安全的地方</li>
                      <li>不要将恢复码存储在与账户相同的设备上</li>
                      <li>定期检查并更新您的验证器应用</li>
                      <li>如果更换手机，请及时重新设置2FA</li>
                    </ul>
                  </n-alert>
                </div>
              </n-collapse-item>
            </n-collapse>
          </n-card>
        </n-grid-item>
      </n-grid>
    </div>

    <!-- 启用2FA向导弹窗 -->
    <n-modal
      v-model:show="showEnableWizard"
      :mask-closable="false"
      :close-on-esc="false"
      preset="card"
      title="启用双因素认证"
      class="enable-wizard-modal"
      :style="{ width: isMobile ? '95%' : '600px' }"
    >
      <Enable2FAWizard @success="handleEnableSuccess" @cancel="showEnableWizard = false" />
    </n-modal>

    <!-- 禁用2FA确认弹窗 -->
    <n-modal
      v-model:show="showDisableModal"
      preset="card"
      title="禁用双因素认证"
      :style="{ width: isMobile ? '95%' : '500px' }"
    >
      <Disable2FAForm @success="handleDisableSuccess" @cancel="showDisableModal = false" />
    </n-modal>

    <!-- 备用码显示弹窗 -->
    <n-modal
      v-model:show="showBackupCodesModal"
      preset="card"
      title="备用恢复码"
      :style="{ width: isMobile ? '95%' : '500px' }"
    >
      <BackupCodesDisplay :backup-codes="backupCodes" @close="showBackupCodesModal = false" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import { useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import { useBreakpoint } from '@/hooks/event/useBreakpoint';
  import { useTwoFactorStore } from '@/store/modules/twoFactor';
  import { useUserStore } from '@/store/modules/user';
  import { sizeEnum } from '@/enums/breakpointEnum';
  import {
    ShieldCheckmarkOutline,
    ShieldOutline,
    EyeOutline,
    RefreshOutline,
    KeyOutline,
    PhonePortraitOutline,
  } from '@vicons/ionicons5';
  import Enable2FAWizard from '@/views/system/security/components/Enable2FAWizard.vue';
  import Disable2FAForm from '@/views/system/security/components/Disable2FAForm.vue';
  import BackupCodesDisplay from '@/views/system/security/components/BackupCodesDisplay.vue';

  // 响应式断点检测
  const breakpoint = useBreakpoint();
  const isMobile = computed(
    () => breakpoint.screenRef.value === sizeEnum.XS || breakpoint.screenRef.value === sizeEnum.SM
  );
  const router = useRouter();
  const message = useMessage();
  const twoFactorStore = useTwoFactorStore();

  // 状态管理
  const loading = ref(false);
  const regenerateLoading = ref(false);
  const disableLoading = ref(false);
  const showingBackupCodes = ref(false);

  // 弹窗状态
  const showEnableWizard = ref(false);
  const showDisableModal = ref(false);
  const showBackupCodesModal = ref(false);
  const backupCodes = ref<string[]>([]);

  // 计算属性
  const twoFactorStatus = computed(() => twoFactorStore.status);

  /**
   * 获取备用码状态颜色
   */
  const getBackupCodesColor = () => {
    const count = twoFactorStatus.value?.backupCodesCount || 0;
    if (count === 0) return '#f5222d';
    if (count <= 2) return '#fa8c16';
    return '#52c41a';
  };

  /**
   * 格式化日期
   */
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-CN');
  };

  /**
   * 加载2FA状态
   */
  const loadTwoFactorStatus = async () => {
    try {
      loading.value = true;
      
      // 检查用户是否已登录
      const userStore = useUserStore();
      const token = userStore.getToken;
      
      if (!token) {
        message.error('请先登录后再访问此页面');
        router.push('/login');
        return;
      }
      
      const success = await twoFactorStore.fetchStatus();
      if (!success) {
        // 如果获取状态失败，检查是否是认证问题
        const errorMsg = twoFactorStore.errors.general;
        if (errorMsg.includes('用户ID不能为空') || errorMsg.includes('用户未登录')) {
          message.error('登录状态已过期，请重新登录');
          userStore.logout();
          router.push('/login');
        } else {
          message.error(`获取2FA状态失败: ${errorMsg}`);
        }
      }
    } catch (error: any) {
      console.error('Failed to load 2FA status:', error);
      const errorMsg = error?.message || error?.data?.message || '获取2FA状态失败';
      
      if (errorMsg.includes('用户ID不能为空') || errorMsg.includes('用户未登录')) {
        message.error('登录状态已过期，请重新登录');
        const userStore = useUserStore();
        userStore.logout();
        router.push('/login');
      } else {
        message.error(errorMsg);
      }
    } finally {
      loading.value = false;
    }
  };

  /**
   * 处理启用2FA
   */
  const handleEnable2FA = () => {
    showEnableWizard.value = true;
  };

  /**
   * 处理启用成功
   */
  const handleEnableSuccess = () => {
    showEnableWizard.value = false;
    message.success('双因素认证启用成功');
    loadTwoFactorStatus();
  };

  /**
   * 处理禁用2FA
   */
  const handleDisable2FA = () => {
    showDisableModal.value = true;
  };

  /**
   * 处理禁用成功
   */
  const handleDisableSuccess = () => {
    showDisableModal.value = false;
    message.success('双因素认证已禁用');
    loadTwoFactorStatus();
  };

  /**
   * 处理查看备用码
   */
  const handleShowBackupCodes = async () => {
    // 这里需要实现密码验证逻辑
    // 暂时直接显示，后续会在组件内部实现密码验证
    showBackupCodesModal.value = true;
  };

  /**
   * 处理重新生成备用码
   */
  const handleRegenerateBackupCodes = async () => {
    try {
      regenerateLoading.value = true;
      const codes = await twoFactorStore.regenerateBackupCodesAuth('');
      if (codes) {
        backupCodes.value = codes;
      } else {
        message.error('重新生成备用码失败');
        return;
      }
      showBackupCodesModal.value = true;
      message.success('备用码重新生成成功');
    } catch (error) {
      message.error('重新生成备用码失败');
      console.error('Failed to regenerate backup codes:', error);
    } finally {
      regenerateLoading.value = false;
    }
  };

  /**
   * 处理返回
   */
  const handleBack = () => {
    router.back();
  };

  // 组件挂载时加载状态
  onMounted(() => {
    loadTwoFactorStatus();
  });
</script>

<style scoped>
  .two-factor-settings {
    padding: 0;
    min-height: 100vh;
    background: #f5f5f5;
  }

  .page-header {
    background: white;
    padding: 16px 24px;
    border-bottom: 1px solid #e8e8e8;
    margin-bottom: 24px;
  }

  .main-content {
    padding: 0 24px 24px;
  }

  .status-card {
    min-height: 400px;
  }

  .status-overview {
    margin-bottom: 32px;
  }

  .status-item {
    display: flex;
    align-items: flex-start;
    gap: 16px;
    padding: 20px;
    border-radius: 8px;
    background: #fafafa;
    border: 1px solid #e8e8e8;
  }

  .status-item.primary {
    background: linear-gradient(135deg, #f6ffed 0%, #f0f9ff 100%);
    border-color: #d9f7be;
  }

  .status-item.secondary {
    background: #f0f9ff;
    border-color: #bae7ff;
  }

  .status-icon {
    flex-shrink: 0;
  }

  .status-content {
    flex: 1;
  }

  .status-title {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    margin-bottom: 4px;
  }

  .status-desc {
    font-size: 14px;
    color: #595959;
    line-height: 1.5;
    margin-bottom: 8px;
  }

  .status-meta {
    font-size: 12px;
    color: #8c8c8c;
  }

  .action-section {
    padding-top: 24px;
    border-top: 1px solid #e8e8e8;
  }

  .primary-action {
    height: 48px;
    font-size: 16px;
    font-weight: 600;
  }

  .help-card {
    height: fit-content;
  }

  .help-content {
    font-size: 14px;
    line-height: 1.6;
  }

  .help-content ul,
  .help-content ol {
    margin: 8px 0;
    padding-left: 20px;
  }

  .help-content li {
    margin: 4px 0;
  }

  .app-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .app-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px;
    border-radius: 6px;
    background: #fafafa;
  }

  .app-name {
    font-weight: 600;
    color: #262626;
  }

  .app-desc {
    font-size: 12px;
    color: #8c8c8c;
  }

  .security-tips {
    margin: 0;
    padding-left: 16px;
  }

  .security-tips li {
    margin: 8px 0;
    color: #fa8c16;
  }

  .enable-wizard-modal {
    max-height: 90vh;
    overflow-y: auto;
  }

  /* 移动端适配 */
  @media (max-width: 768px) {
    .main-content {
      padding: 0 16px 16px;
    }

    .page-header {
      padding: 12px 16px;
    }

    .status-item {
      padding: 16px;
      flex-direction: column;
      text-align: center;
      gap: 12px;
    }

    .status-item.primary {
      flex-direction: row;
      text-align: left;
    }

    .action-section {
      padding-top: 20px;
    }

    .primary-action {
      width: 100%;
      height: 44px;
    }
  }
</style>
