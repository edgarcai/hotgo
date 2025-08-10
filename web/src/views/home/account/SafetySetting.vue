<template>
  <n-grid cols="1" responsive="screen" class="-mt-5">
    <n-grid-item>
      <n-list>
        <n-list-item>
          <template #suffix>
            <n-button type="primary" text @click="openUpdatePassForm">修改</n-button>
          </template>
          <n-thing title="账户密码">
            <template #description
              ><span class="text-gray-400">绑定手机和邮箱，并设置密码，帐号更安全</span></template
            >
          </n-thing>
        </n-list-item>
        <n-list-item>
          <template #suffix>
            <n-button type="primary" text @click="openUpdateMobileForm">修改</n-button>
          </template>
          <n-thing title="绑定手机">
            <template #description
              ><span class="text-gray-400"
                >已绑定手机号：+86{{ userStore.info?.mobile }}</span
              ></template
            >
          </n-thing>
        </n-list-item>
        <n-list-item>
          <template #suffix>
            <n-button type="primary" text @click="openUpdateEmailForm">修改</n-button>
          </template>
          <n-thing title="绑定邮箱">
            <template #description
              ><span class="text-gray-400">已绑定邮箱：{{ userStore.info?.email }}</span></template
            >
          </n-thing>
        </n-list-item>
        <n-list-item>
          <template #suffix>
            <n-button type="primary" text @click="open2FAForm">{{ twoFAStatus ? '管理' : '启用' }}</n-button>
          </template>
          <n-thing title="双因子认证">
            <template #description>
              <span class="text-gray-400">
                {{ twoFAStatus ? '已启用双因子认证，账户更安全' : '启用双因子认证，提升账户安全性' }}
              </span>
            </template>
          </n-thing>
        </n-list-item>
      </n-list>
    </n-grid-item>
  </n-grid>

  <n-modal
    v-model:show="showModal"
    :show-icon="false"
    preset="dialog"
    title="修改登录密码"
    :style="{
      width: dialogWidth,
    }"
  >
    <n-form :label-width="80" :model="formValue" :rules="rules" ref="formRef">
      <n-form-item label="当前密码" path="oldPassword">
        <n-input
          type="password"
          v-model:value="formValue.oldPassword"
          placeholder="请输入当前密码"
        />
      </n-form-item>

      <n-form-item label="新密码" path="newPassword">
        <n-input type="password" v-model:value="formValue.newPassword" placeholder="请输入新密码" />
      </n-form-item>

      <div>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" @click="formSubmit">修改并重新登录</n-button>
        </n-space>
      </div>
    </n-form>
  </n-modal>

  <n-modal
    :block-scroll="false"
    :mask-closable="false"
    v-model:show="showMobileModal"
    :show-icon="false"
    preset="dialog"
    title="修改手机号"
    :style="{
      width: dialogWidth,
    }"
  >
    <n-form :label-width="80" :model="formMobileValue" ref="formMobileRef">
      <n-form-item label="短信验证码" path="code" v-if="userStore.info?.mobile !== ''">
        <n-input-group>
          <n-input v-model:value="formMobileValue.code" placeholder="请输入验证码" />
          <n-button
            type="primary"
            ghost
            @click="sendMobileCode"
            :disabled="isCounting"
            :loading="sendLoading"
          >
            {{ sendLabel }}
          </n-button>
        </n-input-group>

        <template #feedback> 接收号码：+86{{ userStore.info?.mobile }} </template>
      </n-form-item>

      <n-form-item label="换绑手机号" path="mobile">
        <n-input v-model:value="formMobileValue.mobile" placeholder="请输入换绑手机号" />
      </n-form-item>
      <div>
        <n-space justify="end">
          <n-button @click="showMobileModal = false">取消</n-button>
          <n-button type="primary" :loading="formMobileBtnLoading" @click="formMobileSubmit"
            >保存更新</n-button
          >
        </n-space>
      </div>
    </n-form>
  </n-modal>

  <n-modal
    :block-scroll="false"
    :mask-closable="false"
    v-model:show="showEmailModal"
    :show-icon="false"
    preset="dialog"
    title="修改邮箱"
    :style="{
      width: dialogWidth,
    }"
  >
    <n-form :label-width="80" :model="formEmailValue" ref="formEmailRef">
      <n-form-item label="邮箱验证码" path="code" v-if="userStore.info?.email !== ''">
        <n-input-group>
          <n-input v-model:value="formEmailValue.code" placeholder="请输入验证码" />
          <n-button
            type="primary"
            ghost
            @click="sendEmailCode"
            :disabled="isCounting"
            :loading="sendLoading"
          >
            {{ sendLabel }}
          </n-button>
        </n-input-group>
        <template #feedback> 接收邮箱：{{ userStore.info?.email }} </template>
      </n-form-item>

      <n-form-item label="换绑邮箱" path="email">
        <n-input v-model:value="formEmailValue.email" placeholder="请输入换绑邮箱" />
      </n-form-item>
      <div>
        <n-space justify="end">
          <n-button @click="showEmailModal = false">取消</n-button>
          <n-button type="primary" :loading="formEmailBtnLoading" @click="formEmailSubmit"
            >保存更新</n-button
          >
        </n-space>
      </div>
    </n-form>
  </n-modal>

  <!-- 2FA设置模态框 -->
  <n-modal
    :block-scroll="false"
    :mask-closable="false"
    v-model:show="show2FAModal"
    :show-icon="false"
    preset="dialog"
    :title="twoFAStatus ? '管理双因子认证' : '启用双因子认证'"
    :style="{
      width: dialogWidth,
    }"
  >
    <div v-if="!twoFAStatus">
      <!-- 启用2FA流程 -->
      <n-steps :current="currentStep" :status="currentStatus">
        <n-step title="扫描二维码" />
        <n-step title="验证设备" />
        <n-step title="保存备用码" />
      </n-steps>
      
      <div class="mt-4">
        <!-- 步骤1: 扫描二维码 -->
        <div v-if="currentStep === 1">
          <n-alert type="info" class="mb-4">
            请使用身份验证器应用（如Google Authenticator、Microsoft Authenticator等）扫描下方二维码
          </n-alert>
          <div class="text-center mb-4">
            <div v-if="qrCodeUrl" class="inline-block p-4 bg-white border rounded">
              <img :src="qrCodeUrl" alt="2FA QR Code" class="w-48 h-48" />
            </div>
            <n-spin v-else size="large" />
          </div>
          <n-alert type="warning" class="mb-4">
            <template #header>手动输入密钥</template>
            如果无法扫描二维码，请手动输入以下密钥：<br/>
            <n-code :code="secretKey" class="mt-2" />
          </n-alert>
        </div>
        
        <!-- 步骤2: 验证设备 -->
        <div v-if="currentStep === 2">
          <n-alert type="info" class="mb-4">
            请输入身份验证器应用中显示的6位数字验证码
          </n-alert>
          <n-form :model="verify2FAForm" ref="verify2FARef">
            <n-form-item label="验证码" path="code">
              <n-input
                v-model:value="verify2FAForm.code"
                placeholder="请输入6位验证码"
                maxlength="6"
                :input-props="{ autocomplete: 'off' }"
              />
            </n-form-item>
          </n-form>
        </div>
        
        <!-- 步骤3: 保存备用码 -->
        <div v-if="currentStep === 3">
          <n-alert type="success" class="mb-4">
            双因子认证已成功启用！请妥善保存以下备用码，当无法使用身份验证器时可以使用这些代码登录。
          </n-alert>
          <n-alert type="warning" class="mb-4">
            <template #header>重要提醒</template>
            每个备用码只能使用一次，请将其保存在安全的地方。
          </n-alert>
          <div class="bg-gray-50 p-4 rounded mb-4">
            <div class="grid grid-cols-2 gap-2">
              <n-code v-for="code in backupCodes" :key="code" :code="code" class="text-center" />
            </div>
          </div>
        </div>
      </div>
      
      <div class="mt-6">
        <n-space justify="end">
          <n-button @click="show2FAModal = false">取消</n-button>
          <n-button v-if="currentStep === 1" type="primary" @click="nextStep" :disabled="!qrCodeUrl">
            下一步
          </n-button>
          <n-button v-if="currentStep === 2" type="primary" @click="verifyAndEnable" :loading="enabling2FA">
            验证并启用
          </n-button>
          <n-button v-if="currentStep === 3" type="primary" @click="finish2FASetup">
            完成设置
          </n-button>
        </n-space>
      </div>
    </div>
    
    <div v-else>
      <!-- 管理2FA -->
      <n-space vertical>
        <n-alert type="success">
          双因子认证已启用，您的账户受到额外保护。
        </n-alert>
        
        <n-card title="备用码管理" size="small">
          <template #header-extra>
            <n-button text type="primary" @click="regenerateBackupCodes" :loading="regeneratingCodes">
              重新生成
            </n-button>
          </template>
          <n-space>
            <n-button @click="showBackupCodes" ghost>
              查看备用码
            </n-button>
            <n-popconfirm @positive-click="disable2FA" negative-text="取消" positive-text="确认">
              <template #trigger>
                <n-button type="error" ghost>
                  禁用2FA
                </n-button>
              </template>
              确定要禁用双因子认证吗？这将降低您账户的安全性。
            </n-popconfirm>
          </n-space>
        </n-card>
      </n-space>
    </div>
  </n-modal>
  
  <!-- 备用码查看模态框 -->
  <n-modal
    v-model:show="showBackupCodesModal"
    preset="dialog"
    title="备用码"
    :style="{ width: dialogWidth }"
  >
    <n-alert type="warning" class="mb-4">
      请妥善保存这些备用码，每个代码只能使用一次。
    </n-alert>
    <div class="bg-gray-50 p-4 rounded">
      <div class="grid grid-cols-2 gap-2">
        <n-code v-for="code in currentBackupCodes" :key="code" :code="code" class="text-center" />
      </div>
    </div>
    <template #action>
      <n-button @click="showBackupCodesModal = false">关闭</n-button>
    </template>
  </n-modal>
</template>

<script lang="ts" setup>
  import { computed, ref, onMounted } from 'vue';
  import { useMessage } from 'naive-ui';
  import { useRouter, useRoute } from 'vue-router';
  import { useSendCode } from '@/hooks/common';
  import { adaModalWidth } from '@/utils/hotgo';
  import {
    updateMemberPwd,
    updateMemberMobile,
    updateMemberEmail,
    SendBindEmail,
    SendBindSms,
  } from '@/api/system/user';
  import { TABS_ROUTES } from '@/store/mutation-types';
  import { useUserStore } from '@/store/modules/user';
  import {
    get2FAStatus,
    setup2FA,
    verifySetup2FA,
    disable2FA as disable2FAAPI,
    getBackupCodes,
    regenerateBackupCodes as regenerateBackupCodesAPI,
  } from '@/api/system/2fa';

  const { sendLabel, isCounting, loading: sendLoading, activateSend } = useSendCode();
  const userStore = useUserStore();
  const rules = {
    basicName: {
      required: true,
      message: '请输入网站名称',
      trigger: 'blur',
    },
  };

  const formRef: any = ref(null);
  const message = useMessage();
  const router = useRouter();
  const route = useRoute();
  const showModal = ref(false);
  const formValue = ref({
    oldPassword: '',
    newPassword: '',
  });
  const dialogWidth = computed(() => {
    return adaModalWidth(580);
  });

  // 2FA相关状态
  const twoFAStatus = ref(false);
  const show2FAModal = ref(false);
  const currentStep = ref(1);
  const currentStatus = ref('process');
  const qrCodeUrl = ref('');
  const secretKey = ref('');
  const backupCodes = ref([]);
  const enabling2FA = ref(false);
  const regeneratingCodes = ref(false);
  const showBackupCodesModal = ref(false);
  const currentBackupCodes = ref([]);
  
  const verify2FAForm = ref({
    code: '',
  });
  const verify2FARef = ref(null);

  function formSubmit() {
    formRef.value.validate((errors) => {
      if (!errors) {
        updateMemberPwd({
          oldPassword: formValue.value.oldPassword,
          newPassword: formValue.value.newPassword,
        })
          .then((_res) => {
            message.success('更新成功');

            userStore.logout().then(() => {
              message.success('成功注销登录');
              // 移除标签页
              localStorage.removeItem(TABS_ROUTES);
              router
                .replace({
                  name: 'Login',
                  query: {
                    redirect: route.fullPath,
                  },
                })
                .finally(() => location.reload());
            });
          })
          .finally(() => {
            showModal.value = false;
          });
      } else {
        message.error('验证失败，请填写完整信息');
      }
    });
  }

  function openUpdatePassForm() {
    showModal.value = true;
    formValue.value.newPassword = '';
    formValue.value.oldPassword = '';
  }

  const formMobileBtnLoading = ref(false);
  const formMobileRef: any = ref(null);
  const showMobileModal = ref(false);
  const formMobileValue = ref({
    mobile: '',
    code: '',
  });

  function formMobileSubmit() {
    formMobileRef.value.validate((errors) => {
      if (!errors) {
        formMobileBtnLoading.value = true;
        updateMemberMobile({
          mobile: formMobileValue.value.mobile,
          code: formMobileValue.value.code,
        })
          .then((_res) => {
            message.success('更新成功');
            showMobileModal.value = false;
            userStore.GetInfo();
          })
          .finally(() => {
            formMobileBtnLoading.value = false;
          });
      } else {
        message.error('验证失败，请填写完整信息');
      }
    });
  }

  function openUpdateMobileForm() {
    showMobileModal.value = true;
    formMobileValue.value.mobile = '';
    formMobileValue.value.code = '';
  }

  const formEmailBtnLoading = ref(false);
  const formEmailRef: any = ref(null);
  const showEmailModal = ref(false);
  const formEmailValue = ref({
    email: '',
    code: '',
  });

  function formEmailSubmit() {
    formEmailRef.value.validate((errors) => {
      if (!errors) {
        formEmailBtnLoading.value = true;
        updateMemberEmail({
          email: formEmailValue.value.email,
          code: formEmailValue.value.code,
        })
          .then((_res) => {
            message.success('更新成功');
            showEmailModal.value = false;
            userStore.GetInfo();
          })
          .finally(() => {
            formEmailBtnLoading.value = false;
          });
      } else {
        message.error('验证失败，请填写完整信息');
      }
    });
  }

  function openUpdateEmailForm() {
    showEmailModal.value = true;
    formEmailValue.value.email = '';
    formEmailValue.value.code = '';
  }

  function sendMobileCode() {
    activateSend(SendBindSms());
  }

  function sendEmailCode() {
    activateSend(SendBindEmail());
  }

  // 2FA相关方法
  async function load2FAStatus() {
    try {
      const response = await get2FAStatus();
      // 检查是否是错误响应（当isShowErrorMessage为false时返回的错误对象）
      if (response && response.error) {
        console.log('2FA功能未启用或接口不可用');
        twoFAStatus.value = false;
        return;
      }
      twoFAStatus.value = response.data.isEnabled;
    } catch (error) {
      console.error('获取2FA状态失败:', error);
      twoFAStatus.value = false;
    }
  }

  async function open2FAForm() {
    if (twoFAStatus.value) {
      show2FAModal.value = true;
    } else {
      // 生成2FA密钥和二维码
      try {
        const response = await setup2FA();
        secretKey.value = response.data.secret;
        qrCodeUrl.value = response.data.qrCodeUrl;
        currentStep.value = 1;
        show2FAModal.value = true;
      } catch (error) {
        message.error('生成2FA密钥失败');
      }
    }
  }

  function nextStep() {
    if (currentStep.value < 3) {
      currentStep.value++;
    }
  }

  async function verifyAndEnable() {
    if (!verify2FAForm.value.code || verify2FAForm.value.code.length !== 6) {
      message.error('请输入6位验证码');
      return;
    }

    enabling2FA.value = true;
    try {
      const response = await verifySetup2FA({
        secret: secretKey.value,
        code: verify2FAForm.value.code,
      });
      
      if (response.data.success) {
        backupCodes.value = response.data.backupCodes;
        currentStep.value = 3;
        currentStatus.value = 'finish';
        message.success('双因子认证启用成功');
      } else {
        message.error('验证码错误，请重试');
      }
    } catch (error) {
      message.error('启用2FA失败');
    } finally {
      enabling2FA.value = false;
    }
  }

  function finish2FASetup() {
    twoFAStatus.value = true;
    show2FAModal.value = false;
    currentStep.value = 1;
    currentStatus.value = 'process';
    verify2FAForm.value.code = '';
    qrCodeUrl.value = '';
    secretKey.value = '';
    backupCodes.value = [];
  }

  async function disable2FA() {
    try {
      await disable2FAAPI();
      twoFAStatus.value = false;
      show2FAModal.value = false;
      message.success('双因子认证已禁用');
    } catch (error) {
      message.error('禁用2FA失败');
    }
  }

  async function showBackupCodes() {
    try {
      const response = await getBackupCodes();
      currentBackupCodes.value = response.data.codes;
      showBackupCodesModal.value = true;
    } catch (error) {
      message.error('获取备用码失败');
    }
  }

  async function regenerateBackupCodes() {
    regeneratingCodes.value = true;
    try {
      const response = await regenerateBackupCodesAPI();
      currentBackupCodes.value = response.data.codes;
      message.success('备用码已重新生成');
    } catch (error) {
      message.error('重新生成备用码失败');
    } finally {
      regeneratingCodes.value = false;
    }
  }

  onMounted(() => {
    load2FAStatus();
  });
</script>
