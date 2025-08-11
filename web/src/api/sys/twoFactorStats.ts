import { http } from '@/utils/http/axios';
import { BasicResponseModel } from '@/api/system/user';

/**
 * 2FA统计相关接口
 */

// 2FA统计数据模型
export interface TwoFactorStatsData {
  totalUsers: number; // 总用户数
  enabledUsers: number; // 已启用2FA用户数
  disabledUsers: number; // 未启用2FA用户数
  enableRate: number; // 启用率(%)
  todayOperations: number; // 今日操作次数
  weekOperations: number; // 本周操作次数
  monthOperations: number; // 本月操作次数
  successfulLogins: number; // 成功登录次数
  failedLogins: number; // 失败登录次数
  successRate: number; // 成功率(%)
}

// 2FA趋势数据项
export interface TwoFactorTrendItem {
  date: string; // 日期
  operations: number; // 操作次数
  successful: number; // 成功次数
  failed: number; // 失败次数
}

// 2FA趋势查询参数
export interface TwoFactorTrendParams {
  type: 'daily' | 'weekly' | 'monthly'; // 趋势类型
  dateRange: string[]; // 日期范围
}

// 2FA统计查询参数
export interface TwoFactorStatsParams {
  dateRange?: string[]; // 日期范围（可选）
}

/**
 * 获取2FA统计数据
 * @param params 查询参数
 */
export function getTwoFactorStats(params?: TwoFactorStatsParams) {
  return http.request<BasicResponseModel<TwoFactorStatsData>>({
    url: '/admin/api/admin/sys/two-factor-stats/stats',
    method: 'GET',
    params,
  });
}

/**
 * 获取2FA趋势数据
 * @param params 查询参数
 */
export function getTwoFactorTrend(params: TwoFactorTrendParams) {
  return http.request<BasicResponseModel<TwoFactorTrendItem[]>>({
    url: '/admin/api/admin/sys/two-factor-stats/trend',
    method: 'GET',
    params,
  });
}