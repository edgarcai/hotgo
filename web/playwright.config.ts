import { defineConfig, devices } from '@playwright/test';

/**
 * @see https://playwright.dev/docs/test-configuration
 */
export default defineConfig({
  testDir: './tests/e2e',
  /* 并行运行测试 */
  fullyParallel: false,
  /* 失败时禁止重试，用于CI */
  forbidOnly: !!process.env.CI,
  /* CI重试次数 */
  retries: process.env.CI ? 2 : 0,
  /* CI中的工作线程限制 */
  workers: process.env.CI ? 1 : undefined,
  /* 测试报告配置 */
  reporter: [
    ['html', { outputFolder: 'tests/reports/html' }],
    ['junit', { outputFile: 'tests/reports/junit.xml' }],
    ['line']
  ],
  /* 全局测试配置 */
  use: {
    /* 基础URL */
    baseURL: 'http://localhost:8002',
    /* 收集失败测试的痕迹 */
    trace: 'on-first-retry',
    /* 失败时截图 */
    screenshot: 'only-on-failure',
    /* 页面加载状态检测 */
    // 使用默认load/waitUntil策略，不额外指定网络空闲等待
    /* 响应超时时间 */
    actionTimeout: 30000,
    /* 导航超时时间 */
    navigationTimeout: 30000,
  },

  /* 配置不同的浏览器项目 */
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },

    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },

    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },

    /* 移动端测试 */
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },

    /* 品牌浏览器测试 */
    // {
    //   name: 'Microsoft Edge',
    //   use: { ...devices['Desktop Edge'], channel: 'msedge' },
    // },
    // {
    //   name: 'Google Chrome',
    //   use: { ...devices['Desktop Chrome'], channel: 'chrome' },
    // },
  ],

  /* 启动开发服务器 */
  webServer: {
    command: 'pnpm dev',
    url: 'http://localhost:8002',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000,
  },
});