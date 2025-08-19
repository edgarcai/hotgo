<template>
  <div class="search-container">
    <!-- 搜索头部 -->
    <div class="search-header">
      <div class="search-bar">
        <n-input-group>
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索漫画、作者或关键词..."
            size="large"
            @keyup.enter="performSearch"
            clearable
          >
            <template #prefix>
              <n-icon :component="SearchOutlined" />
            </template>
          </n-input>
          <n-button type="primary" size="large" @click="performSearch">
            搜索
          </n-button>
        </n-input-group>
      </div>
      
      <!-- 高级筛选 -->
      <div class="advanced-filters">
        <n-button text @click="showFilters = !showFilters">
          <n-icon :component="FilterOutlined" />
          高级筛选
          <n-icon :component="showFilters ? UpOutlined : DownOutlined" />
        </n-button>
      </div>
    </div>
    
    <!-- 筛选面板 -->
    <n-collapse-transition :show="showFilters">
      <div class="filters-panel">
        <div class="filter-row">
          <div class="filter-group">
            <label class="filter-label">分类</label>
            <n-select
              v-model:value="filters.categoryId"
              :options="categoryOptions"
              placeholder="选择分类"
              clearable
              multiple
            />
          </div>
          
          <div class="filter-group">
            <label class="filter-label">状态</label>
            <n-select
              v-model:value="filters.status"
              :options="statusOptions"
              placeholder="选择状态"
              clearable
            />
          </div>
          
          <div class="filter-group">
            <label class="filter-label">地区</label>
            <n-select
              v-model:value="filters.region"
              :options="regionOptions"
              placeholder="选择地区"
              clearable
            />
          </div>
        </div>
        
        <div class="filter-row">
          <div class="filter-group">
            <label class="filter-label">排序</label>
            <n-select
              v-model:value="filters.sortBy"
              :options="sortOptions"
              placeholder="选择排序方式"
            />
          </div>
          
          <div class="filter-group">
            <label class="filter-label">更新时间</label>
            <n-date-picker
              v-model:value="filters.dateRange"
              type="daterange"
              placeholder="选择时间范围"
              clearable
            />
          </div>
          
          <div class="filter-actions">
            <n-button @click="resetFilters">重置</n-button>
            <n-button type="primary" @click="applyFilters">应用筛选</n-button>
          </div>
        </div>
      </div>
    </n-collapse-transition>
    
    <!-- 搜索结果 -->
    <div class="search-results">
      <!-- 搜索统计 -->
      <div v-if="searchPerformed" class="search-stats">
        <span class="stats-text">
          找到 <strong>{{ totalCount }}</strong> 个结果
          <span v-if="searchQuery">关于 "{{ searchQuery }}"</span>
        </span>
        
        <div class="view-mode-toggle">
          <n-button-group>
            <n-button 
              :type="viewMode === 'grid' ? 'primary' : 'default'"
              @click="viewMode = 'grid'"
            >
              <n-icon :component="AppstoreOutlined" />
            </n-button>
            <n-button 
              :type="viewMode === 'list' ? 'primary' : 'default'"
              @click="viewMode = 'list'"
            >
              <n-icon :component="BarsOutlined" />
            </n-button>
          </n-button-group>
        </div>
      </div>
      
      <!-- 加载状态 -->
      <n-spin :show="loading">
        <!-- 网格视图 -->
        <div v-if="viewMode === 'grid'" class="results-grid">
          <div
            v-for="comic in searchResults"
            :key="comic.id"
            class="comic-card"
            @click="goToComicDetail(comic.id)"
          >
            <div class="comic-cover">
              <img :src="comic.cover" :alt="comic.title" />
              <div class="comic-overlay">
                <div class="comic-stats">
                  <span class="stat-item">
                    <n-icon :component="EyeOutlined" />
                    {{ formatNumber(comic.viewCount) }}
                  </span>
                  <span class="stat-item">
                    <n-icon :component="HeartOutlined" />
                    {{ formatNumber(comic.likeCount) }}
                  </span>
                </div>
              </div>
            </div>
            
            <div class="comic-info">
              <h3 class="comic-title">{{ comic.title }}</h3>
              <p class="comic-author">{{ comic.authorName }}</p>
              <div class="comic-meta">
                <span class="comic-status" :class="comic.status === 1 ? 'ongoing' : 'completed'">
                  {{ comic.status === 1 ? '连载中' : '已完结' }}
                </span>
                <span class="comic-chapter">{{ comic.latestChapter }}</span>
              </div>
              <p class="comic-description">{{ comic.description }}</p>
            </div>
          </div>
        </div>
        
        <!-- 列表视图 -->
        <div v-else class="results-list">
          <div
            v-for="comic in searchResults"
            :key="comic.id"
            class="comic-list-item"
            @click="goToComicDetail(comic.id)"
          >
            <div class="list-cover">
              <img :src="comic.cover" :alt="comic.title" />
            </div>
            
            <div class="list-info">
              <div class="list-header">
                <h3 class="list-title">{{ comic.title }}</h3>
                <div class="list-stats">
                  <span class="stat-item">
                    <n-icon :component="EyeOutlined" />
                    {{ formatNumber(comic.viewCount) }}
                  </span>
                  <span class="stat-item">
                    <n-icon :component="HeartOutlined" />
                    {{ formatNumber(comic.likeCount) }}
                  </span>
                </div>
              </div>
              
              <div class="list-meta">
                <span class="list-author">作者：{{ comic.authorName }}</span>
                <span class="list-status" :class="comic.status === 1 ? 'ongoing' : 'completed'">
                  {{ comic.status === 1 ? '连载中' : '已完结' }}
                </span>
                <span class="list-chapter">{{ comic.latestChapter }}</span>
                <span class="list-update">{{ formatTime(comic.updatedAt) }}</span>
              </div>
              
              <p class="list-description">{{ comic.description }}</p>
              
              <div class="list-tags">
                <n-tag v-for="category in comic.categories" :key="category.id" size="small">
                  {{ category.name }}
                </n-tag>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 空状态 -->
        <div v-if="searchPerformed && !loading && searchResults.length === 0" class="empty-state">
          <n-empty description="没有找到相关漫画">
            <template #icon>
              <n-icon :component="SearchOutlined" size="48" />
            </template>
            <template #extra>
              <n-button @click="resetSearch">重新搜索</n-button>
            </template>
          </n-empty>
        </div>
      </n-spin>
      
      <!-- 分页 -->
      <div v-if="searchResults.length > 0" class="pagination-container">
        <n-pagination
          v-model:page="currentPage"
          :page-count="totalPages"
          :page-size="pageSize"
          :item-count="totalCount"
          show-size-picker
          :page-sizes="[12, 24, 48, 96]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </div>
    
    <!-- 热门搜索 -->
    <div v-if="!searchPerformed" class="hot-searches">
      <h3>热门搜索</h3>
      <div class="hot-tags">
        <n-tag
          v-for="tag in hotSearchTags"
          :key="tag"
          class="hot-tag"
          @click="searchHotTag(tag)"
        >
          {{ tag }}
        </n-tag>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, computed, onMounted } from 'vue';
  import { useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import {
    SearchOutlined,
    FilterOutlined,
    UpOutlined,
    DownOutlined,
    AppstoreOutlined,
    BarsOutlined,
    EyeOutlined,
    HeartOutlined
  } from '@vicons/antd';
  import { searchComics, type ComicInfo } from '@/api/anix/comic';
  import { getAllCategories } from '@/api/anix/category';

  const router = useRouter();
  const message = useMessage();

  // 搜索状态
  const searchQuery = ref('');
  const loading = ref(false);
  const searchPerformed = ref(false);
  const showFilters = ref(false);
  const viewMode = ref<'grid' | 'list'>('grid');

  // 搜索结果
  const searchResults = ref<ComicInfo[]>([]);
  const totalCount = ref(0);
  const currentPage = ref(1);
  const pageSize = ref(24);
  const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

  // 筛选条件
  const filters = reactive({
    categoryId: [] as number[],
    status: null as number | null,
    region: null as string | null,
    sortBy: 'updated_at',
    dateRange: null as [number, number] | null
  });

  // 选项数据
  const categoryOptions = ref<Array<{ label: string; value: number }>>([]);
  const statusOptions = [
    { label: '连载中', value: 1 },
    { label: '已完结', value: 2 }
  ];
  const regionOptions = [
    { label: '国产', value: 'cn' },
    { label: '日本', value: 'jp' },
    { label: '韩国', value: 'kr' },
    { label: '欧美', value: 'us' }
  ];
  const sortOptions = [
    { label: '最新更新', value: 'updated_at' },
    { label: '最多阅读', value: 'view_count' },
    { label: '最多点赞', value: 'like_count' },
    { label: '最新发布', value: 'created_at' }
  ];

  // 热门搜索标签
  const hotSearchTags = ref([
    '恋爱', '校园', '冒险', '奇幻', '悬疑', '搞笑', '热血', '治愈',
    '科幻', '历史', '武侠', '都市', '玄幻', '恐怖', '励志', '青春'
  ]);

  // 格式化数字
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w';
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k';
    }
    return num.toString();
  };

  // 格式化时间
  const formatTime = (dateStr: string): string => {
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    
    if (days === 0) {
      return '今天更新';
    } else if (days === 1) {
      return '昨天更新';
    } else if (days < 7) {
      return `${days}天前更新`;
    } else {
      return date.toLocaleDateString('zh-CN');
    }
  };

  // 执行搜索
  const performSearch = async () => {
    if (!searchQuery.value.trim() && !hasActiveFilters()) {
      message.warning('请输入搜索关键词或设置筛选条件');
      return;
    }

    try {
      loading.value = true;
      searchPerformed.value = true;
      currentPage.value = 1;

      const params = {
        keyword: searchQuery.value.trim(),
        page: currentPage.value,
        pageSize: pageSize.value,
        categoryIds: filters.categoryId.length > 0 ? filters.categoryId : undefined,
        status: filters.status || undefined,
        region: filters.region || undefined,
        sortBy: filters.sortBy,
        startDate: filters.dateRange?.[0] ? new Date(filters.dateRange[0]).toISOString() : undefined,
        endDate: filters.dateRange?.[1] ? new Date(filters.dateRange[1]).toISOString() : undefined
      };

      const response = await searchComics(params);
      searchResults.value = response.data.list;
      totalCount.value = response.data.total;
    } catch (error) {
      console.error('搜索失败:', error);
      message.error('搜索失败，请稍后重试');
    } finally {
      loading.value = false;
    }
  };

  // 检查是否有活跃的筛选条件
  const hasActiveFilters = (): boolean => {
    return filters.categoryId.length > 0 ||
           filters.status !== null ||
           filters.region !== null ||
           filters.dateRange !== null;
  };

  // 应用筛选
  const applyFilters = () => {
    performSearch();
  };

  // 重置筛选
  const resetFilters = () => {
    filters.categoryId = [];
    filters.status = null;
    filters.region = null;
    filters.sortBy = 'updated_at';
    filters.dateRange = null;
  };

  // 重置搜索
  const resetSearch = () => {
    searchQuery.value = '';
    resetFilters();
    searchPerformed.value = false;
    searchResults.value = [];
    totalCount.value = 0;
  };

  // 搜索热门标签
  const searchHotTag = (tag: string) => {
    searchQuery.value = tag;
    performSearch();
  };

  // 跳转到漫画详情
  const goToComicDetail = (comicId: number) => {
    router.push(`/anix/comic/${comicId}`);
  };

  // 页码变化
  const handlePageChange = (page: number) => {
    currentPage.value = page;
    performSearch();
  };

  // 页面大小变化
  const handlePageSizeChange = (size: number) => {
    pageSize.value = size;
    currentPage.value = 1;
    performSearch();
  };

  // 加载分类选项
  const loadCategories = async () => {
    try {
      const response = await getAllCategories();
      categoryOptions.value = response.data.map(category => ({
        label: category.name,
        value: category.id
      }));
    } catch (error) {
      console.error('获取分类失败:', error);
    }
  };

  // 组件挂载
  onMounted(() => {
    loadCategories();
  });
</script>

<style scoped lang="less">
  .search-container {
    min-height: 100vh;
    background: #f5f5f5;
    padding: 24px;
  }

  .search-header {
    max-width: 1200px;
    margin: 0 auto 24px;
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .search-bar {
    margin-bottom: 16px;

    .n-input-group {
      max-width: 600px;
      margin: 0 auto;
    }
  }

  .advanced-filters {
    text-align: center;
  }

  .filters-panel {
    background: #f8f9fa;
    padding: 20px;
    border-radius: 8px;
    margin-top: 16px;
  }

  .filter-row {
    display: flex;
    gap: 16px;
    align-items: end;
    margin-bottom: 16px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .filter-group {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .filter-label {
    font-size: 14px;
    font-weight: 500;
    color: #333;
  }

  .filter-actions {
    display: flex;
    gap: 8px;
  }

  .search-results {
    max-width: 1200px;
    margin: 0 auto;
  }

  .search-stats {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding: 16px 0;
  }

  .stats-text {
    color: #666;
    font-size: 14px;
  }

  .results-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
    margin-bottom: 32px;
  }

  .comic-card {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    }
  }

  .comic-cover {
    position: relative;
    height: 200px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .comic-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(to bottom, transparent 0%, rgba(0, 0, 0, 0.7) 100%);
    display: flex;
    align-items: flex-end;
    padding: 12px;
    opacity: 0;
    transition: opacity 0.3s ease;
  }

  .comic-card:hover .comic-overlay {
    opacity: 1;
  }

  .comic-stats {
    display: flex;
    gap: 12px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 4px;
    color: white;
    font-size: 12px;
  }

  .comic-info {
    padding: 16px;
  }

  .comic-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin: 0 0 8px 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .comic-author {
    color: #666;
    font-size: 14px;
    margin: 0 0 8px 0;
  }

  .comic-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .comic-status {
    font-size: 12px;
    padding: 2px 6px;
    border-radius: 4px;

    &.ongoing {
      background: #e6f7ff;
      color: #1890ff;
    }

    &.completed {
      background: #f6ffed;
      color: #52c41a;
    }
  }

  .comic-chapter {
    font-size: 12px;
    color: #999;
  }

  .comic-description {
    color: #666;
    font-size: 13px;
    line-height: 1.4;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .results-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 32px;
  }

  .comic-list-item {
    display: flex;
    gap: 16px;
    background: white;
    padding: 16px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
    }
  }

  .list-cover {
    flex-shrink: 0;
    width: 80px;
    height: 100px;
    border-radius: 8px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .list-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .list-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
  }

  .list-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin: 0;
    line-height: 1.4;
  }

  .list-stats {
    display: flex;
    gap: 12px;
  }

  .list-meta {
    display: flex;
    gap: 16px;
    font-size: 13px;
    color: #666;
  }

  .list-status {
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 12px;

    &.ongoing {
      background: #e6f7ff;
      color: #1890ff;
    }

    &.completed {
      background: #f6ffed;
      color: #52c41a;
    }
  }

  .list-description {
    color: #666;
    font-size: 13px;
    line-height: 1.4;
    margin: 0;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .list-tags {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
  }

  .pagination-container {
    display: flex;
    justify-content: center;
    padding: 32px 0;
  }

  .hot-searches {
    max-width: 1200px;
    margin: 0 auto;
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

    h3 {
      font-size: 18px;
      font-weight: 600;
      color: #333;
      margin: 0 0 16px 0;
    }
  }

  .hot-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
  }

  .hot-tag {
    cursor: pointer;
    transition: all 0.3s ease;

    &:hover {
      background: #667eea;
      color: white;
    }
  }

  @media (max-width: 768px) {
    .search-container {
      padding: 16px;
    }

    .search-header {
      padding: 16px;
    }

    .filter-row {
      flex-direction: column;
      gap: 12px;
    }

    .filter-actions {
      justify-content: center;
    }

    .search-stats {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }

    .results-grid {
      grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
      gap: 16px;
    }

    .comic-list-item {
      padding: 12px;
    }

    .list-header {
      flex-direction: column;
      gap: 8px;
    }

    .list-meta {
      flex-wrap: wrap;
      gap: 8px;
    }
  }
</style>