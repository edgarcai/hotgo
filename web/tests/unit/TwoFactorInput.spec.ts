import { mount } from '@vue/test-utils';
import { nextTick, h, defineComponent, ref } from 'vue';
import { describe, it, expect, vi } from 'vitest';
import TwoFactorInput from '@/components/TwoFactor/TwoFactorInput.vue';

// 创建Naive UI组件的测试桩，使其在单测环境下渲染为原生元素并透传事件
const NInputStub = defineComponent({
  // 函数说明：模拟Naive UI的n-input，使v-model:value与@input、@keyup.enter事件可用，并暴露focus方法
  props: ['value', 'disabled', 'status', 'placeholder', 'maxlength', 'size'],
  emits: ['update:value', 'input', 'keyup.enter'],
  setup(props, { emit, attrs, expose }) {
    const inputEl = ref<HTMLInputElement | null>(null);
    expose({
      // 暴露focus方法，供父组件通过ref调用
      focus: () => inputEl.value && inputEl.value.focus(),
    });
    return () =>
      h('input', {
        ...attrs,
        ref: inputEl,
        value: (props as any).value,
        disabled: (props as any).disabled,
        onInput: (e: any) => {
          const val = e?.target?.value ?? '';
          emit('update:value', val);
          emit('input', val);
        },
        onKeyup: (e: KeyboardEvent) => {
          if (e.key === 'Enter') emit('keyup.enter', e);
        },
      });
  },
});

const PassThroughStub = defineComponent({
  // 函数说明：简单的透传插槽的容器桩
  setup(_, { slots }) {
    return () => h('div', {}, slots.default ? slots.default() : []);
  },
});

const ButtonStub = defineComponent({
  // 函数说明：模拟n-button，点击时触发click事件
  emits: ['click'],
  setup(_, { emit, slots }) {
    return () =>
      h(
        'button',
        {
          onClick: () => emit('click'),
        },
        slots.default ? slots.default() : []
      );
  },
});

const IconStub = defineComponent({
  // 函数说明：图标组件桩
  setup() {
    return () => h('span');
  },
});

const STUBS = {
  'n-input': NInputStub,
  'n-form-item': PassThroughStub,
  'n-progress': PassThroughStub,
  'n-text': PassThroughStub,
  'n-alert': PassThroughStub,
  'n-button': ButtonStub,
  'n-icon': IconStub,
  // 图标组件
  SafetyCertificateOutlined: IconStub,
  KeyOutlined: IconStub,
};

// TwoFactorInput 组件单元测试
// 覆盖：TOTP模式6位自动提交、备用码模式回车提交、错误信息展示

describe('TwoFactorInput 组件', () => {
  // 基础渲染测试：TOTP模式并输入6位自动提交
  it('渲染TOTP模式并输入6位自动提交', async () => {
    const onSubmit = vi.fn();
    const wrapper = mount(TwoFactorInput as any, {
      props: {
        mode: 'totp',
        loading: false,
        error: undefined,
        onSubmit: onSubmit,
      },
      global: {
        stubs: STUBS,
      },
    });

    // 应至少存在一个输入框
    const inputs = wrapper.findAll('input');
    expect(inputs.length).toBeGreaterThanOrEqual(1);

    // 通过data-testid定位到n-input内部的原生input
    const input = wrapper.find('[data-testid="totp-input"]');
    await input.setValue('123456');

    await nextTick();
    expect(onSubmit).toHaveBeenCalledWith('123456');
  });

  // 备用码模式：输入并按回车提交
  it('切换到备用码模式并触发提交', async () => {
    const onSubmit = vi.fn();
    const onSwitch = vi.fn();
    const wrapper = mount(TwoFactorInput as any, {
      props: {
        mode: 'backup',
        loading: false,
        error: undefined,
        onSubmit,
        onSwitchMode: onSwitch,
      },
      global: {
        stubs: STUBS,
      },
    });

    const backup = wrapper.find('[data-testid="backup-input"]');
    await backup.setValue('ABCD-EFGH-IJKL');
    await backup.trigger('keyup.enter');
    // 组件会标准化为小写且移除非字母数字
    expect(onSubmit).toHaveBeenCalledWith('abcdefghijkl');
  });

  // 错误信息展示
  it('显示错误信息', async () => {
    const wrapper = mount(TwoFactorInput as any, {
      props: {
        mode: 'totp',
        loading: false,
        error: '验证码错误',
        onSubmit: vi.fn(),
        onSwitchMode: vi.fn(),
      },
      global: {
        stubs: STUBS,
      },
    });

    expect(wrapper.find('[data-testid="error-message"]').exists()).toBe(true);
  });
});