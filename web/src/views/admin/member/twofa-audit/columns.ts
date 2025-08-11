import { h } from 'vue';
import { NTag, NButton, NPopconfirm, NSpace } from 'naive-ui';
import { BasicColumn } from '@/components/Table';
import { formatToDateTime } from '@/utils/dateUtil';

/**
 * 2FA操作日志表格列定义
 */
export const columns: BasicColumn[] = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    align: 'center',
  },
  {
    title: '用户名',
    key: 'username',
    width: 120,
    align: 'center',
  },
  {
    title: '操作类型',
    key: 'operationType',
    width: 140,
    align: 'center',
    render(row: any) {
      const typeMap: Record<string, { text: string; type: string }> = {
        '2fa_enable': { text: '启用2FA', type: 'success' },
        '2fa_disable': { text: '禁用2FA', type: 'warning' },
        '2fa_reset': { text: '重置2FA', type: 'error' },
        '2fa_verify_success': { text: '验证成功', type: 'info' },
        '2fa_verify_failed': { text: '验证失败', type: 'error' },
      };
      const config = typeMap[row.operationType] || { text: row.operationType, type: 'default' };
      return h(NTag, { type: config.type as any }, { default: () => config.text });
    },
  },
  {
    title: '操作描述',
    key: 'operationDesc',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: 'IP地址',
    key: 'ip',
    width: 140,
    align: 'center',
  },
  {
    title: '操作时间',
    key: 'createdAt',
    width: 180,
    align: 'center',
    render(row: any) {
      return formatToDateTime(row.createdAt);
    },
  },
];

/**
 * 操作列定义
 */
export const actionColumn = {
  width: 120,
  title: '操作',
  key: 'action',
  fixed: 'right',
  render(record: any) {
    return h(NSpace, { justify: 'center' }, {
      default: () => [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            ghost: true,
            onClick: () => handleView(record),
          },
          { default: () => '查看' }
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDelete(record),
          },
          {
            default: () => '确定删除这条日志吗？',
            trigger: () => h(
              NButton,
              {
                size: 'small',
                type: 'error',
                ghost: true,
              },
              { default: () => '删除' }
            ),
          }
        ),
      ],
    });
  },
};

// 操作处理函数（需要在组件中实现）
function handleView(record: any) {
  // 在组件中实现查看详情逻辑
}

function handleDelete(record: any) {
  // 在组件中实现删除逻辑
}