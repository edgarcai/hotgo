import { http } from '@/utils/http/axios';

// 漫画列表请求参数
export interface ComicListParams {
  page?: number;
  pageSize?: number;
  keyword?: string;
  categoryId?: number;
  authorId?: number;
  status?: number;
  isVip?: number;
  tags?: string;
  orderBy?: string;
  orderType?: string;
}

// 漫画基础信息
export interface ComicInfo {
  id: number;
  title: string;
  description: string;
  cover: string;
  authorId: number;
  authorName: string;
  status: number;
  isVip: number;
  tags: string;
  viewCount: number;
  likeCount: number;
  chapterCount: number;
  lastChapterTitle: string;
  publishedAt: string;
  updatedAt: string;
  createdAt: string;
}

// 漫画详情信息
export interface ComicDetail extends ComicInfo {
  categories: Array<{
    id: number;
    name: string;
  }>;
  chapters: Array<{
    id: number;
    title: string;
    chapterNumber: number;
    pageCount: number;
    isVip: number;
    publishedAt: string;
  }>;
}

// 漫画列表响应
export interface ComicListResponse {
  list: ComicInfo[];
  total: number;
  page: number;
  pageSize: number;
}

// 获取漫画列表
export function getComicList(params: ComicListParams = {}) {
  return http.request<ComicListResponse>({
    url: '/api/anix/comic/list',
    method: 'GET',
    params,
  });
}

// 获取漫画详情
export function getComicDetail(id: number) {
  return http.request<ComicDetail>({
    url: `/api/anix/comic/view`,
    method: 'GET',
    params: { id },
  });
}

// 获取热门漫画
export function getHotComics(limit: number = 10) {
  return http.request<ComicInfo[]>({
    url: '/api/anix/comic/list',
    method: 'GET',
    params: {
      pageSize: limit,
      orderBy: 'view_count',
      orderType: 'desc',
    },
  });
}

// 获取最新更新漫画
export function getLatestComics(limit: number = 10) {
  return http.request<ComicInfo[]>({
    url: '/api/anix/comic/list',
    method: 'GET',
    params: {
      pageSize: limit,
      orderBy: 'updated_at',
      orderType: 'desc',
    },
  });
}

// 获取推荐漫画
export function getRecommendedComics(limit: number = 10) {
  return http.request<ComicInfo[]>({
    url: '/api/anix/comic/list',
    method: 'GET',
    params: {
      pageSize: limit,
      orderBy: 'like_count',
      orderType: 'desc',
    },
  });
}

// 搜索漫画
export function searchComics(params: ComicListParams) {
  return http.request<ComicListResponse>({
    url: '/api/anix/search/comics',
    method: 'GET',
    params,
  });
}