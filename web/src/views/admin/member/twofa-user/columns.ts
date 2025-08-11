import { h } from 'vue';
import { NTag, NButton, NPopconfirm } from 'naive-ui';
import { BasicColumn } from '@/components/Table';
import { formatToDateTime } from '@/utils/dateUtil';

// 2FA用户管理列定义
export const columns: BasicColumn[] = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
  },
  {
    title: '用户名',
    key: 'username',
    width: 120,
  },
  {
    title: '真实姓名',
    key: 'realName',
    width: 100,
  },
  {
    title: '邮箱',
    key: 'email',
    width: 180,
  },
  {
    title: '手机号',
    key: 'mobile',
    width: 120,
  },
  {
    title: '2FA状态',
    key: 'twoFactorEnabled',
    width: 100,
    render(row: any) {
      return h(
        NTag,
        {
          type: row.twoFactorEnabled ? 'success' : 'warning',
          size: 'small',
        },
        {
          default: () => (row.twoFactorEnabled ? '已启用' : '未启用'),
        }
      );
    },
  },
  {
    title: '启用时间',
    key: 'twoFactorEnabledAt',
    width: 160,
    render(row: any) {
      return row.twoFactorEnabledAt ? formatToDateTime(row.twoFactorEnabledAt) : '-';
    },
  },
  {
    title: '最后登录',
    key: 'lastLoginAt',
    width: 160,
    render(row: any) {
      return row.lastLoginAt ? formatToDateTime(row.lastLoginAt) : '-';
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render(row: any) {
      const statusMap: Record<number, { type: any; text: string }> = {
        1: { type: 'success', text: '正常' },
        2: { type: 'error', text: '禁用' },
      };
      const status = statusMap[row.status] || { type: 'default', text: '未知' };
      return h(
        NTag,
        {
          type: status.type as any,
          size: 'small',
        },
        {
          default: () => status.text,
        }
      );
    },
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 160,
    render(row: any) {
      return formatToDateTime(row.createdAt);
    },
  },
];

// 操作列定义
export const actionColumn = {
  width: 200,
  title: '操作',
  key: 'action',
  fixed: 'right',
  render(row: any) {
    return [
      h(
        NButton,
        {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => handleView(row),
        },
        { default: () => '查看' }
      ),
      h(
        NButton,
        {
          size: 'small',
          type: row.twoFactorEnabled ? 'error' : 'success',
          secondary: true,
          style: { marginLeft: '8px' },
          onClick: () => handleToggle2FA(row),
        },
        { default: () => (row.twoFactorEnabled ? '禁用2FA' : '启用2FA') }
      ),
      h(
        NPopconfirm,
        {
          onPositiveClick: () => handleResetBackupCodes(row),
        },
        {
          default: () => '确定要重置该用户的备用码吗？',
          trigger: () =>
            h(
              NButton,
              {
                size: 'small',
                type: 'warning',
                secondary: true,
                style: { marginLeft: '8px' },
                disabled: !row.twoFactorEnabled,
              },
              { default: () => '重置备用码' }
            ),
        }
      ),
    ];
  },
};

// 操作处理函数（需要在组件中实现）
let handleView: (row: any) => void;
let handleToggle2FA: (row: any) => void;
let handleResetBackupCodes: (row: any) => void;

// 设置操作处理函数
export function setActionHandlers(handlers: {
  handleView: (row: any) => void;
  handleToggle2FA: (row: any) => void;
  handleResetBackupCodes: (row: any) => void;
}) {
  handleView = handlers.handleView;
  handleToggle2FA = handlers.handleToggle2FA;
  handleResetBackupCodes = handlers.handleResetBackupCodes;
}