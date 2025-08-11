import { test, expect } from '@playwright/test';

/**
 * 2FA 设置界面端到端测试
 * 覆盖：
 * - 登录
 * - 进入2FA设置页
 * - 启用2FA向导（仅验证UI流程，不校验真实验证码）
 * - 查看备用码弹窗
 * - 禁用2FA表单打开
 */

const ADMIN_USERNAME = process.env.E2E_ADMIN_USERNAME || 'admin';
const ADMIN_PASSWORD = process.env.E2E_ADMIN_PASSWORD || '123456';

// 根据项目路由守卫逻辑，先登录后进入设置页
async function login(page) {
  // 打开登录页
  await page.goto('/#/login');

  // 输入账号密码
  await page.getByPlaceholder('账号').fill(ADMIN_USERNAME);
  await page.getByPlaceholder('密码').fill(ADMIN_PASSWORD);

  // 点击登录
  await page.getByRole('button', { name: '登录' }).click();

  // 断言进入主页（根据项目路由，默认进入控制台或首页）
  await expect(page).toHaveURL(/#\//, { timeout: 15000 });
}

// 导航到 2FA 设置页面（系统里存在两个：系统-安全-TwoFactorAuth 以及 用户-安全-TwoFactorSettings）
async function navigateTo2FASettings(page) {
  // 通过动态菜单搜索进入，或者直接尝试访问用户安全设置页面
  await page.goto('/#/security/TwoFactorSettings');
  // 如果页面不存在，尝试访问系统-安全页面
  if (!(await page.locator('.two-factor-settings').first().isVisible().catch(() => false))) {
    await page.goto('/#/system/security/TwoFactorAuth');
    await expect(page.locator('.two-factor-auth')).toBeVisible();
  }
}

// 基础可见性断言：头部、按钮等
async function assertBaseUI(page) {
  // 标题存在
  const title = page.getByText(/双因素认证/);
  await expect(title).toBeVisible();
}

/**
 * 主流程
 */
 test.describe('2FA 设置界面', () => {
  test('可以登录并访问2FA设置页', async ({ page }) => {
    await login(page);
    await navigateTo2FASettings(page);
    await assertBaseUI(page);
  });

  test('可以打开启用2FA向导并进行UI交互', async ({ page }) => {
    await login(page);
    await navigateTo2FASettings(page);

    // 点击启用按钮（若已启用，此按钮不存在；测试仅验证弹窗流程存在即可）
    const enableBtn = page.getByRole('button', { name: /启用双因素认证/ });
    if (await enableBtn.isVisible().catch(() => false)) {
      await enableBtn.click();

      // 弹出向导
      const wizardModal = page.locator('.enable-wizard-modal');
      await expect(wizardModal).toBeVisible();

      // 步骤1：输入密码（不提交真实后端验证，仅检查控件）
      const pwdInput = wizardModal.getByPlaceholder('请输入当前密码');
      await expect(pwdInput).toBeVisible();

      // 关闭向导
      await page.keyboard.press('Escape');
    }
  });

  test('已启用时可以查看与重新生成备用码UI', async ({ page }) => {
    await login(page);
    await navigateTo2FASettings(page);

    // 查看备用码按钮
    const viewCodes = page.getByRole('button', { name: /查看备用码/ });
    if (await viewCodes.isVisible().catch(() => false)) {
      await viewCodes.click();
      // 弹窗应可见
      await expect(page.getByText('备用恢复码')).toBeVisible();
      // 关闭弹窗
      await page.keyboard.press('Escape');
    }

    // 重新生成备用码
    const regenBtn = page.getByRole('button', { name: /重新生成备用码/ });
    if (await regenBtn.isVisible().catch(() => false)) {
      await regenBtn.click();
      // 页面会出现prompt要求输入密码（TwoFactorAuth.vue中），Playwright默认无法处理prompt，先接受取消
      page.once('dialog', async (dialog) => {
        await dialog.dismiss();
      });
    }
  });

  test('可以打开禁用2FA表单UI', async ({ page }) => {
    await login(page);
    await navigateTo2FASettings(page);

    const disableBtn = page.getByRole('button', { name: /禁用双因素认证/ });
    if (await disableBtn.isVisible().catch(() => false)) {
      await disableBtn.click();
      // 表单可见
      const disableForm = page.locator('.disable-2fa-form');
      await expect(disableForm).toBeVisible();
      // 关闭
      await page.keyboard.press('Escape');
    }
  });
});