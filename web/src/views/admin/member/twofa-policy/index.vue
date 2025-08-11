<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="2FA策略管理">
        配置双因素认证的全局策略，包括强制启用规则、安全设置等
      </n-card>
    </div>

    <n-grid :cols="1" :y-gap="16">
      <!-- 基础策略配置 -->
      <n-grid-item>
        <n-card :bordered="false" title="基础策略">
          <n-form
            ref="basicFormRef"
            :model="basicPolicy"
            :rules="basicRules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="全局启用2FA" path="globalEnabled">
                  <n-switch
                    v-model:value="basicPolicy.globalEnabled"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    启用后所有用户必须配置2FA才能登录
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="新用户强制启用" path="forceNewUsers">
                  <n-switch
                    v-model:value="basicPolicy.forceNewUsers"
                    :round="false"
                  >
                    <template #checked>强制</template>
                    <template #unchecked>可选</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    新注册用户必须立即设置2FA
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="管理员强制启用" path="forceAdmins">
                  <n-switch
                    v-model:value="basicPolicy.forceAdmins"
                    :round="false"
                  >
                    <template #checked>强制</template>
                    <template #unchecked>可选</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    管理员用户必须启用2FA
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="宽限期（天）" path="gracePeriod">
                  <n-input-number
                    v-model:value="basicPolicy.gracePeriod"
                    :min="0"
                    :max="30"
                    placeholder="设置宽限期天数"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    用户可延迟启用2FA的天数
                  </n-text>
                </n-form-item>
              </n-grid-item>
            </n-grid>
          </n-form>
        </n-card>
      </n-grid-item>

      <!-- 安全策略配置 -->
      <n-grid-item>
        <n-card :bordered="false" title="安全策略">
          <n-form
            ref="securityFormRef"
            :model="securityPolicy"
            :rules="securityRules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="验证码有效期（秒）" path="codeValidityPeriod">
                  <n-input-number
                    v-model:value="securityPolicy.codeValidityPeriod"
                    :min="30"
                    :max="300"
                    placeholder="验证码有效期"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    TOTP验证码的有效时间窗口
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="最大失败次数" path="maxFailureAttempts">
                  <n-input-number
                    v-model:value="securityPolicy.maxFailureAttempts"
                    :min="3"
                    :max="10"
                    placeholder="最大失败次数"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    连续验证失败后锁定账户
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="锁定时间（分钟）" path="lockoutDuration">
                  <n-input-number
                    v-model:value="securityPolicy.lockoutDuration"
                    :min="5"
                    :max="60"
                    placeholder="锁定时间"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    账户锁定的持续时间
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="备用码数量" path="backupCodesCount">
                  <n-input-number
                    v-model:value="securityPolicy.backupCodesCount"
                    :min="5"
                    :max="20"
                    placeholder="备用码数量"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    生成的备用码数量
                  </n-text>
                </n-form-item>
              </n-grid-item>
            </n-grid>
          </n-form>
        </n-card>
      </n-grid-item>

      <!-- 通知策略配置 -->
      <n-grid-item>
        <n-card :bordered="false" title="通知策略">
          <n-form
            ref="notificationFormRef"
            :model="notificationPolicy"
            :rules="notificationRules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="启用邮件通知" path="emailNotification">
                  <n-switch
                    v-model:value="notificationPolicy.emailNotification"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    2FA状态变更时发送邮件通知
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="启用短信通知" path="smsNotification">
                  <n-switch
                    v-model:value="notificationPolicy.smsNotification"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    2FA状态变更时发送短信通知
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="通知管理员" path="notifyAdmins">
                  <n-switch
                    v-model:value="notificationPolicy.notifyAdmins"
                    :round="false"
                  >
                    <template #checked>通知</template>
                    <template #unchecked>不通知</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    用户2FA操作时通知管理员
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="异常登录通知" path="suspiciousActivityNotification">
                  <n-switch
                    v-model:value="notificationPolicy.suspiciousActivityNotification"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    检测到异常登录时发送通知
                  </n-text>
                </n-form-item>
              </n-grid-item>
            </n-grid>
          </n-form>
        </n-card>
      </n-grid-item>

      <!-- 高级策略配置 -->
      <n-grid-item>
        <n-card :bordered="false" title="高级策略">
          <n-form
            ref="advancedFormRef"
            :model="advancedPolicy"
            :rules="advancedRules"
            label-placement="left"
            label-width="auto"
            require-mark-placement="right-hanging"
          >
            <n-grid :cols="2" :x-gap="24">
              <n-grid-item>
                <n-form-item label="IP白名单检查" path="ipWhitelistEnabled">
                  <n-switch
                    v-model:value="advancedPolicy.ipWhitelistEnabled"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    仅允许白名单IP绕过2FA
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="设备记忆" path="deviceRememberEnabled">
                  <n-switch
                    v-model:value="advancedPolicy.deviceRememberEnabled"
                    :round="false"
                  >
                    <template #checked>启用</template>
                    <template #unchecked>禁用</template>
                  </n-switch>
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    记住受信任设备，减少2FA频率
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="记忆期限（天）" path="deviceRememberDays">
                  <n-input-number
                    v-model:value="advancedPolicy.deviceRememberDays"
                    :min="1"
                    :max="90"
                    :disabled="!advancedPolicy.deviceRememberEnabled"
                    placeholder="设备记忆天数"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    受信任设备的记忆时长
                  </n-text>
                </n-form-item>
              </n-grid-item>
              <n-grid-item>
                <n-form-item label="强制重新验证间隔（小时）" path="forceReauthInterval">
                  <n-input-number
                    v-model:value="advancedPolicy.forceReauthInterval"
                    :min="1"
                    :max="72"
                    placeholder="重新验证间隔"
                  />
                  <n-text depth="3" style="margin-left: 12px; font-size: 12px">
                    强制用户重新进行2FA验证的间隔
                  </n-text>
                </n-form-item>
              </n-grid-item>
            </n-grid>
          </n-form>
        </n-card>
      </n-grid-item>

      <!-- 操作按钮 -->
      <n-grid-item>
        <n-card :bordered="false">
          <n-space justify="center">
            <n-button type="primary" size="large" @click="saveAllPolicies" :loading="saveLoading">
              <template #icon>
                <n-icon>
                  <SaveOutlined />
                </n-icon>
              </template>
              保存所有策略
            </n-button>
            <n-button size="large" @click="resetPolicies">
              <template #icon>
                <n-icon>
                  <ReloadOutlined />
                </n-icon>
              </template>
              重置为默认值
            </n-button>
            <n-button type="info" size="large" @click="previewPolicies">
              <template #icon>
                <n-icon>
                  <EyeOutlined />
                </n-icon>
              </template>
              预览配置
            </n-button>
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 预览配置弹窗 -->
    <n-modal v-model:show="showPreviewModal" preset="card" title="策略配置预览" style="width: 800px">
      <n-code :code="previewCode" language="json" show-line-numbers />
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, computed } from 'vue';
  import { useMessage, useDialog } from 'naive-ui';
  import { SaveOutlined, ReloadOutlined, EyeOutlined } from '@vicons/antd';

  const message = useMessage();
  const dialog = useDialog();

  // 表单引用
  const basicFormRef = ref();
  const securityFormRef = ref();
  const notificationFormRef = ref();
  const advancedFormRef = ref();

  // 状态管理
  const saveLoading = ref(false);
  const showPreviewModal = ref(false);

  // 基础策略配置
  const basicPolicy = reactive({
    globalEnabled: false,
    forceNewUsers: false,
    forceAdmins: true,
    gracePeriod: 7,
  });

  // 安全策略配置
  const securityPolicy = reactive({
    codeValidityPeriod: 30,
    maxFailureAttempts: 5,
    lockoutDuration: 15,
    backupCodesCount: 10,
  });

  // 通知策略配置
  const notificationPolicy = reactive({
    emailNotification: true,
    smsNotification: false,
    notifyAdmins: true,
    suspiciousActivityNotification: true,
  });

  // 高级策略配置
  const advancedPolicy = reactive({
    ipWhitelistEnabled: false,
    deviceRememberEnabled: true,
    deviceRememberDays: 30,
    forceReauthInterval: 24,
  });

  // 表单验证规则
  const basicRules = {
    gracePeriod: {
      type: 'number',
      required: true,
      message: '请设置宽限期',
      trigger: ['blur', 'change'],
    },
  };

  const securityRules = {
    codeValidityPeriod: {
      type: 'number',
      required: true,
      message: '请设置验证码有效期',
      trigger: ['blur', 'change'],
    },
    maxFailureAttempts: {
      type: 'number',
      required: true,
      message: '请设置最大失败次数',
      trigger: ['blur', 'change'],
    },
    lockoutDuration: {
      type: 'number',
      required: true,
      message: '请设置锁定时间',
      trigger: ['blur', 'change'],
    },
    backupCodesCount: {
      type: 'number',
      required: true,
      message: '请设置备用码数量',
      trigger: ['blur', 'change'],
    },
  };

  const notificationRules = {};

  const advancedRules = {
    deviceRememberDays: {
      type: 'number',
      required: true,
      message: '请设置设备记忆天数',
      trigger: ['blur', 'change'],
    },
    forceReauthInterval: {
      type: 'number',
      required: true,
      message: '请设置重新验证间隔',
      trigger: ['blur', 'change'],
    },
  };

  // 预览配置代码
  const previewCode = computed(() => {
    const allPolicies = {
      basic: basicPolicy,
      security: securityPolicy,
      notification: notificationPolicy,
      advanced: advancedPolicy,
    };
    return JSON.stringify(allPolicies, null, 2);
  });

  // 保存所有策略
  const saveAllPolicies = async () => {
    try {
      // 验证所有表单
      await Promise.all([
        basicFormRef.value?.validate(),
        securityFormRef.value?.validate(),
        notificationFormRef.value?.validate(),
        advancedFormRef.value?.validate(),
      ]);

      saveLoading.value = true;

      // 这里应该调用API保存策略配置
      // await saveTwoFactorPolicies({
      //   basic: basicPolicy,
      //   security: securityPolicy,
      //   notification: notificationPolicy,
      //   advanced: advancedPolicy,
      // });

      message.success('策略配置保存成功');
    } catch (error) {
      message.error('策略配置保存失败');
    } finally {
      saveLoading.value = false;
    }
  };

  // 重置策略为默认值
  const resetPolicies = () => {
    dialog.warning({
      title: '重置策略',
      content: '确定要重置所有策略为默认值吗？此操作不可撤销。',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        // 重置基础策略
        Object.assign(basicPolicy, {
          globalEnabled: false,
          forceNewUsers: false,
          forceAdmins: true,
          gracePeriod: 7,
        });

        // 重置安全策略
        Object.assign(securityPolicy, {
          codeValidityPeriod: 30,
          maxFailureAttempts: 5,
          lockoutDuration: 15,
          backupCodesCount: 10,
        });

        // 重置通知策略
        Object.assign(notificationPolicy, {
          emailNotification: true,
          smsNotification: false,
          notifyAdmins: true,
          suspiciousActivityNotification: true,
        });

        // 重置高级策略
        Object.assign(advancedPolicy, {
          ipWhitelistEnabled: false,
          deviceRememberEnabled: true,
          deviceRememberDays: 30,
          forceReauthInterval: 24,
        });

        message.success('策略已重置为默认值');
      },
    });
  };

  // 预览策略配置
  const previewPolicies = () => {
    showPreviewModal.value = true;
  };

  // 组件挂载时加载现有配置
  const loadPolicies = async () => {
    try {
      // 这里应该调用API加载现有策略配置
      // const response = await getTwoFactorPolicies();
      // Object.assign(basicPolicy, response.data.basic);
      // Object.assign(securityPolicy, response.data.security);
      // Object.assign(notificationPolicy, response.data.notification);
      // Object.assign(advancedPolicy, response.data.advanced);
    } catch (error) {
      message.error('加载策略配置失败');
    }
  };

  // 初始化
  loadPolicies();
</script>

<style scoped>
.n-form-item {
  margin-bottom: 24px;
}

.n-text {
  display: block;
  margin-top: 4px;
}
</style>