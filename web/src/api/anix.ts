import { http } from '@/utils/http/axios';

// 用户注册接口
export interface RegisterRequest {
  username: string;
  email: string;
  nickname: string;
  password: string;
}

export interface RegisterResponse {
  message: string;
}

// 用户登录接口
export interface LoginRequest {
  account: string; // 用户名或邮箱
  password: string;
}

export interface LoginResponse {
  token: string;
  user: UserInfo;
}

// 用户信息接口
export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  email: string;
  avatar: string;
  bio?: string;
  status: number;
  createdAt: string;
  updatedAt: string;
  lastLoginAt?: string;
}

// 更新用户信息接口
export interface UpdateProfileRequest {
  nickname: string;
  email: string;
  avatar?: string;
  bio?: string;
}

// 修改密码接口
export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}

// 头像上传响应
export interface UploadAvatarResponse {
  url: string;
}

// 漫画相关接口
export interface Comic {
  id: number;
  title: string;
  description: string;
  cover: string;
  author: string;
  status: string;
  tags: string[];
  rating: number;
  viewCount: number;
  chapterCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Chapter {
  id: number;
  comicId: number;
  title: string;
  chapterNumber: number;
  pages: string[];
  createdAt: string;
  updatedAt: string;
}

// API 类
class AnixApi {
  // 用户相关 API

  /**
   * 用户注册
   */
  async register(data: RegisterRequest): Promise<RegisterResponse> {
    return http.post('/anix/register', data);
  }

  /**
   * 用户登录
   */
  async login(data: LoginRequest): Promise<LoginResponse> {
    return http.post('/anix/login', data);
  }

  /**
   * 获取用户信息
   */
  async getProfile(): Promise<UserInfo> {
    return http.get('/anix/user/profile');
  }

  /**
   * 更新用户信息
   */
  async updateProfile(data: UpdateProfileRequest): Promise<UserInfo> {
    return http.put('/anix/user/profile', data);
  }

  /**
   * 修改密码
   */
  async changePassword(data: ChangePasswordRequest): Promise<{ message: string }> {
    return http.put('/anix/user/password', data);
  }

  /**
   * 用户登出
   */
  async logout(): Promise<{ message: string }> {
    return http.post('/anix/user/logout');
  }

  /**
   * 上传头像
   */
  async uploadAvatar(formData: FormData): Promise<UploadAvatarResponse> {
    return http.post('/anix/user/avatar', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  }

  // 漫画相关 API

  /**
   * 获取漫画列表
   */
  async getComics(params?: {
    page?: number;
    limit?: number;
    search?: string;
    category?: string;
    status?: string;
    sort?: string;
  }): Promise<{
    list: Comic[];
    total: number;
    page: number;
    limit: number;
  }> {
    return http.get('/anix/comics', { params });
  }

  /**
   * 获取漫画详情
   */
  async getComic(id: number): Promise<Comic> {
    return http.get(`/anix/comics/${id}`);
  }

  /**
   * 获取漫画章节列表
   */
  async getChapters(comicId: number): Promise<Chapter[]> {
    return http.get(`/anix/comics/${comicId}/chapters`);
  }

  /**
   * 获取章节详情
   */
  async getChapter(comicId: number, chapterId: number): Promise<Chapter> {
    return http.get(`/anix/comics/${comicId}/chapters/${chapterId}`);
  }

  /**
   * 获取热门漫画
   */
  async getPopularComics(limit = 10): Promise<Comic[]> {
    return http.get('/anix/comics/popular', { params: { limit } });
  }

  /**
   * 获取最新漫画
   */
  async getLatestComics(limit = 10): Promise<Comic[]> {
    return http.get('/anix/comics/latest', { params: { limit } });
  }

  /**
   * 搜索漫画
   */
  async searchComics(
    keyword: string,
    params?: {
      page?: number;
      limit?: number;
    }
  ): Promise<{
    list: Comic[];
    total: number;
    page: number;
    limit: number;
  }> {
    return http.get('/anix/comics/search', {
      params: {
        keyword,
        ...params,
      },
    });
  }

  // 收藏相关 API（需要登录）

  /**
   * 收藏漫画
   */
  async favoriteComic(comicId: number): Promise<{ message: string }> {
    return http.post(`/anix/user/favorites/${comicId}`);
  }

  /**
   * 取消收藏漫画
   */
  async unfavoriteComic(comicId: number): Promise<{ message: string }> {
    return http.delete(`/anix/user/favorites/${comicId}`);
  }

  /**
   * 获取收藏列表
   */
  async getFavorites(params?: { page?: number; limit?: number }): Promise<{
    list: Comic[];
    total: number;
    page: number;
    limit: number;
  }> {
    return http.get('/anix/user/favorites', { params });
  }

  /**
   * 检查是否已收藏
   */
  async checkFavorite(comicId: number): Promise<{ isFavorite: boolean }> {
    return http.get(`/anix/user/favorites/${comicId}/check`);
  }

  // 阅读历史相关 API（需要登录）

  /**
   * 记录阅读历史
   */
  async recordReadingHistory(comicId: number, chapterId: number): Promise<{ message: string }> {
    return http.post('/anix/user/reading-history', {
      comicId,
      chapterId,
    });
  }

  /**
   * 获取阅读历史
   */
  async getReadingHistory(params?: { page?: number; limit?: number }): Promise<{
    list: Array<{
      comic: Comic;
      chapter: Chapter;
      readAt: string;
    }>;
    total: number;
    page: number;
    limit: number;
  }> {
    return http.get('/anix/user/reading-history', { params });
  }

  /**
   * 清除阅读历史
   */
  async clearReadingHistory(): Promise<{ message: string }> {
    return http.delete('/anix/user/reading-history');
  }
}

// 导出 API 实例
export const anixApi = new AnixApi();

// 导出默认实例
export default anixApi;
