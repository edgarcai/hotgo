import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import {
  getTwoFactorStatus,
  enable2FA,
  verifySetup2FA,
  disable2FA,
  regenerateBackupCodes,
  get2FAQRCode,
  type TwoFactorStatusResponse,
  type Enable2FAParams,
  type Enable2FAResponse,
  type VerifySetup2FAParams,
  type Disable2FAParams,
} from '@/api/auth/twoFactor';

/**
 * 2FA设置状态接口
 */
export interface TwoFactorState {
  // 2FA状态信息
  status: TwoFactorStatusResponse | null;
  // 启用流程状态
  enableFlow: {
    step: number;
    qrCodeUrl: string;
    secretKey: string;
    backupCodes: string[];
  };
  // 加载状态
  loading: {
    status: boolean;
    enable: boolean;
    disable: boolean;
    regenerate: boolean;
    verify: boolean;
  };
  // 错误信息
  errors: {
    general: string;
    verification: string;
    password: string;
  };
}

/**
 * 2FA状态管理Store
 */
export const useTwoFactorStore = defineStore('twoFactor', () => {
  // 状态定义
  const status = ref<TwoFactorStatusResponse | null>(null);
  const enableFlow = ref({
    step: 1,
    qrCodeUrl: '',
    secretKey: '',
    backupCodes: [] as string[],
  });

  const loading = ref({
    status: false,
    enable: false,
    disable: false,
    regenerate: false,
    verify: false,
    verifyingPassword: false,
    verifyingSetup: false,
  });

  const errors = ref({
    general: '',
    verification: '',
    password: '',
  });

  // 计算属性
  const isEnabled = computed(() => status.value?.enabled || false);
  const backupCodesCount = computed(() => status.value?.backupCodesCount || 0);
  const isLowBackupCodes = computed(() => backupCodesCount.value <= 2);
  const hasNoBackupCodes = computed(() => backupCodesCount.value === 0);
  const enabledAt = computed(() => status.value?.enabledAt);
  const lastUsedAt = computed(() => status.value?.lastUsedAt);

  /**
   * 清除错误信息
   */
  const clearErrors = () => {
    errors.value = {
      general: '',
      verification: '',
      password: '',
    };
  };

  /**
   * 设置错误信息
   */
  const setError = (type: keyof typeof errors.value, message: string) => {
    errors.value[type] = message;
  };

  /**
   * 重置启用流程状态
   */
  const resetEnableFlow = () => {
    enableFlow.value = {
      step: 1,
      qrCodeUrl: '',
      secretKey: '',
      backupCodes: [],
    };
  };

  /**
   * 获取2FA状态
   */
  const fetchStatus = async (): Promise<boolean> => {
    try {
      loading.value.status = true;
      clearErrors();

      const response = await getTwoFactorStatus();
      status.value = response.data;
      return true;
    } catch (error: any) {
      console.error('Failed to fetch 2FA status:', error);
      setError('general', error.message || '获取2FA状态失败');
      return false;
    } finally {
      loading.value.status = false;
    }
  };

  /**
   * 启用2FA - 第一步：验证密码并获取QR码
   */
  const startEnable2FA = async (password: string): Promise<boolean> => {
    try {
      loading.value.enable = true;
      clearErrors();

      const params: Enable2FAParams = { password };
      const response = await enable2FA(params);
      
      // 调试：打印响应数据
      console.log('Enable2FA full response:', response);
      
      // 检查响应数据是否存在
      if (!response) {
        throw new Error('服务器响应数据为空');
      }
      
      // 根据axios transformRequestData的处理，response应该直接是Enable2FAResponse类型
      // 但如果transformRequestData没有处理，则需要访问response.data
      const data = (response as any).data || response;
      console.log('Enable2FA processed data:', data);
      
      // 检查必要字段是否存在
      if (!data.qrCodeUrl || !data.secret) {
        throw new Error('服务器返回的数据不完整');
      }

      // 更新启用流程状态
      enableFlow.value = {
        step: 2,
        qrCodeUrl: data.qrCodeUrl,
        secretKey: data.secret,
        backupCodes: data.backupCodes || [],
      };

      return true;
    } catch (error: any) {
      console.error('Failed to start 2FA enable:', error);
      if (error.message?.includes('password')) {
        setError('password', '密码验证失败，请检查密码是否正确');
      } else {
        setError('general', error.message || '启用2FA失败');
      }
      return false;
    } finally {
      loading.value.enable = false;
    }
  };

  /**
   * 验证2FA设置
   */
  const verifyEnable2FA = async (code: string): Promise<boolean> => {
    try {
      loading.value.verify = true;
      clearErrors();

      const params: VerifySetup2FAParams = {
        code,
      };

      await verifySetup2FA(params);

      // 验证成功，更新状态
      enableFlow.value.step = 4;
      await fetchStatus(); // 重新获取状态

      return true;
    } catch (error: any) {
      console.error('Failed to verify 2FA setup:', error);
      if (error.message?.includes('code') || error.message?.includes('验证码')) {
        setError('verification', '验证码不正确，请检查后重试');
      } else {
        setError('general', error.message || '验证2FA设置失败');
      }
      return false;
    } finally {
      loading.value.verify = false;
    }
  };

  /**
   * 禁用2FA
   */
  const disable2FAAuth = async (password: string, code: string): Promise<boolean> => {
    try {
      loading.value.disable = true;
      clearErrors();

      const params: Disable2FAParams = { password, code };
      await disable2FA(params);

      // 禁用成功，重置状态
      status.value = null;
      resetEnableFlow();
      await fetchStatus();

      return true;
    } catch (error: any) {
      console.error('Failed to disable 2FA:', error);
      if (error.message?.includes('password')) {
        setError('password', '密码验证失败');
      } else if (error.message?.includes('code') || error.message?.includes('验证码')) {
        setError('verification', '验证码不正确');
      } else {
        setError('general', error.message || '禁用2FA失败');
      }
      return false;
    } finally {
      loading.value.disable = false;
    }
  };

  /**
   * 重新生成备用码
   */
  const regenerateBackupCodesAuth = async (password: string): Promise<string[] | null> => {
    try {
      loading.value.regenerate = true;
      clearErrors();

      const response = await regenerateBackupCodes(password);
      const newBackupCodes = response.data.backupCodes;

      // 更新状态
      await fetchStatus();

      return newBackupCodes;
    } catch (error: any) {
      console.error('Failed to regenerate backup codes:', error);
      if (error.message?.includes('password')) {
        setError('password', '密码验证失败');
      } else {
        setError('general', error.message || '重新生成备用码失败');
      }
      return null;
    } finally {
      loading.value.regenerate = false;
    }
  };

  /**
   * 获取QR码（用于重新显示）
   */
  const fetchQRCode = async (): Promise<string | null> => {
    try {
      const response = await get2FAQRCode();
      return response.data.qrCodeUrl;
    } catch (error: any) {
      console.error('Failed to fetch QR code:', error);
      setError('general', error.message || '获取二维码失败');
      return null;
    }
  };

  /**
   * 设置启用流程步骤
   */
  const setEnableStep = (step: number) => {
    enableFlow.value.step = step;
  };

  /**
   * 完成启用流程
   */
  const completeEnableFlow = () => {
    resetEnableFlow();
    fetchStatus();
  };

  /**
   * 取消启用流程
   */
  const cancelEnableFlow = () => {
    resetEnableFlow();
    clearErrors();
  };

  // 返回store接口
  return {
    // 状态
    status,
    enableFlow,
    loading,
    errors,

    // 计算属性
    isEnabled,
    backupCodesCount,
    isLowBackupCodes,
    hasNoBackupCodes,
    enabledAt,
    lastUsedAt,

    // 方法
    clearErrors,
    setError,
    resetEnableFlow,
    fetchStatus,
    startEnable2FA,
    verifyEnable2FA,
    disable2FAAuth,
    regenerateBackupCodesAuth,
    fetchQRCode,
    setEnableStep,
    completeEnableFlow,
    cancelEnableFlow,
  };
});

/**
 * 2FA Store类型导出
 */
export type TwoFactorStore = ReturnType<typeof useTwoFactorStore>;
