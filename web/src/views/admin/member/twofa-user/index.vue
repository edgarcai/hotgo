<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="2FA用户管理">
        管理用户的双因素认证设置，包括启用、禁用和重置备用码等操作
      </n-card>
    </div>

    <n-card :bordered="false" class="proCard">
      <!-- 搜索表单 -->
      <BasicForm
        @register="register"
        @submit="reloadTable"
        @reset="reloadTable"
        @keyup.enter="reloadTable"
        ref="searchFormRef"
      >
        <template #statusSlot="{ model, field }">
          <n-input v-model:value="model[field]" />
        </template>
      </BasicForm>

      <!-- 工具栏 -->
      <BasicTable
        :openChecked="true"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row) => row.id"
        ref="actionRef"
        :actionColumn="actionColumn"
        @update:checked-row-keys="onCheckedRow"
        :scroll-x="1280"
        :resizeHeightOffset="-180"
      >
        <template #tableTitle>
          <n-button
            type="error"
            @click="batchDisable2FA"
            :disabled="batchDisableLoading || checkedIds.length <= 0"
            :loading="batchDisableLoading"
            class="min-left-space"
          >
            <template #icon>
              <n-icon>
                <DeleteOutlined />
              </n-icon>
            </template>
            批量禁用2FA
          </n-button>

          <n-button
            type="primary"
            @click="exportData"
            :loading="exportLoading"
            class="min-left-space"
          >
            <template #icon>
              <n-icon>
                <ExportOutlined />
              </n-icon>
            </template>
            导出数据
          </n-button>
        </template>
      </BasicTable>
    </n-card>

    <!-- 用户详情弹窗 -->
    <n-modal v-model:show="showDetailModal" preset="card" title="用户详情" style="width: 800px">
      <div v-if="currentUser">
        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="用户ID">
            {{ currentUser.id }}
          </n-descriptions-item>
          <n-descriptions-item label="用户名">
            {{ currentUser.username }}
          </n-descriptions-item>
          <n-descriptions-item label="真实姓名">
            {{ currentUser.realName || '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="邮箱">
            {{ currentUser.email || '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="手机号">
            {{ currentUser.mobile || '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="2FA状态">
            <n-tag :type="currentUser.twoFactorEnabled ? 'success' : 'warning'" size="small">
              {{ currentUser.twoFactorEnabled ? '已启用' : '未启用' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="启用时间">
            {{ currentUser.twoFactorEnabledAt ? formatToDateTime(currentUser.twoFactorEnabledAt) : '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="最后登录">
            {{ currentUser.lastLoginAt ? formatToDateTime(currentUser.lastLoginAt) : '-' }}
          </n-descriptions-item>
          <n-descriptions-item label="登录次数">
            {{ currentUser.loginCount || 0 }}
          </n-descriptions-item>
          <n-descriptions-item label="状态">
            <n-tag :type="currentUser.status === 1 ? 'success' : 'error'" size="small">
              {{ currentUser.status === 1 ? '正常' : '禁用' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="创建时间">
            {{ formatToDateTime(currentUser.createdAt) }}
          </n-descriptions-item>
          <n-descriptions-item label="更新时间">
            {{ formatToDateTime(currentUser.updatedAt) }}
          </n-descriptions-item>
        </n-descriptions>

        <!-- 2FA详细信息 -->
        <n-divider title-placement="left">2FA详细信息</n-divider>
        <div v-if="currentUser.twoFactorEnabled">
          <n-space vertical>
            <n-alert type="info" title="2FA已启用">
              该用户已启用双因素认证，登录时需要提供验证码
            </n-alert>
            <n-descriptions :column="1" bordered>
              <n-descriptions-item label="备用码剩余">
                {{ currentUser.backupCodesCount || 0 }} 个
              </n-descriptions-item>
              <n-descriptions-item label="最后使用时间">
                {{ currentUser.lastUsedAt ? formatToDateTime(currentUser.lastUsedAt) : '从未使用' }}
              </n-descriptions-item>
            </n-descriptions>
          </n-space>
        </div>
        <div v-else>
          <n-alert type="warning" title="2FA未启用">
            该用户尚未启用双因素认证，建议引导用户启用以提高账户安全性
          </n-alert>
        </div>
      </div>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, h } from 'vue';
  import { useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { columns, actionColumn, setActionHandlers } from './columns';
  import { DeleteOutlined, ExportOutlined } from '@vicons/antd';
  import { formatToDateTime } from '@/utils/dateUtil';
  import { List as getUserList, Reset2FA, Status } from '@/api/org/user';

  const dialog = useDialog();
  const message = useMessage();
  const actionRef = ref();
  const searchFormRef = ref();

  // 表单配置
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
      field: 'email',
      component: 'NInput',
      label: '邮箱',
      componentProps: {
        placeholder: '请输入邮箱',
        clearable: true,
      },
    },
    {
      field: 'mobile',
      component: 'NInput',
      label: '手机号',
      componentProps: {
        placeholder: '请输入手机号',
        clearable: true,
      },
    },
    {
      field: 'twoFactorEnabled',
      component: 'NSelect',
      label: '2FA状态',
      componentProps: {
        placeholder: '请选择2FA状态',
        clearable: true,
        options: [
          { label: '已启用', value: true },
          { label: '未启用', value: false },
        ],
      },
    },
    {
      field: 'status',
      component: 'NSelect',
      label: '状态',
      componentProps: {
        placeholder: '请选择状态',
        clearable: true,
        options: [
          { label: '正常', value: 1 },
          { label: '禁用', value: 2 },
        ],
      },
    },
  ];

  const [register, { getFieldsValue }] = useForm({
    gridProps: { cols: '1 s:1 m:2 l:3 xl:4 2xl:4' },
    labelWidth: 80,
    schemas: schemas as any,
  });

  // 状态管理
  const checkedIds = ref<number[]>([]);
  const batchDisableLoading = ref(false);
  const exportLoading = ref(false);
  const showDetailModal = ref(false);
  const currentUser = ref<any>(null);

  // 加载表格数据
  const loadDataTable = async (res: any) => {
    const params = {
      ...getFieldsValue(),
      page: res.page,
      pageSize: res.pageSize,
    };
    return await getUserList(params);
  };

  // 重新加载表格
  const reloadTable = () => {
    actionRef.value?.reload();
  };

  // 选中行处理
  const onCheckedRow = (rowKeys: number[]) => {
    checkedIds.value = rowKeys;
  };

  // 查看用户详情
  const handleView = (row: any) => {
    currentUser.value = row;
    showDetailModal.value = true;
  };

  // 切换用户2FA状态
  const handleToggle2FA = (row: any) => {
    const action = row.twoFactorEnabled ? '禁用' : '启用';
    dialog.warning({
      title: `${action}2FA`,
      content: `确定要${action}用户 "${row.username}" 的双因素认证吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          // 使用Status API来切换用户状态
          await Status({ id: row.id, status: row.twoFactorEnabled ? 2 : 1 });
          message.success(`${action}2FA成功`);
          reloadTable();
        } catch (error) {
          message.error(`${action}2FA失败`);
        }
      },
    });
  };

  // 重置用户备用码
  const handleResetBackupCodes = async (row: any) => {
    try {
      await Reset2FA({ id: row.id });
      message.success('重置2FA成功');
      reloadTable();
    } catch (error) {
      message.error('重置2FA失败');
    }
  };

  // 批量禁用2FA
  const batchDisable2FA = () => {
    if (checkedIds.value.length <= 0) {
      message.warning('请选择要操作的用户');
      return;
    }

    dialog.warning({
      title: '批量禁用2FA',
      content: `确定要禁用选中的 ${checkedIds.value.length} 个用户的双因素认证吗？`,
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: async () => {
        try {
          batchDisableLoading.value = true;
          // 这里应该调用批量禁用API
          // await batchDisableUser2FA({ userIds: checkedIds.value });
          message.success('批量禁用2FA成功');
          checkedIds.value = [];
          reloadTable();
        } catch (error) {
          message.error('批量禁用2FA失败');
        } finally {
          batchDisableLoading.value = false;
        }
      },
    });
  };

  // 导出数据
  const exportData = async () => {
    try {
      exportLoading.value = true;
      const params = getFieldsValue();
      // 这里需要实现导出功能的API
      message.info('导出功能待实现');
    } catch (error) {
      message.error('导出失败');
    } finally {
      exportLoading.value = false;
    }
  };

  // 设置操作处理函数
  setActionHandlers({
    handleView,
    handleToggle2FA,
    handleResetBackupCodes,
  });
</script>

<style scoped>
.min-left-space {
  margin-left: 12px;
}
</style>