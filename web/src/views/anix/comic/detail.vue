<template>
  <div class="comic-detail-container">
    <n-spin :show="loading">
      <div v-if="comic" class="comic-detail">
        <!-- 漫画基本信息 -->
        <div class="comic-header">
          <div class="comic-cover-section">
            <img :src="comic.cover" :alt="comic.title" class="comic-cover" />
          </div>
          
          <div class="comic-info-section">
            <h1 class="comic-title">{{ comic.title }}</h1>
            
            <div class="comic-meta">
              <div class="meta-item">
                <span class="meta-label">作者：</span>
                <span class="meta-value">{{ comic.authorName }}</span>
              </div>
              
              <div class="meta-item">
                <span class="meta-label">状态：</span>
                <n-tag :type="comic.status === 1 ? 'success' : 'warning'">
                  {{ comic.status === 1 ? '连载中' : '已完结' }}
                </n-tag>
              </div>
              
              <div class="meta-item">
                <span class="meta-label">分类：</span>
                <n-space>
                  <n-tag v-for="category in comic.categories" :key="category.id" size="small">
                    {{ category.name }}
                  </n-tag>
                </n-space>
              </div>
              
              <div class="meta-item">
                <span class="meta-label">标签：</span>
                <n-space>
                  <n-tag v-for="tag in tagList" :key="tag" size="small" type="info">
                    {{ tag }}
                  </n-tag>
                </n-space>
              </div>
            </div>
            
            <div class="comic-stats">
              <div class="stat-item">
                <n-icon :component="EyeOutlined" />
                <span>{{ formatNumber(comic.viewCount) }} 阅读</span>
              </div>
              <div class="stat-item">
                <n-icon :component="HeartOutlined" />
                <span>{{ formatNumber(comic.likeCount) }} 点赞</span>
              </div>
              <div class="stat-item">
                <n-icon :component="BookOutlined" />
                <span>{{ comic.chapterCount }} 章节</span>
              </div>
            </div>
            
            <div class="comic-actions">
              <n-button type="primary" size="large" @click="startReading">
                开始阅读
              </n-button>
              <n-button size="large" @click="toggleFavorite">
                <n-icon :component="isFavorited ? HeartFilled : HeartOutlined" />
                {{ isFavorited ? '已收藏' : '收藏' }}
              </n-button>
            </div>
          </div>
        </div>
        
        <!-- 漫画描述 -->
        <div class="comic-description">
          <h3>作品简介</h3>
          <p>{{ comic.description || '暂无简介' }}</p>
        </div>
        
        <!-- 章节列表 -->
        <div class="chapters-section">
          <div class="section-header">
            <h3>章节列表</h3>
            <n-space>
              <n-button-group>
                <n-button 
                  :type="sortOrder === 'asc' ? 'primary' : 'default'"
                  @click="sortOrder = 'asc'"
                >
                  正序
                </n-button>
                <n-button 
                  :type="sortOrder === 'desc' ? 'primary' : 'default'"
                  @click="sortOrder = 'desc'"
                >
                  倒序
                </n-button>
              </n-button-group>
            </n-space>
          </div>
          
          <div class="chapters-grid">
            <div
              v-for="chapter in sortedChapters"
              :key="chapter.id"
              class="chapter-item"
              :class="{ 'vip-chapter': chapter.isVip }"
              @click="readChapter(chapter.id)"
            >
              <div class="chapter-info">
                <span class="chapter-title">{{ chapter.title }}</span>
                <span class="chapter-number">第{{ chapter.chapterNumber }}话</span>
              </div>
              <div class="chapter-meta">
                <span class="chapter-pages">{{ chapter.pageCount }}页</span>
                <span class="chapter-date">{{ formatDate(chapter.publishedAt) }}</span>
                <n-tag v-if="chapter.isVip" size="small" type="warning">VIP</n-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import { EyeOutlined, HeartOutlined, HeartFilled, BookOutlined } from '@vicons/antd';
  import { getComicDetail, type ComicDetail } from '@/api/anix/comic';

  const route = useRoute();
  const router = useRouter();
  const message = useMessage();

  // 数据状态
  const loading = ref(false);
  const comic = ref<ComicDetail | null>(null);
  const isFavorited = ref(false);
  const sortOrder = ref<'asc' | 'desc'>('asc');

  // 计算属性
  const tagList = computed(() => {
    if (!comic.value?.tags) return [];
    return comic.value.tags.split(',').filter(tag => tag.trim());
  });

  const sortedChapters = computed(() => {
    if (!comic.value?.chapters) return [];
    const chapters = [...comic.value.chapters];
    return chapters.sort((a, b) => {
      if (sortOrder.value === 'asc') {
        return a.chapterNumber - b.chapterNumber;
      } else {
        return b.chapterNumber - a.chapterNumber;
      }
    });
  });

  // 格式化数字
  const formatNumber = (num: number): string => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w';
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k';
    }
    return num.toString();
  };

  // 格式化日期
  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr);
    return date.toLocaleDateString('zh-CN', {
      month: 'short',
      day: 'numeric'
    });
  };

  // 开始阅读
  const startReading = () => {
    if (comic.value?.chapters && comic.value.chapters.length > 0) {
      const firstChapter = sortedChapters.value[0];
      readChapter(firstChapter.id);
    } else {
      message.warning('暂无章节可阅读');
    }
  };

  // 阅读章节
  const readChapter = (chapterId: number) => {
    router.push(`/anix/chapter/${chapterId}`);
  };

  // 切换收藏状态
  const toggleFavorite = () => {
    // TODO: 实现收藏功能
    isFavorited.value = !isFavorited.value;
    message.success(isFavorited.value ? '已添加到收藏' : '已取消收藏');
  };

  // 获取漫画详情
  const loadComicDetail = async () => {
    try {
      loading.value = true;
      const comicId = Number(route.params.id);
      if (!comicId) {
        message.error('漫画ID无效');
        router.back();
        return;
      }

      const response = await getComicDetail(comicId);
      comic.value = response.data;
    } catch (error) {
      console.error('获取漫画详情失败:', error);
      message.error('获取漫画详情失败');
      router.back();
    } finally {
      loading.value = false;
    }
  };

  // 组件挂载时获取数据
  onMounted(() => {
    loadComicDetail();
  });
</script>

<style scoped lang="less">
  .comic-detail-container {
    min-height: 100vh;
    background: #f5f5f5;
    padding: 24px;
  }

  .comic-detail {
    max-width: 1200px;
    margin: 0 auto;
  }

  .comic-header {
    display: flex;
    gap: 32px;
    background: white;
    padding: 32px;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    margin-bottom: 24px;
  }

  .comic-cover-section {
    flex-shrink: 0;
  }

  .comic-cover {
    width: 240px;
    height: 320px;
    object-fit: cover;
    border-radius: 12px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  }

  .comic-info-section {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .comic-title {
    font-size: 32px;
    font-weight: 700;
    color: #333;
    margin: 0;
    line-height: 1.3;
  }

  .comic-meta {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .meta-item {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .meta-label {
    font-weight: 600;
    color: #666;
    min-width: 60px;
  }

  .meta-value {
    color: #333;
  }

  .comic-stats {
    display: flex;
    gap: 24px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #666;
    font-size: 14px;

    .n-icon {
      font-size: 16px;
    }
  }

  .comic-actions {
    display: flex;
    gap: 16px;
  }

  .comic-description {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    margin-bottom: 24px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      color: #333;
      margin: 0 0 16px 0;
    }

    p {
      color: #666;
      line-height: 1.6;
      margin: 0;
    }
  }

  .chapters-section {
    background: white;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      font-size: 18px;
      font-weight: 600;
      color: #333;
      margin: 0;
    }
  }

  .chapters-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 12px;
  }

  .chapter-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 2px solid transparent;

    &:hover {
      background: #e9ecef;
      border-color: #667eea;
    }

    &.vip-chapter {
      background: linear-gradient(135deg, #fff3cd 0%, #ffeaa7 100%);
      
      &:hover {
        background: linear-gradient(135deg, #ffeaa7 0%, #fdcb6e 100%);
      }
    }
  }

  .chapter-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .chapter-title {
    font-weight: 600;
    color: #333;
    font-size: 14px;
  }

  .chapter-number {
    font-size: 12px;
    color: #666;
  }

  .chapter-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 12px;
    color: #999;
  }

  .chapter-pages,
  .chapter-date {
    white-space: nowrap;
  }

  @media (max-width: 768px) {
    .comic-detail-container {
      padding: 16px;
    }

    .comic-header {
      flex-direction: column;
      padding: 20px;
      gap: 20px;
    }

    .comic-cover {
      width: 180px;
      height: 240px;
      margin: 0 auto;
    }

    .comic-title {
      font-size: 24px;
      text-align: center;
    }

    .comic-stats {
      justify-content: center;
    }

    .comic-actions {
      justify-content: center;
    }

    .chapters-grid {
      grid-template-columns: 1fr;
    }

    .chapter-item {
      padding: 12px;
    }

    .section-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }
  }
</style>