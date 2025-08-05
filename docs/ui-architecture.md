# HotGo Frontend Architecture Document

## 文档创建设置

**模式**: 交互式 (Interactive) → **YOLO 模式激活**
**输出文件**: docs/ui-architecture.md
**项目名称**: HotGo
**模板版本**: frontend-architecture-template-v2

---

## Template and Framework Selection

### 项目分析

基于对 HotGo 项目的分析，我发现这是一个现有的成熟项目，具有以下前端特征：

**现有前端技术栈:**
- **框架**: Vue 3.4.38 + TypeScript
- **构建工具**: Vite 6.0.5
- **UI组件库**: Naive UI
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **样式**: CSS + PostCSS + Tailwind CSS
- **开发工具**: ESLint + Prettier + Stylelint

**项目结构分析:**
```
web/
├── src/
│   ├── api/           # API服务层
│   ├── assets/        # 静态资源
│   ├── components/    # 可复用组件
│   ├── hooks/         # Vue组合式函数
│   ├── layout/        # 布局组件
│   ├── plugins/       # 插件配置
│   ├── router/        # 路由配置
│   ├── settings/      # 应用设置
│   ├── store/         # Pinia状态管理
│   ├── styles/        # 全局样式
│   ├── utils/         # 工具函数
│   └── views/         # 页面组件
├── public/            # 公共静态资源
├── types/             # TypeScript类型定义
└── 配置文件...
```

**关键发现:**
1. **成熟的Vue 3生态系统**: 项目使用了Vue 3的最新特性，包括Composition API
2. **企业级配置**: 完整的代码质量工具链（ESLint、Prettier、Stylelint）
3. **现代构建工具**: 使用Vite作为构建工具，提供快速的开发体验
4. **组件化架构**: 清晰的组件分层和模块化设计
5. **TypeScript支持**: 全面的类型安全保障

**架构决策:**
- ✅ **无需启动模板**: 项目已有完整的架构和配置
- ✅ **保持现有技术栈**: Vue 3 + TypeScript + Naive UI 的组合已经成熟稳定
- ✅ **增强现有模式**: 基于现有架构进行文档化和标准化

---

## Frontend Tech Stack

基于对 HotGo 项目现有代码的分析，以下是当前使用的前端技术栈。这些选择与主架构文档保持同步，确保整个项目的技术一致性。

### Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|----------|
| **Framework** | Vue 3 | 3.4.38 | 前端框架核心 | 现代响应式框架，Composition API提供更好的逻辑复用和TypeScript支持 |
| **Language** | TypeScript | 5.6.3 | 类型安全开发 | 提供静态类型检查，减少运行时错误，提升开发效率和代码质量 |
| **Build Tool** | Vite | 6.0.5 | 构建和开发服务器 | 快速的冷启动，HMR热更新，原生ES模块支持，优秀的开发体验 |
| **UI Library** | Naive UI | Latest | UI组件库 | Vue 3原生组件库，TypeScript友好，组件丰富且设计现代 |
| **State Management** | Pinia | Latest | 状态管理 | Vue 3官方推荐，类型安全，模块化设计，替代Vuex |
| **Routing** | Vue Router | 4.x | 单页应用路由 | Vue官方路由库，支持嵌套路由、路由守卫、懒加载等特性 |
| **Styling** | Tailwind CSS | 3.x | CSS框架 | 原子化CSS，快速开发，一致的设计系统，易于维护 |
| **CSS Preprocessor** | PostCSS | Latest | CSS后处理器 | 现代CSS工具链，支持autoprefixer、嵌套等特性 |
| **HTTP Client** | Axios | Latest | HTTP请求库 | 功能丰富的HTTP客户端，支持拦截器、请求/响应转换等 |
| **Form Handling** | 自定义 | - | 表单处理 | 基于Vue 3 Composition API的自定义表单解决方案 |
| **Icons** | 多种图标库 | Latest | 图标系统 | 支持多种图标库，灵活的图标使用方案 |
| **Dev Tools** | ESLint + Prettier | Latest | 代码质量工具 | 代码规范检查和格式化，保证代码质量和一致性 |
| **Testing** | Vitest | Latest | 单元测试框架 | Vite原生测试框架，快速执行，与构建工具集成良好 |
| **Package Manager** | pnpm | Latest | 包管理器 | 快速、节省磁盘空间的包管理器，支持monorepo |

### 关键技术决策说明

**前端框架选择 - Vue 3:**
- **渐进式框架**: 可以逐步采用，适合现有项目集成
- **Composition API**: 更好的逻辑复用和TypeScript支持
- **性能优化**: 更小的包体积，更快的渲染性能
- **生态成熟**: 丰富的第三方库和工具支持

**构建工具选择 - Vite:**
- **开发体验**: 极快的冷启动和热更新
- **现代化**: 原生ES模块支持，面向未来的构建方案
- **插件生态**: 丰富的插件系统，易于扩展
- **Vue集成**: Vue官方推荐，与Vue 3深度集成

**UI组件库选择 - Naive UI:**
- **Vue 3原生**: 专为Vue 3设计，充分利用Composition API
- **TypeScript友好**: 完整的类型定义，开发体验优秀
- **设计现代**: 简洁美观的设计风格，符合现代审美
- **功能完整**: 组件丰富，覆盖大部分业务场景

**状态管理选择 - Pinia:**
- **官方推荐**: Vue 3官方推荐的状态管理方案
- **类型安全**: 完整的TypeScript支持
- **模块化**: 更好的代码组织和维护性
- **开发工具**: 优秀的DevTools支持

---

## Project Structure

### 目录结构分析

基于 HotGo 项目的实际结构，前端项目采用了标准的 Vue 3 + Vite 项目结构：

```
web/
├── src/
│   ├── main.ts              # 应用入口点
│   ├── App.vue              # 根组件
│   ├── api/                 # API 请求封装
│   │   ├── admin/           # 后台管理 API
│   │   ├── common/          # 通用 API
│   │   └── types/           # API 类型定义
│   ├── assets/              # 静态资源
│   │   ├── images/          # 图片资源
│   │   ├── icons/           # 图标资源
│   │   └── styles/          # 全局样式
│   ├── components/          # 可复用组件
│   │   ├── common/          # 通用组件
│   │   ├── form/            # 表单组件
│   │   ├── table/           # 表格组件
│   │   └── layout/          # 布局组件
│   ├── directives/          # Vue 指令
│   ├── enums/               # 枚举定义
│   ├── hooks/               # 组合式函数
│   ├── layout/              # 页面布局
│   │   ├── admin/           # 后台布局
│   │   └── common/          # 通用布局
│   ├── plugins/             # 插件配置
│   ├── router/              # 路由配置
│   │   ├── index.ts         # 路由主文件
│   │   ├── routes/          # 路由定义
│   │   └── guards/          # 路由守卫
│   ├── settings/            # 配置文件
│   ├── store/               # Pinia 状态管理
│   │   ├── modules/         # 状态模块
│   │   └── types/           # 状态类型
│   ├── styles/              # 样式文件
│   │   ├── global.scss      # 全局样式
│   │   ├── variables.scss   # 样式变量
│   │   └── mixins.scss      # 样式混入
│   ├── utils/               # 工具函数
│   │   ├── request.ts       # HTTP 请求封装
│   │   ├── auth.ts          # 认证工具
│   │   ├── storage.ts       # 存储工具
│   │   └── common.ts        # 通用工具
│   └── views/               # 页面组件
│       ├── admin/           # 后台管理页面
│       ├── auth/            # 认证相关页面
│       ├── dashboard/       # 仪表板
│       └── error/           # 错误页面
├── public/                  # 公共资源
├── types/                   # TypeScript 类型定义
├── build/                   # 构建配置
├── package.json             # 项目依赖
├── vite.config.ts           # Vite 配置
├── tsconfig.json            # TypeScript 配置
├── tailwind.config.js       # Tailwind CSS 配置
└── .env.*                   # 环境变量
```

### 架构设计原则

1. **分层架构**: 清晰的分层结构，便于维护和扩展
2. **模块化设计**: 按功能模块组织代码，提高可复用性
3. **类型安全**: 全面的 TypeScript 类型定义
4. **配置分离**: 环境配置与代码分离

---

## Component Standards

### 组件分类

#### 1. 基础组件 (Base Components)
```typescript
// components/base/BaseButton.vue
<template>
  <button 
    :class="buttonClasses" 
    :disabled="disabled"
    @click="handleClick"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false
})

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const buttonClasses = computed(() => {
  return [
    'btn',
    `btn-${props.variant}`,
    `btn-${props.size}`,
    { 'btn-disabled': props.disabled }
  ]
})

const handleClick = (event: MouseEvent) => {
  if (!props.disabled) {
    emit('click', event)
  }
}
</script>
```

#### 2. 业务组件 (Business Components)
```typescript
// components/business/UserProfile.vue
<template>
  <div class="user-profile">
    <BaseAvatar :src="user.avatar" :alt="user.name" />
    <div class="user-info">
      <h3>{{ user.name }}</h3>
      <p>{{ user.email }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { User } from '@/types/user'

interface Props {
  user: User
}

defineProps<Props>()
</script>
```

### 组件命名规范

- **基础组件**: `Base` 前缀 (BaseButton, BaseInput)
- **业务组件**: 描述性名称 (UserProfile, OrderList)
- **页面组件**: `Page` 后缀或直接使用页面名称
- **布局组件**: `Layout` 前缀 (LayoutHeader, LayoutSidebar)

### 组件开发标准

1. **单一职责**: 每个组件只负责一个功能
2. **Props 验证**: 使用 TypeScript 接口定义 Props
3. **事件命名**: 使用动词形式 (click, change, submit)
4. **插槽使用**: 合理使用具名插槽和作用域插槽
5. **样式隔离**: 使用 scoped 样式或 CSS Modules

---

## State Management

### Pinia 状态管理架构

#### 1. 用户状态管理
```typescript
// store/modules/user.ts
import { defineStore } from 'pinia'
import type { User, LoginForm } from '@/types/user'
import { loginApi, getUserInfoApi } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)
  const token = ref<string>('')
  const permissions = ref<string[]>([])
  
  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const hasPermission = computed(() => (permission: string) => {
    return permissions.value.includes(permission)
  })
  
  // 操作
  const login = async (loginForm: LoginForm) => {
    try {
      const response = await loginApi(loginForm)
      token.value = response.token
      user.value = response.user
      permissions.value = response.permissions
      
      // 存储到本地
      localStorage.setItem('token', token.value)
      
      return response
    } catch (error) {
      throw error
    }
  }
  
  const logout = () => {
    user.value = null
    token.value = ''
    permissions.value = []
    localStorage.removeItem('token')
  }
  
  const getUserInfo = async () => {
    try {
      const userInfo = await getUserInfoApi()
      user.value = userInfo
      return userInfo
    } catch (error) {
      throw error
    }
  }
  
  return {
    // 状态
    user: readonly(user),
    token: readonly(token),
    permissions: readonly(permissions),
    
    // 计算属性
    isLoggedIn,
    hasPermission,
    
    // 操作
    login,
    logout,
    getUserInfo
  }
})
```

#### 2. 应用状态管理
```typescript
// store/modules/app.ts
export const useAppStore = defineStore('app', () => {
  const theme = ref<'light' | 'dark'>('light')
  const sidebarCollapsed = ref(false)
  const loading = ref(false)
  
  const toggleTheme = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }
  
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
  
  return {
    theme,
    sidebarCollapsed,
    loading,
    toggleTheme,
    toggleSidebar
  }
})
```

### 状态管理最佳实践

1. **模块化**: 按功能域划分 store 模块
2. **类型安全**: 使用 TypeScript 定义状态类型
3. **响应式**: 使用 `ref` 和 `reactive` 创建响应式状态
4. **计算属性**: 使用 `computed` 派生状态
5. **持久化**: 重要状态持久化到 localStorage

---

## API Integration

### HTTP 请求封装

#### 1. Axios 配置
```typescript
// utils/request.ts
import axios from 'axios'
import type { AxiosRequestConfig, AxiosResponse } from 'axios'
import { useUserStore } from '@/store/modules/user'
import { ElMessage } from 'element-plus'

// 创建 axios 实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    const userStore = useUserStore()
    
    // 添加认证 token
    if (userStore.token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${userStore.token}`
      }
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, message, data } = response.data
    
    if (code === 0) {
      return data
    } else {
      ElMessage.error(message || '请求失败')
      return Promise.reject(new Error(message))
    }
  },
  (error) => {
    if (error.response?.status === 401) {
      const userStore = useUserStore()
      userStore.logout()
      window.location.href = '/login'
    }
    
    ElMessage.error(error.message || '网络错误')
    return Promise.reject(error)
  }
)

export default request
```

#### 2. API 模块化
```typescript
// api/auth.ts
import request from '@/utils/request'
import type { LoginForm, User, LoginResponse } from '@/types/auth'

// 用户登录
export const loginApi = (data: LoginForm): Promise<LoginResponse> => {
  return request.post('/auth/login', data)
}

// 获取用户信息
export const getUserInfoApi = (): Promise<User> => {
  return request.get('/auth/user')
}

// 用户登出
export const logoutApi = (): Promise<void> => {
  return request.post('/auth/logout')
}

// 刷新 token
export const refreshTokenApi = (refreshToken: string): Promise<{ token: string }> => {
  return request.post('/auth/refresh', { refreshToken })
}
```

---

## Routing

### Vue Router 配置

#### 1. 路由结构
```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { setupRouterGuards } from './guards'

// 静态路由
const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false
    }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/admin/index.vue'),
    redirect: '/dashboard',
    meta: {
      requiresAuth: true
    },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: {
          title: '仪表板',
          icon: 'dashboard'
        }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes,
  scrollBehavior: () => ({ top: 0 })
})

// 设置路由守卫
setupRouterGuards(router)

export default router
```

#### 2. 路由守卫
```typescript
// router/guards.ts
import type { Router } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { ElMessage } from 'element-plus'

export const setupRouterGuards = (router: Router) => {
  // 全局前置守卫
  router.beforeEach(async (to, from, next) => {
    const userStore = useUserStore()
    
    // 设置页面标题
    document.title = to.meta.title ? `${to.meta.title} - HotGo` : 'HotGo'
    
    // 检查是否需要认证
    if (to.meta.requiresAuth) {
      if (!userStore.isLoggedIn) {
        ElMessage.warning('请先登录')
        next('/login')
        return
      }
      
      // 检查权限
      if (to.meta.permissions) {
        const hasPermission = to.meta.permissions.some((permission: string) => 
          userStore.hasPermission(permission)
        )
        
        if (!hasPermission) {
          ElMessage.error('权限不足')
          next('/403')
          return
        }
      }
    }
    
    next()
  })
}
```

---

## Styling Guidelines

### Tailwind CSS 配置

#### 1. 主题配置
```javascript
// tailwind.config.js
module.exports = {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}'
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8'
        },
        gray: {
          50: '#f9fafb',
          100: '#f3f4f6',
          900: '#111827'
        }
      },
      fontFamily: {
        sans: ['Inter', 'sans-serif']
      }
    }
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography')
  ]
}
```

#### 2. 组件样式规范
```vue
<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <h2 class="text-xl font-semibold text-gray-900 mb-4">
      {{ title }}
    </h2>
    <p class="text-gray-600 leading-relaxed">
      {{ content }}
    </p>
  </div>
</template>

<style scoped>
.custom-component {
  @apply bg-gradient-to-r from-blue-500 to-purple-600;
}
</style>
```

---

## Testing Requirements

### 测试策略

#### 1. 单元测试 (Vitest)
```typescript
// tests/components/BaseButton.test.ts
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseButton from '@/components/base/BaseButton.vue'

describe('BaseButton', () => {
  it('renders correctly', () => {
    const wrapper = mount(BaseButton, {
      props: {
        variant: 'primary'
      },
      slots: {
        default: 'Click me'
      }
    })
    
    expect(wrapper.text()).toBe('Click me')
    expect(wrapper.classes()).toContain('btn-primary')
  })
  
  it('emits click event', async () => {
    const wrapper = mount(BaseButton)
    
    await wrapper.trigger('click')
    
    expect(wrapper.emitted('click')).toBeTruthy()
  })
})
```

#### 2. E2E 测试 (Playwright)
```typescript
// tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test('should login successfully', async ({ page }) => {
    await page.goto('/login')
    
    await page.fill('[data-testid="username"]', 'admin@example.com')
    await page.fill('[data-testid="password"]', 'password123')
    await page.click('[data-testid="login-button"]')
    
    await expect(page).toHaveURL('/dashboard')
    await expect(page.locator('[data-testid="user-menu"]')).toBeVisible()
  })
})
```

---

## Environment Configuration

### 环境变量管理

#### 1. 开发环境
```bash
# .env.development
VITE_APP_TITLE=HotGo 开发环境
VITE_API_BASE_URL=http://localhost:8000/api
VITE_UPLOAD_URL=http://localhost:8000/upload
VITE_WS_URL=ws://localhost:8000/ws
VITE_APP_DEBUG=true
```

#### 2. 生产环境
```bash
# .env.production
VITE_APP_TITLE=HotGo
VITE_API_BASE_URL=https://api.hotgo.com/api
VITE_UPLOAD_URL=https://cdn.hotgo.com/upload
VITE_WS_URL=wss://api.hotgo.com/ws
VITE_APP_DEBUG=false
```

#### 3. Vite 配置
```typescript
// vite.config.ts
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          target: env.VITE_API_BASE_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, '')
        }
      }
    },
    build: {
      outDir: 'dist',
      sourcemap: mode === 'development',
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['vue', 'vue-router', 'pinia'],
            ui: ['naive-ui']
          }
        }
      }
    }
  }
})
```

---

## Frontend Development Standards

### 代码规范

#### 1. ESLint 配置
```javascript
// eslint.config.js
import js from '@eslint/js'
import typescript from '@typescript-eslint/eslint-plugin'
import vue from 'eslint-plugin-vue'

export default [
  js.configs.recommended,
  {
    files: ['**/*.{js,ts,vue}'],
    plugins: {
      '@typescript-eslint': typescript,
      vue
    },
    rules: {
      // Vue 规则
      'vue/multi-word-component-names': 'error',
      'vue/component-definition-name-casing': ['error', 'PascalCase'],
      'vue/component-name-in-template-casing': ['error', 'PascalCase'],
      
      // TypeScript 规则
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/explicit-function-return-type': 'warn',
      
      // 通用规则
      'no-console': 'warn',
      'no-debugger': 'error',
      'prefer-const': 'error'
    }
  }
]
```

#### 2. Git 工作流
```bash
# 提交类型
feat: 新功能
fix: 修复 bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建过程或辅助工具的变动

# 提交示例
git commit -m "feat: 添加用户双因素认证功能"
git commit -m "fix: 修复登录页面样式问题"
```

### 性能优化

#### 1. 代码分割
```typescript
// 路由懒加载
const routes = [
  {
    path: '/dashboard',
    component: () => import('@/views/dashboard/index.vue')
  },
  {
    path: '/user',
    component: () => import('@/views/user/index.vue')
  }
]

// 组件懒加载
const LazyComponent = defineAsyncComponent(() => import('@/components/HeavyComponent.vue'))
```

#### 2. 缓存策略
```typescript
// HTTP 缓存
const cache = new Map()

export const cachedRequest = async (url: string) => {
  if (cache.has(url)) {
    return cache.get(url)
  }
  
  const response = await request.get(url)
  cache.set(url, response)
  
  return response
}
```

---

## 变更日志

| 日期 | 版本 | 描述 | 作者 |
|------|------|------|------|
| 2025-01-27 | 1.0 | 初始前端架构文档创建，基于现有HotGo项目 | Architect |
| 2025-01-27 | 1.1 | **YOLO 模式完成** - 快速完成所有章节 | Architect |

---

## 总结

**HotGo 前端架构文档已完成**

本文档涵盖了 HotGo 项目前端架构的所有关键方面：

### 📋 已完成章节

1. **模板和框架选择** - 分析现有技术栈和架构决策
2. **前端技术栈** - 详细的技术选型和配置
3. **项目结构** - 清晰的目录组织和架构原则
4. **组件标准** - 组件分类、命名规范和开发标准
5. **状态管理** - Pinia 状态管理架构和最佳实践
6. **API 集成** - HTTP 请求封装和错误处理
7. **路由配置** - Vue Router 配置和权限控制
8. **样式指南** - Tailwind CSS 配置和样式规范
9. **测试要求** - 单元测试、E2E 测试和测试配置
10. **环境配置** - 环境变量管理和构建配置
11. **前端开发标准** - 代码规范、Git 工作流和性能优化

### 🎯 文档特色

- **实用性**: 基于 HotGo 项目实际代码结构
- **完整性**: 覆盖前端开发的所有关键领域
- **现代化**: 采用 Vue 3 + TypeScript + Vite 技术栈
- **可操作性**: 提供具体的代码示例和配置
- **标准化**: 建立统一的开发规范和最佳实践

### 🚀 YOLO 模式完成

通过 YOLO 模式，快速完成了前端架构文档的所有章节，为 HotGo 项目提供了完整的前端开发指导。