import { mount, flushPromises } from '@vue/test-utils'
import { nextTick, h, defineComponent } from 'vue'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'

// 模拟 naive-ui 的 useMessage，使用全局对象避免 hoisting 问题
vi.mock('naive-ui', () => ({
  useMessage: () => (globalThis as any).__useMessage || { error: vi.fn(), success: vi.fn(), warning: vi.fn() },
}))

// 模拟图标组件
vi.mock('@vicons/ionicons5', () => ({
  ShieldCheckmarkOutline: defineComponent({
    name: 'ShieldCheckmarkOutline',
    render: () => h('svg'),
  }),
  ArrowBackOutline: defineComponent({
    name: 'ArrowBackOutline',
    render: () => h('svg'),
  }),
}))

// 模拟环境配置
vi.mock('@/utils/env', () => ({
  getAppEnvConfig: () => ({
    VITE_GLOB_APP_TITLE: 'HotGo Test',
  }),
}))

// 模拟 useUserStore
const userStoreMock = {
  setToken: vi.fn(),
  GetInfo: vi.fn(),
}

vi.mock('@/store/modules/user', () => ({
  useUserStore: () => userStoreMock,
}))

// 导入需要测试的组件和相关依赖
import Verify2FA from '@/views/auth/Verify2FA.vue'
import * as twoFactorAPI from '@/api/auth/twoFactor'
import { ResultEnum } from '@/enums/httpEnum'

// 创建测试用的 Naive UI 组件桩
const NCardStub = defineComponent({
  name: 'NCard',
  props: ['bordered'],
  setup(_, { slots }) {
    return () => h('div', { class: 'n-card' }, slots.default?.())
  },
})

const NSpaceStub = defineComponent({
  name: 'NSpace',
  props: ['justify', 'align', 'size'],
  setup(_, { slots }) {
    return () => h('div', { class: 'n-space' }, slots.default?.())
  },
})

const NIconStub = defineComponent({
  name: 'NIcon',
  props: ['size', 'color'],
  setup(_, { slots }) {
    return () => h('span', { class: 'n-icon' }, slots.default?.())
  },
})

const NAvatarStub = defineComponent({
  name: 'NAvatar',
  props: ['size', 'src'],
  setup(_, { slots }) {
    return () => h('img', { class: 'n-avatar' }, slots.default?.())
  },
})

const NButtonStub = defineComponent({
  name: 'NButton',
  props: ['text', 'disabled'],
  emits: ['click'],
  setup(props, { emit, slots }) {
    return () =>
      h(
        'button',
        {
          class: 'n-button',
          onClick: () => emit('click'),
          disabled: props.disabled,
        },
        slots.default?.()
      )
  },
})

const NModalStub = defineComponent({
  name: 'NModal',
  props: ['show', 'preset', 'title'],
  emits: ['update:show'],
  setup(props, { slots }) {
    // 根据 props.show 动态渲染，避免仅在 setup 时判断导致后续无法更新
    return () => (props.show ? h('div', { class: 'n-modal' }, slots.default?.()) : null)
  },
})

// 模拟 TwoFactorInput 组件
const TwoFactorInputStub = defineComponent({
  name: 'TwoFactorInput',
  props: ['mode', 'loading', 'error'],
  emits: ['submit', 'switch-mode', 'update:error'],
  setup(props, { emit, expose }) {
    const clear = vi.fn()
    const focus = vi.fn()

    // 暴露方法给父组件使用
    expose({ clear, focus })

    return () =>
      h('div', {
        'data-testid': 'two-factor-input',
        class: 'two-factor-input-stub',
      })
  },
})

// 模拟 2FA API
const mockVerifyLogin2FA = vi.fn()
vi.spyOn(twoFactorAPI, 'verifyLogin2FA').mockImplementation(mockVerifyLogin2FA)

describe('Verify2FA 页面集成测试', () => {
  let router: any
  let pinia: any

  beforeEach(() => {
    // 重置所有mock
    vi.clearAllMocks()

    // 创建路由实例
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Home</div>' } },
        { path: '/login', name: 'Login', component: { template: '<div>Login</div>' } },
        { path: '/verify-2fa', name: 'Verify2FA', component: Verify2FA },
      ],
    })

    // 创建 Pinia 实例
    pinia = createPinia()

    // 模拟localStorage
    Object.defineProperty(window, 'localStorage', {
      value: {
        setItem: vi.fn(),
        getItem: vi.fn(),
        removeItem: vi.fn(),
      },
    })

    // 设置 useMessage mock 的默认返回值到全局
    ;(globalThis as any).__useMessage = {
      error: vi.fn(),
      success: vi.fn(),
      warning: vi.fn(),
    }
  })

  it('应正确渲染页面基础结构', async () => {
    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
        username: 'testuser',
      },
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 检查页面主体结构
    expect(wrapper.find('.verify-2fa-container').exists()).toBe(true)
    expect(wrapper.find('.verify-2fa-card').exists()).toBe(true)

    // 检查头部信息
    expect(wrapper.find('.verify-header').exists()).toBe(true)
    expect(wrapper.text()).toContain('双因素认证')
    expect(wrapper.text()).toContain('请完成身份验证以继续登录')

    // 检查用户信息区域
    expect(wrapper.find('.user-info').exists()).toBe(true)

    // 检查 TwoFactorInput 组件
    expect(wrapper.findComponent(TwoFactorInputStub).exists()).toBe(true)

    // 检查底部操作按钮
    expect(wrapper.find('.footer-actions').exists()).toBe(true)
    expect(wrapper.text()).toContain('返回登录')
    expect(wrapper.text()).toContain('需要帮助？')
  })

  it('缺少 tempToken 时应显示错误并跳转到登录页', async () => {
    const message = { error: vi.fn(), success: vi.fn(), warning: vi.fn() }
    ;(globalThis as any).__useMessage = message

    // 不提供 tempToken 参数
    await router.push('/verify-2fa')

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()
    await flushPromises()

    // 应该显示错误消息
    expect(message.error).toHaveBeenCalledWith('缺少必要的验证参数，请重新登录')
    
    // 应该跳转到登录页面
    await nextTick()
    await flushPromises()
    expect(router.currentRoute.value.path).toBe('/login')
  })

  it('2FA验证成功应设置token并跳转', async () => {
    const message = { error: vi.fn(), success: vi.fn(), warning: vi.fn() }
    ;(globalThis as any).__useMessage = message

    // 提供必要的路由参数
    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
        username: 'testuser',
        redirect: '/dashboard',
      },
    })

    // 模拟API成功响应
    mockVerifyLogin2FA.mockResolvedValue({
      code: ResultEnum.SUCCESS,
      data: {
        token: 'final-access-token',
      },
    })

    // 模拟用户信息获取成功
    userStoreMock.GetInfo.mockResolvedValue({
      id: 1,
      username: 'testuser',
      realName: '测试用户',
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 获取 TwoFactorInput 组件并触发提交
    const twoFactorInput = wrapper.findComponent(TwoFactorInputStub)
    await twoFactorInput.vm.$emit('submit', '123456')

    await nextTick()
    await flushPromises()

    // 验证API调用
    expect(mockVerifyLogin2FA).toHaveBeenCalledWith({
      tempToken: 'test-temp-token',
      code: '123456',
      type: 'totp',
    })

    // 验证用户Store操作
    expect(userStoreMock.setToken).toHaveBeenCalledWith('final-access-token')
    expect(userStoreMock.GetInfo).toHaveBeenCalled()

    // 验证localStorage操作
    expect(localStorage.setItem).toHaveBeenCalledWith('ACCESS_TOKEN', 'final-access-token')
    expect(localStorage.setItem).toHaveBeenCalledWith(
      'CURRENT_USER',
      JSON.stringify({
        id: 1,
        username: 'testuser',
        realName: '测试用户',
      })
    )

    // 验证成功消息
    expect(message.success).toHaveBeenCalledWith('验证成功，正在跳转...')

    // 验证跳转到目标页面
    await nextTick()
    await flushPromises()
    expect(router.currentRoute.value.path).toBe('/dashboard')
  })

  it('2FA验证失败应显示错误消息', async () => {
    const message = { error: vi.fn(), success: vi.fn(), warning: vi.fn() }
    ;(globalThis as any).__useMessage = message

    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
        username: 'testuser',
      },
    })

    // 模拟API失败响应
    mockVerifyLogin2FA.mockResolvedValue({
      code: -1,
      message: '验证码错误',
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 触发提交
    const twoFactorInput = wrapper.findComponent(TwoFactorInputStub)
    await twoFactorInput.vm.$emit('submit', '000000')

    await nextTick()

    // 验证API调用
    expect(mockVerifyLogin2FA).toHaveBeenCalledWith({
      tempToken: 'test-temp-token',
      code: '000000',
      type: 'totp',
    })

    // 不应该设置token
    expect(userStoreMock.setToken).not.toHaveBeenCalled()
    
    // 验证错误状态透传到子组件props
    const updatedProps = wrapper.findComponent(TwoFactorInputStub).props()
    expect(updatedProps.error).toBe('验证码错误')
  })

  it('应支持验证模式切换', async () => {
    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
        username: 'testuser',
      },
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 初始模式应为 totp
    const twoFactorInput = wrapper.findComponent(TwoFactorInputStub)
    expect(twoFactorInput.props('mode')).toBe('totp')

    // 触发模式切换
    await twoFactorInput.vm.$emit('switch-mode')
    await nextTick()

    // 模式应切换为 backup（通过子组件props断言）
    expect(wrapper.findComponent(TwoFactorInputStub).props('mode')).toBe('backup')

    // 再次切换
    await twoFactorInput.vm.$emit('switch-mode')
    await nextTick()

    // 模式应切换回 totp（通过子组件props断言）
    expect(wrapper.findComponent(TwoFactorInputStub).props('mode')).toBe('totp')
  })

  it('返回登录按钮应正确跳转', async () => {
    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
        redirect: '/dashboard',
      },
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 查找返回登录按钮
    const buttons = wrapper.findAllComponents(NButtonStub)
    const backButton = buttons.find(button => 
      button.text().includes('返回登录')
    )
    expect(backButton).toBeTruthy()
    
    await backButton!.trigger('click')
    await nextTick()
    await flushPromises()
    
    // 验证跳转到登录页面，并携带redirect参数
    expect(router.currentRoute.value.path).toBe('/login')
    expect(router.currentRoute.value.query.redirect).toBe('/dashboard')
  })

  it('应正确显示帮助对话框', async () => {
    await router.push({
      path: '/verify-2fa',
      query: {
        tempToken: 'test-temp-token',
      },
    })

    const wrapper = mount(Verify2FA, {
      global: {
        plugins: [router, pinia],
        stubs: {
          NCard: NCardStub,
          NSpace: NSpaceStub,
          NIcon: NIconStub,
          NAvatar: NAvatarStub,
          NButton: NButtonStub,
          NModal: NModalStub,
          TwoFactorInput: TwoFactorInputStub,
        },
      },
    })

    await nextTick()

    // 初始时帮助对话框应该隐藏（通过DOM判断）
    expect(wrapper.find('.n-modal').exists()).toBe(false)

    // 查找帮助按钮
    const buttons = wrapper.findAllComponents(NButtonStub)
    const helpButton = buttons.find(button => 
      button.text().includes('需要帮助？')
    )
    expect(helpButton).toBeTruthy()
    
    await helpButton!.trigger('click')
    await nextTick()
    
    // 帮助对话框应该显示（通过DOM判断）
    expect(wrapper.find('.n-modal').exists()).toBe(true)
  })
})