<template>
  <div class="anix-home-container">
    <!-- 顶部导航栏 -->
    <div class="home-header">
      <div class="header-content">
        <div class="logo-section">
          <h1 class="app-title">AniX</h1>
          <span class="app-subtitle">漫画阅读平台</span>
        </div>

        <div class="nav-section">
          <n-space>
            <n-button text @click="scrollToSection('hot')">热门漫画</n-button>
            <n-button text @click="scrollToSection('latest')">最新更新</n-button>
            <n-button text @click="scrollToSection('recommended')">推荐作品</n-button>
            <n-button text @click="goToSearch">搜索</n-button>
          </n-space>
        </div>

        <div class="user-section">
          <template v-if="isLoggedIn">
            <n-dropdown :options="userMenuOptions" @select="handleUserMenuSelect">
              <n-button text>
                <n-avatar :size="32" :src="userInfo?.avatar">
                  {{ userInfo?.nickname?.charAt(0) || 'U' }}
                </n-avatar>
                <span class="username">{{ userInfo?.nickname || userInfo?.username }}</span>
              </n-button>
            </n-dropdown>
          </template>
          <template v-else>
            <n-space>
              <n-button @click="goToLogin">登录</n-button>
              <n-button type="primary" @click="goToRegister">注册</n-button>
            </n-space>
          </template>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="home-content">
      <!-- 欢迎区域 -->
      <div class="welcome-section">
        <div class="welcome-content">
          <h2 class="welcome-title">欢迎来到 AniX 漫画世界</h2>
          <p class="welcome-description"> 发现精彩的漫画作品，享受沉浸式的阅读体验 </p>
          <template v-if="!isLoggedIn">
            <n-space class="welcome-actions">
              <n-button type="primary" size="large" @click="goToRegister"> 立即注册 </n-button>
              <n-button size="large" @click="goToLogin"> 已有账号？登录 </n-button>
            </n-space>
          </template>
        </div>
      </div>

      <!-- 功能介绍区域 -->
      <div class="features-section">
        <h3 class="section-title">平台特色</h3>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon"> 📚 </div>
            <h4 class="feature-title">丰富内容</h4>
            <p class="feature-description"> 海量优质漫画作品，涵盖各种题材和风格 </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon"> 📱 </div>
            <h4 class="feature-title">流畅阅读</h4>
            <p class="feature-description"> 优化的阅读体验，支持多种阅读模式 </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon"> 💾 </div>
            <h4 class="feature-title">个人收藏</h4>
            <p class="feature-description"> 收藏喜爱的作品，记录阅读进度 </p>
          </div>

          <div class="feature-card">
            <div class="feature-icon"> 🔄 </div>
            <h4 class="feature-title">实时更新</h4>
            <p class="feature-description"> 第一时间获取最新章节更新 </p>
          </div>
        </div>
      </div>

      <!-- 热门漫画区域 -->
      <div id="hot" class="comics-section">
        <div class="section-header">
          <h3 class="section-title">🔥 热门漫画</h3>
          <n-button text type="primary" @click="viewMore('hot')">查看更多</n-button>
        </div>
        <n-spin :show="loading.hot">
          <div class="comics-grid">
            <div
              v-for="comic in hotComics"
              :key="comic.id"
              class="comic-card"
              @click="goToComicDetail(comic.id)"
            >
              <div class="comic-cover">
                <img :src="comic.cover" :alt="comic.title" loading="lazy" />
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
                <h4 class="comic-title">{{ comic.title }}</h4>
                <p class="comic-author">{{ comic.authorName }}</p>
                <p class="comic-chapter">{{ comic.lastChapterTitle }}</p>
              </div>
            </div>
          </div>
        </n-spin>
      </div>

      <!-- 最新更新区域 -->
      <div id="latest" class="comics-section">
        <div class="section-header">
          <h3 class="section-title">🆕 最新更新</h3>
          <n-button text type="primary" @click="viewMore('latest')">查看更多</n-button>
        </div>
        <n-spin :show="loading.latest">
          <div class="comics-grid">
            <div
              v-for="comic in latestComics"
              :key="comic.id"
              class="comic-card"
              @click="goToComicDetail(comic.id)"
            >
              <div class="comic-cover">
                <img :src="comic.cover" :alt="comic.title" loading="lazy" />
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
                <h4 class="comic-title">{{ comic.title }}</h4>
                <p class="comic-author">{{ comic.authorName }}</p>
                <p class="comic-chapter">{{ comic.lastChapterTitle }}</p>
                <p class="comic-update-time">{{ formatTime(comic.updatedAt) }}</p>
              </div>
            </div>
          </div>
        </n-spin>
      </div>

      <!-- 推荐作品区域 -->
      <div id="recommended" class="comics-section">
        <div class="section-header">
          <h3 class="section-title">⭐ 推荐作品</h3>
          <n-button text type="primary" @click="viewMore('recommended')">查看更多</n-button>
        </div>
        <n-spin :show="loading.recommended">
          <div class="comics-grid">
            <div
              v-for="comic in recommendedComics"
              :key="comic.id"
              class="comic-card"
              @click="goToComicDetail(comic.id)"
            >
              <div class="comic-cover">
                <img :src="comic.cover" :alt="comic.title" loading="lazy" />
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
                <h4 class="comic-title">{{ comic.title }}</h4>
                <p class="comic-author">{{ comic.authorName }}</p>
                <p class="comic-description">{{ comic.description }}</p>
              </div>
            </div>
          </div>
        </n-spin>
      </div>

      <!-- 用户状态区域 -->
      <template v-if="isLoggedIn">
        <div class="user-status-section">
          <h3 class="section-title">我的阅读</h3>
          <div class="status-cards">
            <div class="status-card">
              <div class="status-number">0</div>
              <div class="status-label">收藏作品</div>
            </div>
            <div class="status-card">
              <div class="status-number">0</div>
              <div class="status-label">阅读历史</div>
            </div>
            <div class="status-card">
              <div class="status-number">0</div>
              <div class="status-label">关注更新</div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import { EyeOutlined, HeartOutlined } from '@vicons/antd';
  import { useAnixUserStore } from '@/store/modules/anixUser';
  import { getHotComics, getLatestComics, getRecommendedComics, type ComicInfo } from '@/api/anix/comic';

  const router = useRouter();
  const message = useMessage();
  const anixUserStore = useAnixUserStore();

  // 计算属性
  const isLoggedIn = computed(() => anixUserStore.isLoggedIn);
  const userInfo = computed(() => anixUserStore.userInfo);

  // 漫画数据
  const hotComics = ref<ComicInfo[]>([]);
  const latestComics = ref<ComicInfo[]>([]);
  const recommendedComics = ref<ComicInfo[]>([]);

  // 加载状态
  const loading = reactive({
    hot: false,
    latest: false,
    recommended: false,
  });

  // 用户菜单选项
  const userMenuOptions = [
    {
      label: '个人中心',
      key: 'profile',
    },
    {
      label: '退出登录',
      key: 'logout',
    },
  ];

  // 处理用户菜单选择
  const handleUserMenuSelect = async (key: string) => {
    switch (key) {
      case 'profile':
        router.push('/anix/profile');
        break;
      case 'logout':
        try {
          await anixUserStore.logout();
          message.success('已退出登录');
        } catch (error: any) {
          console.error('退出登录失败:', error);
          message.error('退出登录失败');
        }
        break;
    }
  };

  // 跳转到登录页
  const goToLogin = () => {
    router.push('/anix/login');
  };

  // 跳转到注册页
  const goToRegister = () => {
    router.push('/anix/register');
  };

  // 跳转到搜索页
  const goToSearch = () => {
    router.push('/anix/search');
  };

  // 跳转到漫画详情页
  const goToComicDetail = (id: number) => {
    router.push(`/anix/comic/${id}`);
  };

  // 滚动到指定区域
  const scrollToSection = (sectionId: string) => {
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  // 查看更多
  const viewMore = (type: string) => {
    router.push(`/anix/comics?type=${type}`);
  };

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
  const formatTime = (timeStr: string): string => {
    const time = new Date(timeStr);
    const now = new Date();
    const diff = now.getTime() - time.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    
    if (days === 0) {
      const hours = Math.floor(diff / (1000 * 60 * 60));
      if (hours === 0) {
        const minutes = Math.floor(diff / (1000 * 60));
        return `${minutes}分钟前`;
      }
      return `${hours}小时前`;
    } else if (days < 7) {
      return `${days}天前`;
    } else {
      return time.toLocaleDateString();
    }
  };

  // 获取热门漫画
  const loadHotComics = async () => {
    try {
      loading.hot = true;
      const response = await getHotComics(8);
      hotComics.value = response.data?.list || response.data || [];
    } catch (error) {
      console.error('获取热门漫画失败:', error);
      message.error('获取热门漫画失败');
    } finally {
      loading.hot = false;
    }
  };

  // 获取最新更新漫画
  const loadLatestComics = async () => {
    try {
      loading.latest = true;
      const response = await getLatestComics(8);
      latestComics.value = response.data?.list || response.data || [];
    } catch (error) {
      console.error('获取最新漫画失败:', error);
      message.error('获取最新漫画失败');
    } finally {
      loading.latest = false;
    }
  };

  // 获取推荐漫画
  const loadRecommendedComics = async () => {
    try {
      loading.recommended = true;
      const response = await getRecommendedComics(8);
      recommendedComics.value = response.data?.list || response.data || [];
    } catch (error) {
      console.error('获取推荐漫画失败:', error);
      message.error('获取推荐漫画失败');
    } finally {
      loading.recommended = false;
    }
  };

  // 组件挂载时获取数据
  onMounted(async () => {
    // 获取用户信息
    if (isLoggedIn.value) {
      try {
        await anixUserStore.getUserInfo();
      } catch (error) {
        console.error('获取用户信息失败:', error);
      }
    }

    // 并行获取漫画数据
    await Promise.all([
      loadHotComics(),
      loadLatestComics(),
      loadRecommendedComics(),
    ]);
  });
</script>

<style scoped lang="less">
  .anix-home-container {
    min-height: 100vh;
    background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  }

  .home-header {
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .header-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 16px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 24px;
  }

  .nav-section {
    flex: 1;
    display: flex;
    justify-content: center;
  }

  .logo-section {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .app-title {
    font-size: 28px;
    font-weight: 700;
    color: #333;
    margin: 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .app-subtitle {
    font-size: 14px;
    color: #666;
  }

  .user-section {
    .username {
      margin-left: 8px;
      color: #333;
    }
  }

  .home-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 40px 24px;
  }

  .welcome-section {
    text-align: center;
    padding: 60px 0;
    background: white;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    margin-bottom: 40px;
  }

  .welcome-title {
    font-size: 36px;
    font-weight: 600;
    color: #333;
    margin: 0 0 16px 0;
  }

  .welcome-description {
    font-size: 18px;
    color: #666;
    margin: 0 0 32px 0;
    line-height: 1.6;
  }

  .welcome-actions {
    justify-content: center;
  }

  .section-title {
    font-size: 24px;
    font-weight: 600;
    color: #333;
    text-align: center;
    margin: 0 0 32px 0;
  }

  .features-section {
    margin-bottom: 40px;
  }

  .features-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 24px;
  }

  .feature-card {
    background: white;
    padding: 32px 24px;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    text-align: center;
    transition:
      transform 0.3s ease,
      box-shadow 0.3s ease;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    }
  }

  .feature-icon {
    font-size: 48px;
    margin-bottom: 16px;
  }

  .feature-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
    margin: 0 0 12px 0;
  }

  .feature-description {
    font-size: 14px;
    color: #666;
    line-height: 1.6;
    margin: 0;
  }

  .user-status-section {
    background: white;
    padding: 32px;
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }

  .status-cards {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 24px;
  }

  .status-card {
    text-align: center;
    padding: 24px;
    background: #f8f9fa;
    border-radius: 12px;
    border: 2px solid transparent;
    transition: border-color 0.3s ease;

    &:hover {
      border-color: #667eea;
    }
  }

  .status-number {
    font-size: 32px;
    font-weight: 700;
    color: #667eea;
    margin-bottom: 8px;
  }

  .status-label {
    font-size: 14px;
    color: #666;
  }

  // 漫画区域样式
  .comics-section {
    margin-bottom: 48px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
  }

  .comics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
  }

  .comic-card {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    cursor: pointer;

    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    }
  }

  .comic-cover {
    position: relative;
    width: 100%;
    height: 280px;
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform 0.3s ease;
    }

    &:hover img {
      transform: scale(1.05);
    }
  }

  .comic-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(to bottom, transparent 0%, rgba(0, 0, 0, 0.7) 100%);
    opacity: 0;
    transition: opacity 0.3s ease;
    display: flex;
    align-items: flex-end;
    padding: 16px;

    .comic-card:hover & {
      opacity: 1;
    }
  }

  .comic-stats {
    display: flex;
    gap: 16px;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 4px;
    color: white;
    font-size: 12px;

    .n-icon {
      font-size: 14px;
    }
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
    font-size: 14px;
    color: #666;
    margin: 0 0 4px 0;
  }

  .comic-chapter {
    font-size: 12px;
    color: #999;
    margin: 0 0 4px 0;
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .comic-update-time {
    font-size: 12px;
    color: #999;
    margin: 0;
  }

  .comic-description {
    font-size: 12px;
    color: #666;
    margin: 0;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  @media (max-width: 768px) {
    .header-content {
      padding: 12px 16px;
    }

    .app-title {
      font-size: 24px;
    }

    .home-content {
      padding: 24px 16px;
    }

    .welcome-title {
      font-size: 28px;
    }

    .welcome-description {
      font-size: 16px;
    }

    .features-grid {
      grid-template-columns: 1fr;
    }

    .status-cards {
      grid-template-columns: 1fr;
    }

    .nav-section {
      display: none;
    }

    .comics-grid {
      grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
      gap: 16px;
    }

    .comic-cover {
      height: 200px;
    }

    .comic-info {
      padding: 12px;
    }

    .comic-title {
      font-size: 14px;
    }

    .section-header {
      margin-bottom: 16px;
    }

    .comics-section {
      margin-bottom: 32px;
    }
  }
</style>
