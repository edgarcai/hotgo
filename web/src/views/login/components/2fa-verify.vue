<template>
  <div class="2fa-verify-container">
    <n-card :bordered="false" class="2fa-verify-card">
      <template #header>
        <div class="text-center">
          <n-icon size="48" color="#18a058" class="mb-4">
            <ShieldCheckmarkOutline />
          </n-icon>
          <h3 class="text-lg font-semibold">双因子认证</h3>
          <p class="text-gray-500 text-sm mt-2">请输入身份验证器中的6位数字验证码</p>
        </div>
      </template>
      
      <n-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleVerify">
        <n-form-item path="code">
          <n-input
            v-model:value="form.code"
            placeholder="请输入6位验证码"
            size="large"
            maxlength="6"
            :input-props="{ autocomplete: 'off', inputmode: 'numeric' }"
            @keyup.enter="handleVerify"
            class="text-center text-2xl tracking-widest"
          >
            <template #prefix>
              <n-icon size="18" color="#808695">
                <SafetyCertificateOutlined />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>
        
        <n-space vertical size="large">
          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            @click="handleVerify"
            :disabled="form.code.length !== 6"
          >
            验证
          </n-button>
          
          <div class="text-center">
            <n-button text @click="showBackupCodeInput = !showBackupCodeInput">
              使用备用码登录
            </n-button>
          </div>
          
          <!-- 备用码输入 -->
          <div v-if="showBackupCodeInput" class="backup-code-section">
            <n-divider>备用码登录</n-divider>
            <n-form-item path="backupCode">
              <n-input
                v-model:value="form.backupCode"
                placeholder="请输入备用码"
                size="large"
                :input-props="{ autocomplete: 'off' }"
                @keyup.enter="handleBackupCodeVerify"
              >
                <template #prefix>
                  <n-icon size="18" color="#808695">
                    <KeyOutline />
                  </n-icon>
                </template>
              </n-input>
            </n-form-item>
            <n-button
              type="primary"
              size="large"
              block
              :loading="loading"
              @click="handleBackupCodeVerify"
              :disabled="!form.backupCode.trim()"
            >
              使用备用码验证
            </n-button>
          </div>
          
          <div class="text-center">
            <n-button text @click="handleBack">
              返回登录
            </n-button>
          </div>
        </n-space>
      </n-form>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive } from 'vue';
  import { useMessage } from 'naive-ui';
  import { ShieldCheckmarkOutline, KeyOutline } from '@vicons/ionicons5';
  import { SafetyCertificateOutlined } from '@vicons/antd';
  import { verify2FA, useBackupCode } from '@/api/system/2fa';
  import { useUserStore } from '@/store/modules/user';
  import { useRouter, useRoute } from 'vue-router';
  import { ResultEnum } from '@/enums/httpEnum';
  import { PageEnum } from '@/enums/pageEnum';

  interface Props {
    loginData?: any;
  }

  const props = withDefaults(defineProps<Props>(), {
    loginData: null,
  });

  const emit = defineEmits(['back', 'success']);

  const formRef = ref();
  const message = useMessage();
  const loading = ref(false);
  const showBackupCodeInput = ref(false);
  const userStore = useUserStore();
  const router = useRouter();
  const route = useRoute();

  const form = reactive({
    code: '',
    backupCode: '',
  });

  const rules = {
    code: {
      required: true,
      message: '请输入6位验证码',
      trigger: 'blur',
      validator: (rule: any, value: string) => {
        if (!value) {
          return new Error('请输入验证码');
        }
        if (!/^\d{6}$/.test(value)) {
          return new Error('验证码必须是6位数字');
        }
        return true;
      },
    },
    backupCode: {
      required: true,
      message: '请输入备用码',
      trigger: 'blur',
    },
  };

  async function handleVerify() {
    if (!form.code || form.code.length !== 6) {
      message.error('请输入6位验证码');
      return;
    }

    loading.value = true;
    try {
      const response = await verify2FA({ code: form.code });
      if (response.code === ResultEnum.SUCCESS) {
        message.success('验证成功');
        await completeLogin();
      } else {
        message.error(response.message || '验证失败');
        form.code = '';
      }
    } catch (error) {
      message.error('验证失败，请重试');
      form.code = '';
    } finally {
      loading.value = false;
    }
  }

  async function handleBackupCodeVerify() {
    if (!form.backupCode.trim()) {
      message.error('请输入备用码');
      return;
    }

    loading.value = true;
    try {
      const response = await useBackupCode({ code: form.backupCode });
      if (response.code === ResultEnum.SUCCESS) {
        message.success('验证成功');
        await completeLogin();
      } else {
        message.error(response.message || '备用码无效');
        form.backupCode = '';
      }
    } catch (error) {
      message.error('验证失败，请重试');
      form.backupCode = '';
    } finally {
      loading.value = false;
    }
  }

  async function completeLogin() {
    try {
      // 获取用户信息
      await userStore.GetInfo();
      
      const toPath = decodeURIComponent((route.query?.redirect || '/') as string);
      message.success('登录成功，即将进入系统');
      
      if (route.name === PageEnum.BASE_LOGIN_NAME) {
        await router.replace('/');
      } else {
        await router.replace(toPath);
      }
      
      emit('success');
    } catch (error) {
      message.error('获取用户信息失败');
    }
  }

  function handleBack() {
    emit('back');
  }
</script>

<style lang="less" scoped>
  .2fa-verify-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  .2fa-verify-card {
    width: 100%;
    max-width: 400px;
    margin: 0 auto;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
  }

  .backup-code-section {
    margin-top: 16px;
    padding-top: 16px;
  }

  :deep(.n-input__input-el) {
    text-align: center;
    font-size: 18px;
    letter-spacing: 4px;
    font-weight: 600;
  }

  :deep(.n-divider__title) {
    font-size: 14px;
    color: #666;
  }
</style>