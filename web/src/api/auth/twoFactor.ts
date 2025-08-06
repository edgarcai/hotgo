import { http } from '@/utils/http/axios'
import { BasicResponseModel } from '@/api/system/user'

/**
 * 2FA验证相关接口
 */

// 2FA验证请求参数
export interface VerifyLogin2FAParams {
  tempToken: string  // 临时token
  code: string      // 验证码（TOTP或备用码）
  type: 'totp' | 'backup'  // 验证类型
}

// 2FA验证响应数据
export interface VerifyLogin2FAResponse {
  token: string     // 最终的访问token
  expires: number   // token过期时间
  user: {
    id: number
    username: string
    realName: string
    avatar: string
    email: string
    mobile: string
    permissions: string[]
    [key: string]: any
  }
}

// 2FA状态查询响应
export interface TwoFactorStatusResponse {
  enabled: boolean      // 是否已启用2FA
  enabledAt?: string   // 启用时间
  lastUsedAt?: string  // 最后使用时间
  backupCodesCount: number  // 剩余备用码数量
}

// 2FA启用请求参数
export interface Enable2FAParams {
  password: string  // 当前密码确认
}

// 2FA启用响应数据
export interface Enable2FAResponse {
  qrCodeUrl: string    // 二维码URL
  secretKey: string    // 密钥（用于手动输入）
  backupCodes: string[] // 备用恢复码
}

// 2FA设置验证参数
export interface VerifySetup2FAParams {
  code: string      // TOTP验证码
  secretKey: string // 密钥
}

// 2FA禁用参数
export interface Disable2FAParams {
  password: string  // 当前密码确认
  code: string     // TOTP验证码
}

/**
 * 验证登录2FA
 * @param params 验证参数
 */
export function verifyLogin2FA(params: VerifyLogin2FAParams) {
  return http.request<BasicResponseModel<VerifyLogin2FAResponse>>({
    url: '/api/admin/auth/verify-login-2fa',
    method: 'POST',
    params
  }, {
    isTransformResponse: false
  })
}

/**
 * 获取2FA状态
 */
export function getTwoFactorStatus() {
  return http.request<BasicResponseModel<TwoFactorStatusResponse>>({
    url: '/api/admin/auth/2fa-status',
    method: 'GET'
  })
}

/**
 * 启用2FA
 * @param params 启用参数
 */
export function enable2FA(params: Enable2FAParams) {
  return http.request<BasicResponseModel<Enable2FAResponse>>({
    url: '/api/admin/auth/enable-2fa',
    method: 'POST',
    params
  })
}

/**
 * 验证2FA设置
 * @param params 验证参数
 */
export function verifySetup2FA(params: VerifySetup2FAParams) {
  return http.request<BasicResponseModel<{ backupCodes: string[] }>>({
    url: '/api/admin/auth/verify-setup-2fa',
    method: 'POST',
    params
  })
}

/**
 * 禁用2FA
 * @param params 禁用参数
 */
export function disable2FA(params: Disable2FAParams) {
  return http.request<BasicResponseModel<{ success: boolean }>>({
    url: '/api/admin/auth/disable-2fa',
    method: 'POST',
    params
  })
}

/**
 * 重新生成备用码
 * @param password 当前密码
 */
export function regenerateBackupCodes(password: string) {
  return http.request<BasicResponseModel<{ backupCodes: string[] }>>({
    url: '/api/admin/auth/regenerate-backup-codes',
    method: 'POST',
    params: { password }
  })
}

/**
 * 获取2FA设置二维码
 */
export function get2FAQRCode() {
  return http.request<BasicResponseModel<{ qrCodeUrl: string; secretKey: string }>>({
    url: '/api/admin/auth/2fa-qrcode',
    method: 'GET'
  })
}