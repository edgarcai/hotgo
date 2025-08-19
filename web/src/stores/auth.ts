import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { UserInfo } from '@/api/anix';

// Token 存储键名
const TOKEN_KEY = 'anix_token';
const USER_KEY = 'anix_user';
const REFRESH_TOKEN_KEY = 'anix_refresh_token';

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string>('');
  const refreshToken = ref<string>('');
  const user = ref<UserInfo | null>(null);
  const isLoading = ref(false);

  // 计算属性
  const isAuthenticated = computed(() => {
    return !!token.value && !!user.value;
  });

  const userAvatar = computed(() => {
    return user.value?.avatar || '/default-avatar.png';
  });

  const userName = computed(() => {
    return user.value?.nickname || user.value?.username || '未知用户';
  });

  // 初始化认证状态（从本地存储恢复）
  const initAuth = () => {
    try {
      const savedToken = localStorage.getItem(TOKEN_KEY);
      const savedRefreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
      const savedUser = localStorage.getItem(USER_KEY);

      if (savedToken) {
        token.value = savedToken;
      }

      if (savedRefreshToken) {
        refreshToken.value = savedRefreshToken;
      }

      if (savedUser) {
        user.value = JSON.parse(savedUser);
      }
    } catch (error) {
      console.error('初始化认证状态失败:', error);
      clearAuth();
    }
  };

  // 设置 Token
  const setToken = (newToken: string, newRefreshToken?: string) => {
    token.value = newToken;
    localStorage.setItem(TOKEN_KEY, newToken);

    if (newRefreshToken) {
      refreshToken.value = newRefreshToken;
      localStorage.setItem(REFRESH_TOKEN_KEY, newRefreshToken);
    }
  };

  // 设置用户信息
  const setUser = (userInfo: UserInfo) => {
    user.value = userInfo;
    localStorage.setItem(USER_KEY, JSON.stringify(userInfo));
  };

  // 更新用户信息
  const updateUser = (updates: Partial<UserInfo>) => {
    if (user.value) {
      user.value = { ...user.value, ...updates };
      localStorage.setItem(USER_KEY, JSON.stringify(user.value));
    }
  };

  // 清除认证状态
  const clearAuth = () => {
    token.value = '';
    refreshToken.value = '';
    user.value = null;
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    localStorage.removeItem('anix_remember');
  };

  // 登录
  const login = async (loginData: { token: string; user: UserInfo; refreshToken?: string }) => {
    isLoading.value = true;
    try {
      setToken(loginData.token, loginData.refreshToken);
      setUser(loginData.user);
      return true;
    } catch (error) {
      console.error('登录失败:', error);
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  // 登出
  const logout = async () => {
    isLoading.value = true;
    try {
      // 这里可以调用登出 API
      // await anixApi.logout()
      clearAuth();
      return true;
    } catch (error) {
      console.error('登出失败:', error);
      // 即使 API 调用失败，也要清除本地状态
      clearAuth();
      return false;
    } finally {
      isLoading.value = false;
    }
  };

  // 检查 Token 是否有效
  const checkTokenValidity = () => {
    if (!token.value) {
      return false;
    }

    try {
      // 解析 JWT Token 检查过期时间
      const payload = JSON.parse(atob(token.value.split('.')[1]));
      const currentTime = Math.floor(Date.now() / 1000);

      if (payload.exp && payload.exp < currentTime) {
        // Token 已过期
        clearAuth();
        return false;
      }

      return true;
    } catch (error) {
      console.error('Token 解析失败:', error);
      clearAuth();
      return false;
    }
  };

  // 刷新 Token
  const refreshAuthToken = async () => {
    if (!refreshToken.value) {
      clearAuth();
      return false;
    }

    try {
      // 这里应该调用刷新 Token 的 API
      // const response = await anixApi.refreshToken(refreshToken.value)
      // setToken(response.token, response.refreshToken)
      // return true

      // 暂时返回 false，等待后端实现
      return false;
    } catch (error) {
      console.error('刷新 Token 失败:', error);
      clearAuth();
      return false;
    }
  };

  // 获取认证头
  const getAuthHeader = () => {
    return token.value ? `Bearer ${token.value}` : '';
  };

  // 检查用户权限
  const hasPermission = (permission: string) => {
    // 这里可以根据用户角色或权限列表进行检查
    // 暂时返回 true，等待权限系统实现
    return isAuthenticated.value;
  };

  // 检查用户角色
  const hasRole = (role: string) => {
    // 这里可以根据用户角色进行检查
    // 暂时返回 true，等待角色系统实现
    return isAuthenticated.value;
  };

  // 获取用户设置
  const getUserSetting = (key: string, defaultValue: any = null) => {
    try {
      const settings = localStorage.getItem(`anix_user_settings_${user.value?.id}`);
      if (settings) {
        const parsed = JSON.parse(settings);
        return parsed[key] !== undefined ? parsed[key] : defaultValue;
      }
    } catch (error) {
      console.error('获取用户设置失败:', error);
    }
    return defaultValue;
  };

  // 设置用户设置
  const setUserSetting = (key: string, value: any) => {
    if (!user.value) return;

    try {
      const settingsKey = `anix_user_settings_${user.value.id}`;
      const settings = localStorage.getItem(settingsKey);
      const parsed = settings ? JSON.parse(settings) : {};
      parsed[key] = value;
      localStorage.setItem(settingsKey, JSON.stringify(parsed));
    } catch (error) {
      console.error('设置用户设置失败:', error);
    }
  };

  // 清除用户设置
  const clearUserSettings = () => {
    if (!user.value) return;

    try {
      localStorage.removeItem(`anix_user_settings_${user.value.id}`);
    } catch (error) {
      console.error('清除用户设置失败:', error);
    }
  };

  return {
    // 状态
    token: readonly(token),
    refreshToken: readonly(refreshToken),
    user: readonly(user),
    isLoading: readonly(isLoading),

    // 计算属性
    isAuthenticated,
    userAvatar,
    userName,

    // 方法
    initAuth,
    setToken,
    setUser,
    updateUser,
    clearAuth,
    login,
    logout,
    checkTokenValidity,
    refreshAuthToken,
    getAuthHeader,
    hasPermission,
    hasRole,
    getUserSetting,
    setUserSetting,
    clearUserSettings,
  };
});

// 导出类型
export type AuthStore = ReturnType<typeof useAuthStore>;
