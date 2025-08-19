import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { anixUserApi, UserInfo, LoginParams, RegisterParams } from '@/api/anix/user';
import { storage } from '@/utils/Storage';
import { ACCESS_TOKEN, CURRENT_USER } from '@/store/mutation-types';
import router from '@/router';
import { PageEnum } from '@/enums/pageEnum';

/**
 * AniX用户状态管理
 */
export const useAnixUserStore = defineStore('anixUser', () => {
  // 状态
  const token = ref<string>(storage.get(ACCESS_TOKEN, '') || '');
  const userInfo = ref<UserInfo | null>(storage.get(CURRENT_USER, null));
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value);

  // 设置Token
  const setToken = (newToken: string) => {
    token.value = newToken;
    storage.set(ACCESS_TOKEN, newToken);
  };

  // 设置用户信息
  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info;
    storage.set(CURRENT_USER, info);
  };

  // 清除用户信息
  const clearUserInfo = () => {
    token.value = '';
    userInfo.value = null;
    storage.remove(ACCESS_TOKEN);
    storage.remove(CURRENT_USER);
  };

  // 用户注册
  const register = async (params: RegisterParams) => {
    try {
      const result = await anixUserApi.register(params);
      return result;
    } catch (error) {
      throw error;
    }
  };

  // 用户登录
  const login = async (params: LoginParams) => {
    try {
      const result = await anixUserApi.login(params);

      // 保存token和用户信息
      setToken(result.token);
      setUserInfo(result.user);

      return result;
    } catch (error) {
      throw error;
    }
  };

  // 获取用户信息
  const getUserInfo = async () => {
    try {
      const result = await anixUserApi.getProfile();
      setUserInfo(result.user);
      return result.user;
    } catch (error) {
      // 如果获取用户信息失败，清除本地存储
      clearUserInfo();
      throw error;
    }
  };

  // 更新用户信息
  const updateProfile = async (params: any) => {
    try {
      const result = await anixUserApi.updateProfile(params);
      // 更新成功后重新获取用户信息
      await getUserInfo();
      return result;
    } catch (error) {
      throw error;
    }
  };

  // 修改密码
  const changePassword = async (params: any) => {
    try {
      const result = await anixUserApi.changePassword(params);
      return result;
    } catch (error) {
      throw error;
    }
  };

  // 用户登出
  const logout = async () => {
    try {
      await anixUserApi.logout();
    } catch (error) {
      console.error('登出请求失败:', error);
    } finally {
      // 无论请求是否成功，都清除本地存储
      clearUserInfo();
      // 跳转到登录页
      router.push(PageEnum.BASE_LOGIN);
    }
  };

  // 初始化用户信息（应用启动时调用）
  const initUserInfo = async () => {
    if (token.value && !userInfo.value) {
      try {
        await getUserInfo();
      } catch (error) {
        console.error('初始化用户信息失败:', error);
        clearUserInfo();
      }
    }
  };

  return {
    // 状态
    token: computed(() => token.value),
    userInfo: computed(() => userInfo.value),
    isLoggedIn,

    // 方法
    setToken,
    setUserInfo,
    clearUserInfo,
    register,
    login,
    getUserInfo,
    updateProfile,
    changePassword,
    logout,
    initUserInfo,
  };
});
