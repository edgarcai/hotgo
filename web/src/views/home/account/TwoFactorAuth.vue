<template>
  <div class="two-factor-auth">
    <n-card :bordered="false" class="mb-4">
      <template #header>
        <n-space align="center">
          <n-icon size="20" color="#18a058">
            <ShieldCheckmarkOutline />
          </n-icon>
          <span>双因子认证状态</span>
        </n-space>
      </template>
      
      <n-space vertical size="large">
        <n-alert
          v-if="!status.isEnabled"
          type="warning"
          title="未启用双因子认证"
          show-icon
        >
          为了提高账户安全性，建议启用双因子认证。启用后，登录时需要输入手机验证器中的6位数字验证码。
        </n-alert>
        
        <n-alert
          v-else
          type="success"
          title="双因子认证已启用"
          show-icon
        >
          您的账户已启用双因子认证，安全性得到有效保障。
        </n-alert>
        
        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="认证状态">
            <n-tag :type="status.isEnabled ? 'success' : 'warning'">
              {{ status.isEnabled ? '已启用' : '未启用' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="验证状态">
            <n-tag :type="status.isVerified ? 'success' : 'default'">
              {{ status.isVerified ? '已验证' : '未验证' }}
            </n-tag>
          </n-descriptions-item>
        </n-descriptions>
      </n-space>
    </n-card>
    
    <!-- 启用2FA -->
    <n-card v-if="!status.isEnabled" :bordered="false" class="mb-4">
      <template #header>
        <n-space align="center">
          <n-icon size="20" color="#2080f0">
            <SettingsOutline />
          </n-icon>
          <span>启用双因子认证</span>
        </n-space>
      </template>
      
      <n-steps :current="setupStep" size="small" class="mb-6">
        <n-step title="生成密钥" description="获取TOTP密钥和二维码" />
        <n-step title="配置应用" description="使用验证器应用扫描二维码" />
        <n-step title="验证设置" description="输入验证码完成设置" />
      </n-steps>
      
      <!-- 步骤1：生成密钥 -->
      <div v-if="setupStep === 1">
        <n-space vertical size="large">
          <n-alert type="info" show-icon>
            点击下方按钮生成TOTP密钥和二维码，用于配置您的身份验证器应用。
          </n-alert>
          <n-button
            type="primary"
            size="large"
            :loading="loading"
            @click="generateSecret"
          >
            生成密钥和二维码
          </n-button>
        </n-space>
      </div>
      
      <!-- 步骤2：显示二维码 -->
      <div v-if="setupStep === 2">
        <n-space vertical size="large">
          <n-alert type="info" show-icon>
            请使用Google Authenticator、Microsoft Authenticator等应用扫描下方二维码，或手动输入密钥。
          </n-alert>
          
          <n-grid cols="1 600:2" :x-gap="24">
            <n-grid-item>
              <n-card title="扫描二维码" size="small">
                <div class="text-center">
                  <img :src="setupData.qrCode" alt="2FA QR Code" class="qr-code" />
                </div>
              </n-card>
            </n-grid-item>
            <n-grid-item>
              <n-card title="手动输入" size="small">
                <n-space vertical>
                  <div>
                    <n-text strong>密钥：</n-text>
                    <n-text code>{{ setupData.secretKey }}</n-text>
                    <n-button
                      text
                      size="small"
                      @click="copySecret"
                      class="ml-2"
                    >
                      复制
                    </n-button>
                  </div>
                  <n-text depth="3" class="text-sm">
                    如果无法扫描二维码，请在验证器应用中手动添加账户，并输入上述密钥。
                  </n-text>
                </n-space>
              </n-card>
            </n-grid-item>
          </n-grid>
          
          <n-space>
            <n-button @click="setupStep = 1">上一步</n-button>
            <n-button type="primary" @click="setupStep = 3">下一步</n-button>
          </n-space>
        </n-space>
      </div>
      
      <!-- 步骤3：验证设置 -->
      <div v-if="setupStep === 3">
        <n-space vertical size="large">
          <n-alert type="info" show-icon>
            请输入验证器应用中显示的6位数字验证码来完成设置。
          </n-alert>
          
          <n-form ref="verifyFormRef" :model="verifyForm" :rules="verifyRules">
            <n-form-item path="code" label="验证码">
              <n-input
                v-model:value="verifyForm.code"
                placeholder="请输入6位验证码"
                maxlength="6"
                :input-props="{ autocomplete: 'off', inputmode: 'numeric' }"
                @keyup.enter="verifySetup"
              />
            </n-form-item>
          </n-form>
          
          <n-space>
            <n-button @click="setupStep = 2">上一步</n-button>
            <n-button
              type="primary"
              :loading="loading"
              @click="verifySetup"
              :disabled="verifyForm.code.length !== 6"
            >
              完成设置
            </n-button>
          </n-space>
        </n-space>
      </div>
    </n-card>
    
    <!-- 管理2FA -->
    <n-card v-if="status.isEnabled" :bordered="false" class="mb-4">
      <template #header>
        <n-space align="center">
          <n-icon size="20" color="#f0a020">
            <SettingsOutline />
          </n-icon>
          <span>管理双因子认证</span>
        </n-space>
      </template>
      
      <n-space vertical size="large">
        <n-grid cols="1 600:2" :x-gap="24">
          <n-grid-item>
            <n-card title="备用码管理" size="small">
              <n-space vertical>
                <n-text>
                  备用码可在无法使用验证器应用时用于登录，每个备用码只能使用一次。
                </n-text>
                <n-space>
                  <n-button @click="showBackupCodes">查看备用码</n-button>
                  <n-button @click="showRegenerateModal = true">重新生成</n-button>
                </n-space>
              </n-space>
            </n-card>
          </n-grid-item>
          <n-grid-item>
            <n-card title="禁用认证" size="small">
              <n-space vertical>
                <n-text>
                  禁用双因子认证将降低账户安全性，请谨慎操作。
                </n-text>
                <n-button type="error" @click="showDisableModal = true">
                  禁用双因子认证
                </n-button>
              </n-space>
            </n-card>
          </n-grid-item>
        </n-grid>
      </n-space>
    </n-card>
    
    <!-- 备用码显示模态框 -->
    <n-modal v-model:show="showBackupCodesModal" preset="card" title="备用码" style="width: 500px">
      <n-space vertical>
        <n-alert type="warning" show-icon>
          请将这些备用码保存在安全的地方。每个备用码只能使用一次。
        </n-alert>
        <n-grid cols="2" :x-gap="12" :y-gap="8">
          <n-grid-item v-for="code in backupCodes" :key="code">
            <n-input :value="code" readonly>
              <template #suffix>
                <n-button text size="small" @click="copyText(code)">
                  复制
                </n-button>
              </template>
            </n-input>
          </n-grid-item>
        </n-grid>
      </n-space>
    </n-modal>
    
    <!-- 重新生成备用码模态框 -->
    <n-modal v-model:show="showRegenerateModal" preset="card" title="重新生成备用码" style="width: 400px">
      <n-space vertical>
        <n-alert type="warning" show-icon>
          重新生成备用码将使现有的备用码失效。请输入验证码确认操作。
        </n-alert>
        <n-form ref="regenerateFormRef" :model="regenerateForm" :rules="codeRules">
          <n-form-item path="code" label="验证码">
            <n-input
              v-model:value="regenerateForm.code"
              placeholder="请输入6位验证码"
              maxlength="6"
              :input-props="{ autocomplete: 'off', inputmode: 'numeric' }"
            />
          </n-form-item>
        </n-form>
        <n-space justify="end">
          <n-button @click="showRegenerateModal = false">取消</n-button>
          <n-button
            type="primary"
            :loading="loading"
            @click="regenerateBackupCodes"
            :disabled="regenerateForm.code.length !== 6"
          >
            确认重新生成
          </n-button>
        </n-space>
      </n-space>
    </n-modal>
    
    <!-- 禁用2FA模态框 -->
    <n-modal v-model:show="showDisableModal" preset="card" title="禁用双因子认证" style="width: 400px">
      <n-space vertical>
        <n-alert type="error" show-icon>
          禁用双因子认证将降低账户安全性。请输入验证码确认操作。
        </n-alert>
        <n-form ref="disableFormRef" :model="disableForm" :rules="codeRules">
          <n-form-item path="code" label="验证码">
            <n-input
              v-model:value="disableForm.code"
              placeholder="请输入6位验证码"
              maxlength="6"
              :input-props="{ autocomplete: 'off', inputmode: 'numeric' }"
            />
          </n-form-item>
        </n-form>
        <n-space justify="end">
          <n-button @click="showDisableModal = false">取消</n-button>
          <n-button
            type="error"
            :loading="loading"
            @click="disable2FA"
            :disabled="disableForm.code.length !== 6"
          >
            确认禁用
          </n-button>
        </n-space>
      </n-space>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, onMounted } from 'vue';
  import { useMessage } from 'naive-ui';
  import { ShieldCheckmarkOutline, SettingsOutline } from '@vicons/ionicons5';
  import {
    get2FAStatus,
    setup2FA,
    verifySetup2FA,
    disable2FA as disableAPI,
    getBackupCodes,
    regenerateBackupCodes as regenerateAPI,
  } from '@/api/system/2fa';
  import { ResultEnum } from '@/enums/httpEnum';

  const message = useMessage();
  const loading = ref(false);
  const setupStep = ref(1);
  const showBackupCodesModal = ref(false);
  const showRegenerateModal = ref(false);
  const showDisableModal = ref(false);

  const status = reactive({
    isEnabled: false,
    isVerified: false,
  });

  const setupData = reactive({
    secretKey: '',
    qrCode: '',
  });

  const verifyForm = reactive({
    code: '',
  });

  const regenerateForm = reactive({
    code: '',
  });

  const disableForm = reactive({
    code: '',
  });

  const backupCodes = ref<string[]>([]);

  const verifyRules = {
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
  };

  const codeRules = {
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
  };

  onMounted(() => {
    loadStatus();
  });

  async function loadStatus() {
    try {
      const response = await get2FAStatus();
      // 检查是否是错误响应（当isShowErrorMessage为false时返回的错误对象）
      if (response && response.error) {
        console.log('2FA功能未启用或接口不可用');
        return;
      }
      if (response.code === ResultEnum.SUCCESS) {
        Object.assign(status, response.data);
      }
    } catch (error) {
      console.error('获取2FA状态失败:', error);
    }
  }

  async function generateSecret() {
    loading.value = true;
    try {
      const response = await setup2FA();
      if (response.code === ResultEnum.SUCCESS) {
        Object.assign(setupData, response.data);
        setupStep.value = 2;
        message.success('密钥生成成功');
      } else {
        message.error(response.message || '生成密钥失败');
      }
    } catch (error) {
      message.error('生成密钥失败');
    } finally {
      loading.value = false;
    }
  }

  async function verifySetup() {
    loading.value = true;
    try {
      const response = await verifySetup2FA({ code: verifyForm.code });
      if (response.code === ResultEnum.SUCCESS) {
        message.success('双因子认证设置成功');
        await loadStatus();
        resetSetup();
      } else {
        message.error(response.message || '验证失败');
        verifyForm.code = '';
      }
    } catch (error) {
      message.error('验证失败');
      verifyForm.code = '';
    } finally {
      loading.value = false;
    }
  }

  async function showBackupCodes() {
    loading.value = true;
    try {
      const response = await getBackupCodes();
      if (response.code === ResultEnum.SUCCESS) {
        backupCodes.value = response.data.backupCodes;
        showBackupCodesModal.value = true;
      } else {
        message.error(response.message || '获取备用码失败');
      }
    } catch (error) {
      message.error('获取备用码失败');
    } finally {
      loading.value = false;
    }
  }

  async function regenerateBackupCodes() {
    loading.value = true;
    try {
      const response = await regenerateAPI({ code: regenerateForm.code });
      if (response.code === ResultEnum.SUCCESS) {
        backupCodes.value = response.data.backupCodes;
        showRegenerateModal.value = false;
        showBackupCodesModal.value = true;
        regenerateForm.code = '';
        message.success('备用码重新生成成功');
      } else {
        message.error(response.message || '重新生成失败');
        regenerateForm.code = '';
      }
    } catch (error) {
      message.error('重新生成失败');
      regenerateForm.code = '';
    } finally {
      loading.value = false;
    }
  }

  async function disable2FA() {
    loading.value = true;
    try {
      const response = await disableAPI({ code: disableForm.code });
      if (response.code === ResultEnum.SUCCESS) {
        message.success('双因子认证已禁用');
        await loadStatus();
        showDisableModal.value = false;
        disableForm.code = '';
      } else {
        message.error(response.message || '禁用失败');
        disableForm.code = '';
      }
    } catch (error) {
      message.error('禁用失败');
      disableForm.code = '';
    } finally {
      loading.value = false;
    }
  }

  function copySecret() {
    copyText(setupData.secretKey);
  }

  function copyText(text: string) {
    navigator.clipboard.writeText(text).then(() => {
      message.success('已复制到剪贴板');
    }).catch(() => {
      message.error('复制失败');
    });
  }

  function resetSetup() {
    setupStep.value = 1;
    verifyForm.code = '';
    setupData.secretKey = '';
    setupData.qrCode = '';
  }
</script>

<style lang="less" scoped>
  .two-factor-auth {
    .qr-code {
      max-width: 200px;
      height: auto;
      border: 1px solid #e0e0e6;
      border-radius: 6px;
    }
  }
</style>