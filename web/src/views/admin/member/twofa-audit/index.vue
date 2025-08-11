<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="2FA操作日志">
        查看和管理双因素认证相关的操作日志，包括启用、禁用、重置和验证记录
      </n-card>
    </div>
    <n-card :bordered="false" class="proCard">
      <BasicForm
        @register="register"
        @submit="reloadTable"
        @reset="reloadTable"
        @keyup.enter="reloadTable"
        ref="searchFormRef"
      />

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        ref="actionRef"
        :actionColumn="actionColumn"
        :scroll-x="scrollX"
        :resizeHeightOffset="-10000"
        size="small"
      >
        <template #tableTitle>
          <n-button
            type="error"
            @click="handleBatchDelete"
            :disabled="batchDeleteDisabled"
            class="min-left-space"
            v-if="hasPermission(['/admin/2fa/logs/delete'])"
          >
            <template #icon>
              <n-icon>
                <DeleteOutlined />
              </n-icon>
            </template>
            批量删除
          </n-button>
        </template>
      </BasicTable>
    </n-card>

    <!-- 日志详情弹窗 -->
    <n-modal v-model:show="showDetailModal" :show-icon="false" preset="dialog" title="日志详情" style="width: 600px">
      <div v-if="currentLog" class="log-detail">
        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="用户名">
            {{ currentLog.username }}
          </n-descriptions-item>
          <n-descriptions-item label="操作类型">
            <n-tag :type="getOperationTypeColor(currentLog.operationType)">
              {{ getOperationTypeText(currentLog.operationType) }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="操作描述">
            {{ currentLog.operationDesc }}
          </n-descriptions-item>
          <n-descriptions-item label="IP地址">
            {{ currentLog.ip }}
          </n-descriptions-item>
          <n-descriptions-item label="用户代理" :span="2">
            {{ currentLog.userAgent }}
          </n-descriptions-item>
          <n-descriptions-item label="操作时间">
            {{ formatToDateTime(currentLog.createdAt) }}
          </n-descriptions-item>
          <n-descriptions-item label="详细信息" :span="2" v-if="currentLog.details">
            <n-code :code="currentLog.details" language="json" />
          </n-descriptions-item>
        </n-descriptions>
      </div>
      <template #action>
        <n-space>
          <n-button @click="showDetailModal = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, reactive, ref } from 'vue';
  import { useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { usePermission } from '@/hooks/web/usePermission';
  import { getTwoFactorLogList, getTwoFactorLogDetail, deleteTwoFactorLog } from '@/api/sys/twoFactorLog';
  import { columns } from './columns';
  import { DeleteOutlined } from '@vicons/antd';
  import { adaTableScrollX } from '@/utils/hotgo';
  import { formatToDateTime } from '@/utils/dateUtil';
  import type { TwoFactorLogDetail } from '@/api/sys/twoFactorLog';

  const dialog = useDialog();
  const message = useMessage();
  const { hasPermission } = usePermission();
  const actionRef = ref();
  const searchFormRef = ref<any>({});
  const batchDeleteDisabled = ref(true);
  const checkedIds = ref([]);
  const showDetailModal = ref(false);
  const currentLog = ref<TwoFactorLogDetail | null>(null);

  // 搜索表单配置
  const schemas = [
    {
      field: 'username',
      component: 'NInput',
      label: '用户名',
      componentProps: {
        placeholder: '请输入用户名',
        clearable: true,
      },
    },
    {
      field: 'operationType',
      component: 'NSelect',
      label: '操作类型',
      componentProps: {
        placeholder: '请选择操作类型',
        clearable: true,
        options: [
          { label: '启用2FA', value: '2fa_enable' },
          { label: '禁用2FA', value: '2fa_disable' },
          { label: '重置2FA', value: '2fa_reset' },
          { label: '验证成功', value: '2fa_verify_success' },
          { label: '验证失败', value: '2fa_verify_failed' },
        ],
      },
    },
    {
      field: 'ip',
      component: 'NInput',
      label: 'IP地址',
      componentProps: {
        placeholder: '请输入IP地址',
        clearable: true,
      },
    },
    {
      field: 'createdAt',
      component: 'NDatePicker',
      label: '操作时间',
      componentProps: {
        type: 'datetimerange',
        clearable: true,
        format: 'yyyy-MM-dd HH:mm:ss',
        valueFormat: 'yyyy-MM-dd HH:mm:ss',
      },
    },
  ];

  // 操作列配置
  const actionColumn = reactive({
    width: 120,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: any) {
      return h(TableAction as any, {
        style: 'button',
        actions: [
          {
            label: '查看',
            onClick: handleView.bind(null, record),
            auth: ['/admin/2fa/logs/view'],
          },
          {
            label: '删除',
            onClick: handleDelete.bind(null, record),
            auth: ['/admin/2fa/logs/delete'],
          },
        ],
      });
    },
  });

  const scrollX = computed(() => {
    return adaTableScrollX(columns, actionColumn.width);
  });

  const [register] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas,
  });

  // 加载表格数据
  const loadDataTable = async (res: any) => {
    const params = { ...(searchFormRef.value?.formModel ?? {}), ...res };
    return await getTwoFactorLogList(params);
  };

  // 重新加载表格
  const reloadTable = () => {
    actionRef.value?.reload();
  };

  // 查看详情
  const handleView = async (record: any) => {
    try {
      const response = await getTwoFactorLogDetail(record.id);
      currentLog.value = response.data;
      showDetailModal.value = true;
    } catch (error) {
      message.error('获取日志详情失败');
    }
  };

  // 删除单条记录
  const handleDelete = (record: any) => {
    dialog.warning({
      title: '警告',
      content: `确定删除用户 ${record.username} 的操作日志吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          await deleteTwoFactorLog(record.id);
          message.success('删除成功');
          reloadTable();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  // 批量删除
  const handleBatchDelete = () => {
    if (checkedIds.value.length === 0) {
      message.warning('请选择要删除的记录');
      return;
    }
    dialog.warning({
      title: '警告',
      content: `确定删除选中的 ${checkedIds.value.length} 条日志记录吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          // 批量删除逻辑
          for (const id of checkedIds.value) {
            await deleteTwoFactorLog(id as number);
          }
          message.success('批量删除成功');
          checkedIds.value = [];
          reloadTable();
        } catch (error) {
          message.error('批量删除失败');
        }
      },
    });
  };

  // 获取操作类型文本
  const getOperationTypeText = (type: string) => {
    const typeMap: Record<string, string> = {
      '2fa_enable': '启用2FA',
      '2fa_disable': '禁用2FA',
      '2fa_reset': '重置2FA',
      '2fa_verify_success': '验证成功',
      '2fa_verify_failed': '验证失败',
    };
    return typeMap[type] || type;
  };

  // 获取操作类型颜色
  const getOperationTypeColor = (type: string) => {
    const colorMap: Record<string, string> = {
      '2fa_enable': 'success',
      '2fa_disable': 'warning',
      '2fa_reset': 'error',
      '2fa_verify_success': 'info',
      '2fa_verify_failed': 'error',
    };
    return colorMap[type] || 'default';
  };

  // 选中行变化
  const onCheckedRow = (rowKeys: any[]) => {
    checkedIds.value = rowKeys;
    batchDeleteDisabled.value = rowKeys.length === 0;
  };
</script>

<style scoped>
.log-detail {
  padding: 16px 0;
}
</style>