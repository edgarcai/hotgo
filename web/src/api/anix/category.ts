import { http } from '@/utils/http/axios';

// 分类信息
export interface CategoryInfo {
  id: number;
  name: string;
  description: string;
  icon: string;
  sort: number;
  status: number;
  createdAt: string;
  updatedAt: string;
}

// 分类列表请求参数
export interface CategoryListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  status?: number;
}

// 分类列表响应
export interface CategoryListResponse {
  list: CategoryInfo[];
  total: number;
  page: number;
  pageSize: number;
}

// 获取分类列表
export function getCategoryList(params: CategoryListParams = {}) {
  return http.request<CategoryListResponse>({
    url: '/api/anix/category/list',
    method: 'GET',
    params,
  });
}

// 获取所有启用的分类（用于筛选）
export function getAllCategories() {
  return http.request<CategoryInfo[]>({
    url: '/api/anix/category/list',
    method: 'GET',
    params: {
      status: 1,
      pageSize: 100,
    },
  });
}

// 获取分类详情
export function getCategoryDetail(id: number) {
  return http.request<CategoryInfo>({
    url: `/api/anix/category/view`,
    method: 'GET',
    params: { id },
  });
}