<template>
  <div class="enable-2fa-wizard">
    <n-steps :current="twoFactorStore.enableFlow.step" :status="stepStatus">
      <n-step title="验证密码" description="确认您的身份" />
      <n-step title="扫描二维码" description="使用验证器应用扫描" />
      <n-step title="验证设置" description="输入验证码完成设置" />
      <n-step title="保存备用码" description="保存恢复码以备不时之需" />
    </n-steps>

    <div class="step-content">
      <!-- 步骤1: 验证密码 -->
      <div v-if="twoFactorStore.enableFlow.step === 1" class="step-panel">
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

        <div v-if="twoFactorStore.errors.password" class="error-message">
          {{ twoFactorStore.errors.password }}
        </div>

        <div class="step-actions">
          <n-space>
            <n-button @click="handleCancel">取消</n-button>
            <n-button
              type="primary"
              @click="handlePasswordVerify"
              :loading="twoFactorStore.loading.enable"
            >
              验证密码
            </n-button>
          </n-space>
        </div>
      </div>

      <!-- 步骤2: 扫描二维码 -->
      <div v-if="twoFactorStore.enableFlow.step === 2" class="step-panel">
        <div class="step-title">扫描二维码</div>
        <div class="step-desc">使用您的验证器应用扫描下方二维码，然后输入应用显示的6位验证码。</div>

        <div class="qr-section">
          <div class="qr-code">
            <img
              v-if="twoFactorStore.enableFlow.qrCodeUrl"
              :src="twoFactorStore.enableFlow.qrCodeUrl"
              alt="2FA QR Code"
              @click="handleQRCodeClick"
              style="cursor: pointer"
            />
            <n-spin v-else size="large" />
          </div>

          <div class="qr-info">
            <div class="secret-key">
              <div class="label">手动输入密钥（如果无法扫描）：</div>
              <div class="key-value">
                <n-input readonly :value="twoFactorStore.enableFlow.secretKey" />
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
            <n-button @click="twoFactorStore.setEnableStep(1)">上一步</n-button>
            <n-button type="primary" @click="twoFactorStore.setEnableStep(3)">
              我已扫描二维码
            </n-button>
          </n-space>
        </div>
      </div>

      <!-- 步骤3: 验证设置 -->
      <div v-if="twoFactorStore.enableFlow.step === 3" class="step-panel">
        <div class="step-title">验证设置</div>
        <div class="step-desc">请输入验证器应用显示的6位验证码以完成设置。</div>

        <div class="verify-section">
          <TwoFactorInput
            v-model:code="verificationCode"
            :mode="'totp'"
            :error="twoFactorStore.errors.verification"
            @submit="handleVerifySetup"
          />
        </div>

        <div class="step-actions">
          <n-space>
            <n-button @click="twoFactorStore.setEnableStep(2)">上一步</n-button>
            <n-button
              type="primary"
              @click="handleVerifySetup"
              :loading="twoFactorStore.loading.verify"
            >
              验证并启用
            </n-button>
          </n-space>
        </div>
      </div>

      <!-- 步骤4: 保存备用码 -->
      <div v-if="twoFactorStore.enableFlow.step === 4" class="step-panel">
        <div class="step-title">保存备用恢复码</div>
        <div class="step-desc"
          >请将以下备用恢复码保存在安全的地方。当您无法使用验证器应用时，可以使用这些代码登录。</div
        >

        <div class="backup-codes-section">
          <BackupCodesDisplay :backup-codes="twoFactorStore.enableFlow.backupCodes" />
        </div>

        <div class="security-warning">
          <n-alert type="warning" title="重要提示">
            <ul>
              <li>每个恢复码只能使用一次</li>
              <li>请将恢复码保存在安全的地方，不要与账户存储在同一设备上</li>
              <li>如果丢失恢复码，您可以在设置中重新生成</li>
              <li>建议打印或写在纸上，存放在安全位置</li>
            </ul>
          </n-alert>
        </div>

        <div class="confirmation-section">
          <n-checkbox v-model:checked="backupCodesSaved"> 我已安全保存了备用恢复码 </n-checkbox>
        </div>

        <div class="step-actions">
          <n-space>
            <n-button @click="handleComplete" type="primary" :disabled="!backupCodesSaved">
              完成设置
            </n-button>
          </n-space>
        </div>
      </div>
    </div>

    <!-- QR码放大查看弹窗 -->
    <n-modal v-model:show="showQRModal" preset="card" title="二维码" style="width: 400px">
      <div class="qr-modal-content">
        <img
          v-if="twoFactorStore.enableFlow.qrCodeUrl"
          :src="twoFactorStore.enableFlow.qrCodeUrl"
          alt="2FA QR Code"
          class="qr-modal-image"
        />
        <div class="qr-modal-hint"> 使用验证器应用扫描此二维码 </div>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue';
  import { useMessage } from 'naive-ui';
  import { CopyOutline } from '@vicons/ionicons5';
  import { useTwoFactorStore } from '@/store/modules/twoFactor';
  import TwoFactorInput from '@/components/TwoFactor/TwoFactorInput.vue';
  import BackupCodesDisplay from './BackupCodesDisplay.vue';

  // 组件属性
  const emit = defineEmits<{
    success: [];
    cancel: [];
  }>();

  // Store和工具
  const twoFactorStore = useTwoFactorStore();
  const message = useMessage();

  // 响应式数据
  const verificationCode = ref('');
  const backupCodesSaved = ref(false);
  const showQRModal = ref(false);

  // 表单数据
  const passwordForm = ref({
    password: '',
  });

  // 表单验证规则
  const passwordRules = {
    password: {
      required: true,
      message: '请输入密码',
      trigger: 'blur',
    },
  };

  // 计算属性
  const stepStatus = computed(() => {
    return twoFactorStore.enableFlow.step === 4 ? 'finish' : 'process';
  });

  /**
   * 处理密码验证
   */
  const handlePasswordVerify = async () => {
    const success = await twoFactorStore.startEnable2FA(passwordForm.value.password);
    if (success) {
      message.success('密码验证成功');
      passwordForm.value.password = ''; // 清空密码
    }
  };

  /**
   * 处理验证设置
   */
  const handleVerifySetup = async () => {
    if (!verificationCode.value || verificationCode.value.length !== 6) {
      twoFactorStore.setError('verification', '请输入6位验证码');
      return;
    }

    const success = await twoFactorStore.verifyEnable2FA(verificationCode.value);
    if (success) {
      message.success('2FA设置验证成功');
      verificationCode.value = '';
    }
  };

  /**
   * 复制密钥
   */
  const copySecretKey = async () => {
    try {
      await navigator.clipboard.writeText(twoFactorStore.enableFlow.secretKey);
      message.success('密钥已复制到剪贴板');
    } catch {
      message.error('复制失败，请手动复制');
    }
  };

  /**
   * 处理QR码点击
   */
  const handleQRCodeClick = () => {
    showQRModal.value = true;
  };

  /**
   * 处理完成设置
   */
  const handleComplete = () => {
    twoFactorStore.completeEnableFlow();
    emit('success');
  };

  /**
   * 处理取消
   */
  const handleCancel = () => {
    twoFactorStore.cancelEnableFlow();
    emit('cancel');
  };
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
