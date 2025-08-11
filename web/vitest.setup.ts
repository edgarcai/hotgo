import { config } from '@vue/test-utils';

// 解决部分浏览器API在jsdom中的缺失
class ResizeObserverMock {
  observe() {}
  unobserve() {}
  disconnect() {}
}

// @ts-ignore
global.ResizeObserver = global.ResizeObserver || ResizeObserverMock;

config.global.mocks = {
  $t: (t: string) => t,
};