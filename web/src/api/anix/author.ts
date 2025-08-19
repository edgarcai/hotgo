import { http } from '@/utils/http/axios';

// 作者信息
export interface AuthorInfo {
  id: number;
  name: string;
  avatar: string;
  bio: string;
  comicCount: number;
  totalViews: number;
  status: number;
  createdAt: string;
  updatedAt: string;
}

// 作者列表请求参数
export interface AuthorListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: number;
}

// 作者列表响应
export interface AuthorListResponse {
  list: AuthorInfo[];
  total: number;
  page: number;
  pageSize: number;
}

// 获取作者列表
export function getAuthorList(params: AuthorListParams = {}) {
  return http.request<AuthorListResponse>({
    url: '/api/anix/author/list',
    method: 'GET',
    params,
  });
}

// 获取所有启用的作者（用于筛选）
export function getAllAuthors() {
  return http.request<AuthorInfo[]>({
    url: '/api/anix/author/list',
    method: 'GET',
    params: {
      status: 1,
      pageSize: 100,
    },
  });
}

// 获取作者详情
export function getAuthorDetail(id: number) {
  return http.request<AuthorInfo>({
    url: `/api/anix/author/view`,
    method: 'GET',
    params: { id },
  });
}