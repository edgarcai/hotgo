<template>
  <div class="verify-2fa-container">
    <div class="verify-2fa-card">
      <n-card :bordered="false">
        <!-- 头部 -->
        <header class="verify-header">
          <n-space justify="center" align="center" :size="16">
            <n-icon size="32" color="#18a058">
              <ShieldCheckmarkOutline />
            </n-icon>
            <div>
              <h2 class="verify-title">双因素认证</h2>
              <p class="verify-subtitle">请完成身份验证以继续登录</p>
            </div>
          </n-space>
        </header>

        <!-- 用户信息 -->
        <div class="user-info">
          <n-space justify="center" align="center" :size="12">
            <n-avatar :size="40" :src="userAvatar" />
            <div>
              <div class="username">{{ username }}</div>
              <div class="user-hint">正在登录到 {{ projectName }}</div>
            </div>
          </n-space>
        </div>

        <!-- 2FA输入组件 -->
        <div class="two-factor-section">
          <TwoFactorInput
            ref="twoFactorInputRef"
            :mode="verificationMode"
            :loading="loading"
            :error="errorMessage"
            @submit="handleVerificationSubmit"
            @switch-mode="handleSwitchMode"
            @update:error="errorMessage = $event"
          />
        </div>

        <!-- 底部操作 -->
        <div class="footer-actions">
          <n-space justify="space-between" align="center">
            <n-button text @click="handleBackToLogin">
              <template #icon>
                <n-icon><ArrowBackOutline /></n-icon>
              </template>
              返回登录
            </n-button>

            <n-button text @click="handleNeedHelp"> 需要帮助？ </n-button>
          </n-space>
        </div>
      </n-card>
    </div>

    <!-- 帮助对话框 -->
    <n-modal v-model:show="showHelpModal" preset="dialog" title="双因素认证帮助">
      <div class="help-content">
        <h4>如何获取验证码？</h4>
        <ul>
          <li>打开您的身份验证器应用（如 Google Authenticator、Microsoft Authenticator）</li>
          <li>找到对应的账户条目</li>
          <li>输入显示的6位数字验证码</li>
        </ul>

        <h4>如果无法访问验证器应用？</h4>
        <ul>
          <li>点击"使用备用恢复码"</li>
          <li>输入您在设置2FA时保存的备用恢复码</li>
          <li>每个恢复码只能使用一次</li>
        </ul>

        <h4>仍然无法登录？</h4>
        <p>请联系系统管理员寻求帮助。</p>
      </div>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, computed } from 'vue';
  import { useRouter, useRoute } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import { ShieldCheckmarkOutline, ArrowBackOutline } from '@vicons/ionicons5';
  import TwoFactorInput from '@/components/TwoFactor/TwoFactorInput.vue';
  import { useUserStore } from '@/store/modules/user';
  import { verifyLogin2FA } from '@/api/auth/twoFactor';
  import { ResultEnum } from '@/enums/httpEnum';
  import { PageEnum } from '@/enums/pageEnum';
  import { getAppEnvConfig } from '@/utils/env';

  // 响应式数据
  const router = useRouter();
  const route = useRoute();
  const message = useMessage();
  const userStore = useUserStore();

  const loading = ref(false);
  const errorMessage = ref<string | undefined>();
  const verificationMode = ref<'totp' | 'backup'>('totp');
  const showHelpModal = ref(false);
  const twoFactorInputRef = ref();

  // 从路由参数或store获取临时token和用户信息
  const tempToken = ref((route.query.tempToken as string) || '');
  const username = ref((route.query.username as string) || '');
  const userAvatar = ref((route.query.avatar as string) || '');

  // 计算属性
  const projectName = computed(() => {
    return getAppEnvConfig().VITE_GLOB_APP_TITLE || 'HotGo';
  });

  const redirectPath = computed(() => {
    const r = (route.query?.redirect as string) || '/';
    try {
      return decodeURIComponent(r);
    } catch (e) {
      return '/';
    }
  });

  // 方法
  const handleVerificationSubmit = async (code: string) => {
    if (!tempToken.value) {
      message.error('临时token无效，请重新登录');
      handleBackToLogin();
      return;
    }

    loading.value = true;
    errorMessage.value = undefined;

    try {
      const params = {
        tempToken: tempToken.value,
        code: code,
        type: (verificationMode.value === 'totp' ? 'totp' : 'backup') as 'totp' | 'backup',
      };

      const response = await verifyLogin2FA(params);

      if (response.code === ResultEnum.SUCCESS) {
        // 验证成功，保存token，并拉取用户信息
        const { token } = response.data;

        // 更新用户store
        userStore.setToken(token);
        // 保存到localStorage
        const ex = 30 * 24 * 60 * 60 * 1000;
        localStorage.setItem('ACCESS_TOKEN', token);
        // 拉取并设置用户信息
        try {
          const info = await userStore.GetInfo();
          localStorage.setItem('CURRENT_USER', JSON.stringify(info));
        } catch (e) {
          // 忽略错误，后续进入应用会再次拉取
          console.warn('获取用户信息失败，将在进入应用后重试');
        }

        message.success('验证成功，正在跳转...');

        // 跳转到目标页面
        await router.replace(redirectPath.value || '/');
      } else {
        errorMessage.value = response.message || '验证失败，请重试';

        // 清空输入
        if (twoFactorInputRef.value) {
          twoFactorInputRef.value.clear();
          twoFactorInputRef.value.focus();
        }
      }
    } catch (error: any) {
      console.error('2FA验证失败:', error);
      errorMessage.value = error.message || '验证失败，请重试';

      // 清空输入
      if (twoFactorInputRef.value) {
        twoFactorInputRef.value.clear();
        twoFactorInputRef.value.focus();
      }
    } finally {
      loading.value = false;
    }
  };

  const handleSwitchMode = () => {
    verificationMode.value = verificationMode.value === 'totp' ? 'backup' : 'totp';
    errorMessage.value = undefined;
  };

  const handleBackToLogin = () => {
    router.push({
      path: PageEnum.BASE_LOGIN,
      query: {
        redirect: redirectPath.value || '/',
      },
    });
  };

  const handleNeedHelp = () => {
    showHelpModal.value = true;
  };

  // 生命周期
  onMounted(() => {
    // 检查是否有必要的参数
    if (!tempToken.value) {
      message.error('缺少必要的验证参数，请重新登录');
      handleBackToLogin();
      return;
    }

    // 自动聚焦输入框
    if (twoFactorInputRef.value) {
      twoFactorInputRef.value.focus();
    }
  });
</script>

<style lang="less" scoped>
  .verify-2fa-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
  }

  .verify-2fa-card {
    width: 100%;
    max-width: 420px;

    :deep(.n-card) {
      border-radius: 12px;
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
    }

    :deep(.n-card__content) {
      padding: 32px;
    }
  }

  .verify-header {
    text-align: center;
    margin-bottom: 24px;

    .verify-title {
      margin: 0;
      font-size: 24px;
      font-weight: 600;
      color: #333;
    }

    .verify-subtitle {
      margin: 4px 0 0 0;
      font-size: 14px;
      color: #666;
    }
  }

  .user-info {
    text-align: center;
    margin-bottom: 32px;
    padding: 16px;
    background-color: #f8f9fa;
    border-radius: 8px;

    .username {
      font-size: 16px;
      font-weight: 500;
      color: #333;
      margin-bottom: 2px;
    }

    .user-hint {
      font-size: 12px;
      color: #666;
    }
  }

  .two-factor-section {
    margin-bottom: 24px;
  }

  .footer-actions {
    margin-top: 24px;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;
  }

  .help-content {
    h4 {
      margin: 16px 0 8px 0;
      color: #333;
      font-size: 14px;

      &:first-child {
        margin-top: 0;
      }
    }

    ul {
      margin: 0 0 16px 0;
      padding-left: 20px;

      li {
        margin-bottom: 4px;
        font-size: 13px;
        color: #666;
        line-height: 1.4;
      }
    }

    p {
      margin: 0;
      font-size: 13px;
      color: #666;
      line-height: 1.4;
    }
  }

  // 响应式设计
  @media (max-width: 768px) {
    .verify-2fa-container {
      padding: 16px;
    }

    .verify-2fa-card {
      :deep(.n-card__content) {
        padding: 24px 20px;
      }
    }

    .verify-header {
      .verify-title {
        font-size: 20px;
      }
    }

    .user-info {
      padding: 12px;

      .username {
        font-size: 14px;
      }
    }
  }
</style>
