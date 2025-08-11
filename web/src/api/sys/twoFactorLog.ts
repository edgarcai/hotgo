import { http } from '@/utils/http/axios';
import { BasicResponseModel, BasicPageParams } from '@/api/system/user';

/**
 * 2FA操作日志相关接口
 */

// 2FA操作日志列表项
export interface TwoFactorLogItem {
  id: number;
  memberId: number;
  username: string;
  operationType: string; // 操作类型：2fa_enable,2fa_disable,2fa_reset,2fa_verify_success,2fa_verify_failed
  operationDesc: string; // 操作描述
  ip: string;
  createdAt: string;
}

// 2FA操作日志详情
export interface TwoFactorLogDetail {
  id: number;
  memberId: number;
  username: string;
  operationType: string;
  operationDesc: string;
  ip: string;
  userAgent: string;
  details: string;
  createdAt: string;
}

// 2FA操作日志列表查询参数
export interface TwoFactorLogListParams extends BasicPageParams {
  memberId?: number;
  username?: string;
  operationType?: string;
  ip?: string;
  createdAt?: string[];
}

// 2FA操作日志列表响应
export interface TwoFactorLogListResponse {
  list: TwoFactorLogItem[];
  totalCount: number;
  page: number;
  perPage: number;
}

/**
 * 获取2FA操作日志列表
 * @param params 查询参数
 */
export function getTwoFactorLogList(params: TwoFactorLogListParams) {
  return http.request<BasicResponseModel<TwoFactorLogListResponse>>({
    url: '/api/admin/sys/two-factor-log/list',
    method: 'GET',
    params,
  });
}

/**
 * 获取2FA操作日志详情
 * @param id 日志ID
 */
export function getTwoFactorLogDetail(id: number) {
  return http.request<BasicResponseModel<TwoFactorLogDetail>>({
    url: '/api/admin/sys/two-factor-log/view',
    method: 'GET',
    params: { id },
  });
}

/**
 * 删除2FA操作日志
 * @param id 日志ID
 */
export function deleteTwoFactorLog(id: number) {
  return http.request<BasicResponseModel<{ message: string }>>({
    url: '/api/admin/sys/two-factor-log/delete',
    method: 'POST',
    data: { id },
  });
}