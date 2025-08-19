import { http } from '@/utils/http/axios';

// 用户信息接口
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  email: string;
}

// 用户注册接口
export interface RegisterParams {
  username: string;
  password: string;
  email: string;
  nickname?: string;
}

export interface RegisterResult {
  message: string;
}

// 用户登录接口
export interface LoginParams {
  username: string;
  password: string;
}

export interface LoginResult {
  token: string;
  user: UserInfo;
}

// 更新用户信息接口
export interface UpdateProfileParams {
  nickname?: string;
  avatar?: string;
  email?: string;
}

export interface UpdateProfileResult {
  message: string;
}

// 修改密码接口
export interface ChangePasswordParams {
  oldPassword: string;
  newPassword: string;
}

export interface ChangePasswordResult {
  message: string;
}

// 登出接口
export interface LogoutResult {
  message: string;
}

/**
 * AniX用户相关API
 */
export const anixUserApi = {
  // 用户注册
  register: (params: RegisterParams) => {
    return http.request<RegisterResult>({
      url: '/anix/user/register',
      method: 'POST',
      data: params,
    });
  },

  // 用户登录
  login: (params: LoginParams) => {
    return http.request<LoginResult>({
      url: '/anix/user/login',
      method: 'POST',
      data: params,
    });
  },

  // 获取用户信息
  getProfile: () => {
    return http.request<{ user: UserInfo }>({
      url: '/anix/user/info',
      method: 'GET',
    });
  },

  // 更新用户信息
  updateProfile: (params: UpdateProfileParams) => {
    return http.request<UpdateProfileResult>({
      url: '/anix/user/profile',
      method: 'PUT',
      data: params,
    });
  },

  // 修改密码
  changePassword: (params: ChangePasswordParams) => {
    return http.request<ChangePasswordResult>({
      url: '/anix/user/password',
      method: 'PUT',
      data: params,
    });
  },

  // 用户登出
  logout: () => {
    return http.request<LogoutResult>({
      url: '/anix/user/logout',
      method: 'POST',
    });
  },
};
