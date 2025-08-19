import type { RouteRecordRaw } from 'vue-router';

const anixRoutes: RouteRecordRaw[] = [
  {
    path: '/anix',
    name: 'Anix',
    redirect: '/anix/home',
    meta: {
      title: 'AniX漫画',
      sort: 1,
    },
    children: [
      {
        path: '/anix/home',
        name: 'AnixHome',
        component: () => import('@/views/anix/home/index.vue'),
        meta: {
          title: '首页',
        },
      },
      {
        path: '/anix/register',
        name: 'AnixRegister',
        component: () => import('@/views/anix/register/index.vue'),
        meta: {
          title: '用户注册',
        },
      },
      {
        path: '/anix/login',
        name: 'AnixLogin',
        component: () => import('@/views/anix/login/index.vue'),
        meta: {
          title: '用户登录',
        },
      },
      {
        path: '/anix/profile',
        name: 'AnixProfile',
        component: () => import('@/views/anix/profile/index.vue'),
        meta: {
          title: '个人中心',
          requiresAuth: true, // 需要登录才能访问
        },
      },
      {
        path: '/anix/comic/:id',
        name: 'AnixComicDetail',
        component: () => import('@/views/anix/comic/detail.vue'),
        meta: {
          title: '漫画详情',
        },
      },
      {
        path: '/anix/chapter/:id',
        name: 'AnixChapterReader',
        component: () => import('@/views/anix/chapter/index.vue'),
        meta: {
          title: '章节阅读',
        },
      },
      {
        path: '/anix/search',
        name: 'AnixSearch',
        component: () => import('@/views/anix/search/index.vue'),
        meta: {
          title: '搜索漫画',
        },
      },
    ],
  },
];

export default anixRoutes;
