import { http } from '@/utils/http/axios';

/**
 * 获取用户2FA状态
 */
export function get2FAStatus() {
  return http.request({
    url: '/admin2fa/status',
    method: 'get',
  }, {
    // 禁用错误消息显示，因为2FA状态检查失败是正常情况
    isShowErrorMessage: false,
  });
}

/**
 * 设置2FA（生成密钥和二维码）
 */
export function setup2FA() {
  return http.request({
    url: '/admin2fa/setup',
    method: 'post',
  });
}

/**
 * 验证2FA设置
 */
export function verifySetup2FA(params: { code: string }) {
  return http.request({
    url: '/admin2fa/verify-setup',
    method: 'post',
    params,
  });
}

/**
 * 禁用2FA
 */
export function disable2FA(params: { code: string }) {
  return http.request({
    url: '/admin2fa/disable',
    method: 'post',
    params,
  });
}

/**
 * 验证2FA代码
 */
export function verify2FA(params: { code: string }) {
  return http.request({
    url: '/admin2fa/verify',
    method: 'post',
    params,
  });
}

/**
 * 获取备用码
 */
export function getBackupCodes() {
  return http.request({
    url: '/admin2fa/backup-codes',
    method: 'get',
  });
}

/**
 * 重新生成备用码
 */
export function regenerateBackupCodes(params: { code: string }) {
  return http.request({
    url: '/admin2fa/regenerate-backup-codes',
    method: 'post',
    params,
  });
}

/**
 * 使用备用码验证
 */
export function useBackupCode(params: { code: string }) {
  return http.request({
    url: '/admin2fa/use-backup-code',
    method: 'post',
    params,
  });
}

/**
 * 获取2FA日志列表
 */
export function get2FALogList(params: any) {
  return http.request({
    url: '/admin2fa/log/list',
    method: 'get',
    params,
  });
}