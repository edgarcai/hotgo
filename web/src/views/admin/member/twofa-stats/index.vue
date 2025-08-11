<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="2FA统计报告">
        查看双因素认证的使用统计和趋势分析，帮助了解系统安全状况
      </n-card>
    </div>

    <!-- 统计卡片 -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" class="mb-4">
      <n-grid-item>
        <n-card :bordered="false" class="stat-card">
          <n-statistic label="总用户数" :value="statsData.totalUsers">
            <template #prefix>
              <n-icon color="#18a058">
                <UserOutlined />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :bordered="false" class="stat-card">
          <n-statistic label="已启用2FA" :value="statsData.enabledUsers">
            <template #prefix>
              <n-icon color="#2080f0">
                <SafetyOutlined />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :bordered="false" class="stat-card">
          <n-statistic label="启用率" :value="statsData.enableRate" suffix="%">
            <template #prefix>
              <n-icon color="#f0a020">
                <PercentageOutlined />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :bordered="false" class="stat-card">
          <n-statistic label="今日操作" :value="statsData.todayOperations">
            <template #prefix>
              <n-icon color="#d03050">
                <BarChartOutlined />
              </n-icon>
            </template>
          </n-statistic>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 详细统计 -->
    <n-grid :cols="2" :x-gap="16" :y-gap="16" class="mb-4">
      <n-grid-item>
        <n-card :bordered="false" title="操作统计">
          <n-space vertical>
            <n-statistic label="本周操作次数" :value="statsData.weekOperations" />
            <n-statistic label="本月操作次数" :value="statsData.monthOperations" />
          </n-space>
        </n-card>
      </n-grid-item>
      <n-grid-item>
        <n-card :bordered="false" title="验证统计">
          <n-space vertical>
            <n-statistic label="成功登录" :value="statsData.successfulLogins" />
            <n-statistic label="失败登录" :value="statsData.failedLogins" />
            <n-statistic label="成功率" :value="statsData.successRate" suffix="%" />
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 趋势图表 -->
    <n-card :bordered="false" title="操作趋势" class="mb-4">
      <template #header-extra>
        <n-space>
          <n-select
            v-model:value="trendType"
            :options="trendTypeOptions"
            @update:value="loadTrendData"
            style="width: 120px"
          />
          <n-date-picker
            v-model:value="dateRange"
            type="daterange"
            clearable
            @update:value="loadTrendData"
          />
        </n-space>
      </template>
      <div ref="chartRef" style="height: 400px"></div>
    </n-card>

    <!-- 用户分布 -->
    <n-card :bordered="false" title="用户2FA状态分布">
      <div class="user-distribution">
        <n-progress
          type="circle"
          :percentage="statsData.enableRate"
          :stroke-width="8"
          style="margin-right: 32px"
        >
          <span class="progress-text">
            {{ statsData.enableRate.toFixed(1) }}%
          </span>
        </n-progress>
        <div class="distribution-info">
          <div class="info-item">
            <n-tag type="success" size="small">已启用</n-tag>
            <span class="info-value">{{ statsData.enabledUsers }} 用户</span>
          </div>
          <div class="info-item">
            <n-tag type="warning" size="small">未启用</n-tag>
            <span class="info-value">{{ statsData.disabledUsers }} 用户</span>
          </div>
        </div>
      </div>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted, onUnmounted, nextTick } from 'vue';
  import { useMessage } from 'naive-ui';
  import { getTwoFactorStats, getTwoFactorTrend } from '@/api/sys/twoFactorStats';
  import type { TwoFactorStatsData, TwoFactorTrendItem } from '@/api/sys/twoFactorStats';
  import { UserOutlined, SafetyOutlined, PercentageOutlined, BarChartOutlined } from '@vicons/antd';
  import * as echarts from 'echarts';

  const message = useMessage();
  const chartRef = ref<HTMLElement>();
  let chartInstance: echarts.ECharts | null = null;

  // 统计数据
  const statsData = ref<TwoFactorStatsData>({
    totalUsers: 0,
    enabledUsers: 0,
    disabledUsers: 0,
    enableRate: 0,
    todayOperations: 0,
    weekOperations: 0,
    monthOperations: 0,
    successfulLogins: 0,
    failedLogins: 0,
    successRate: 0,
  });

  // 趋势数据
  const trendData = ref<TwoFactorTrendItem[]>([]);
  const trendType = ref('daily');
  const dateRange = ref<[number, number] | null>(null);

  // 趋势类型选项
  const trendTypeOptions = [
    { label: '日趋势', value: 'daily' },
    { label: '周趋势', value: 'weekly' },
    { label: '月趋势', value: 'monthly' },
  ];

  // 加载统计数据
  const loadStatsData = async () => {
    try {
      const response = await getTwoFactorStats();
      statsData.value = response.data;
    } catch (error) {
      message.error('加载统计数据失败');
    }
  };

  // 加载趋势数据
  const loadTrendData = async () => {
    try {
      const params = {
        type: trendType.value as 'daily' | 'weekly' | 'monthly',
        dateRange: dateRange.value ? [
          new Date(dateRange.value[0]).toISOString().split('T')[0],
          new Date(dateRange.value[1]).toISOString().split('T')[0],
        ] : [],
      };
      const response = await getTwoFactorTrend(params);
      trendData.value = response.data;
      updateChart();
    } catch (error) {
      message.error('加载趋势数据失败');
    }
  };

  // 初始化图表
  const initChart = () => {
    if (!chartRef.value) return;
    chartInstance = echarts.init(chartRef.value);
    updateChart();
  };

  // 更新图表
  const updateChart = () => {
    if (!chartInstance || !trendData.value.length) return;

    const dates = trendData.value.map(item => item.date);
    const operations = trendData.value.map(item => item.operations);
    const successful = trendData.value.map(item => item.successful);
    const failed = trendData.value.map(item => item.failed);

    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
        },
      },
      legend: {
        data: ['总操作', '成功', '失败'],
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dates,
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          name: '总操作',
          type: 'line',
          stack: 'Total',
          data: operations,
          itemStyle: {
            color: '#5470c6',
          },
        },
        {
          name: '成功',
          type: 'line',
          stack: 'Total',
          data: successful,
          itemStyle: {
            color: '#91cc75',
          },
        },
        {
          name: '失败',
          type: 'line',
          stack: 'Total',
          data: failed,
          itemStyle: {
            color: '#ee6666',
          },
        },
      ],
    };

    chartInstance.setOption(option);
  };

  // 组件挂载后初始化
  onMounted(async () => {
    await loadStatsData();
    await loadTrendData();
    await nextTick();
    initChart();

    // 监听窗口大小变化
    window.addEventListener('resize', () => {
      chartInstance?.resize();
    });
  });

  // 组件卸载时清理
  onUnmounted(() => {
    chartInstance?.dispose();
    window.removeEventListener('resize', () => {
      chartInstance?.resize();
    });
  });
</script>

<style scoped>
.stat-card {
  text-align: center;
  padding: 16px;
}

.user-distribution {
  display: flex;
  align-items: center;
  padding: 24px;
}

.progress-text {
  font-size: 18px;
  font-weight: bold;
}

.distribution-info {
  flex: 1;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.info-item .n-tag {
  margin-right: 12px;
  min-width: 60px;
}

.info-value {
  font-size: 16px;
  font-weight: 500;
}
</style>