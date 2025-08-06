<template>
  <div class="disable-2fa-form">
    <n-alert type="warning" title="警告" style="margin-bottom: 24px">
      禁用双因素认证将降低您账户的安全性。请确认您真的要禁用此功能。
    </n-alert>

    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
      <n-form-item path="password" label="当前密码">
        <n-input
          v-model:value="formData.password"
          type="password"
          placeholder="请输入当前密码以确认身份"
          show-password-on="click"
        />
      </n-form-item>

      <n-form-item path="code" label="验证码">
        <n-input
          v-model:value="formData.code"
          placeholder="请输入验证器应用显示的6位验证码"
          maxlength="6"
          @keyup.enter="handleSubmit"
        />
        <template #feedback>
          <span class="form-tip">请打开您的验证器应用，输入显示的6位数字</span>
        </template>
      </n-form-item>
    </n-form>

    <div class="form-actions">
      <n-space>
        <n-button @click="emit('cancel')">取消</n-button>
        <n-button type="error" @click="handleSubmit" :loading="twoFactorStore.loading.disable">
          确认禁用
        </n-button>
      </n-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useMessage } from 'naive-ui';
  import { useTwoFactorStore } from '@/store/modules/twoFactor';
  import type { FormInst } from 'naive-ui';

  // 事件定义
  const emit = defineEmits<{
    success: [];
    cancel: [];
  }>();

  // Store和工具
  const twoFactorStore = useTwoFactorStore();
  const message = useMessage();
  const formRef = ref<FormInst | null>(null);

  // 表单数据
  const formData = ref({
    password: '',
    code: '',
  });

  // 表单验证规则
  const formRules = {
    password: {
      required: true,
      message: '请输入当前密码',
      trigger: 'blur',
    },
    code: {
      required: true,
      message: '请输入验证码',
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
  };

  /**
   * 处理表单提交
   */
  const handleSubmit = async () => {
    if (!formRef.value) return;

    try {
      await formRef.value.validate();
    } catch {
      return;
    }

    const success = await twoFactorStore.disable2FAAuth(
      formData.value.password,
      formData.value.code
    );

    if (success) {
      message.success('双因素认证已成功禁用');
      formData.value = { password: '', code: '' }; // 清空表单
      emit('success');
    }
  };
</script>

<style lang="less" scoped>
  .disable-2fa-form {
    .form-tip {
      font-size: 12px;
      color: #999;
    }

    .form-actions {
      margin-top: 24px;
      padding-top: 16px;
      border-top: 1px solid #f0f0f0;
      text-align: right;
    }
  }
</style>
