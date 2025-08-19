<template>
  <div class="chapter-reader-container">
    <n-spin :show="loading">
      <div v-if="chapter" class="chapter-reader">
        <!-- 阅读器头部 -->
        <div class="reader-header" :class="{ 'hidden': isHeaderHidden }">
          <div class="header-left">
            <n-button text @click="goBack">
              <n-icon :component="ArrowLeftOutlined" />
              返回
            </n-button>
            <div class="comic-info">
              <span class="comic-title">{{ chapter.comicTitle }}</span>
              <span class="chapter-title">{{ chapter.title }}</span>
            </div>
          </div>
          
          <div class="header-right">
            <n-button-group>
              <n-button @click="toggleReadingMode">
                <n-icon :component="readingMode === 'vertical' ? ViewListOutlined : readingMode === 'horizontal' ? ViewModuleOutlined : BookOutlined" />
                {{ readingMode === 'vertical' ? '横向' : readingMode === 'horizontal' ? '双页' : '纵向' }}
              </n-button>
              <n-button @click="toggleSettings">
                <n-icon :component="SettingOutlined" />
                设置
              </n-button>
              <n-button @click="toggleFullscreen">
                <n-icon :component="isFullscreen ? FullscreenExitOutlined : FullscreenOutlined" />
                {{ isFullscreen ? '退出全屏' : '全屏' }}
              </n-button>
            </n-button-group>
          </div>
        </div>
        
        <!-- 章节导航 -->
        <div class="chapter-navigation" :class="{ 'hidden': isHeaderHidden }">
          <n-button 
            :disabled="!prevChapter" 
            @click="goToChapter(prevChapter?.id)"
            size="large"
          >
            <n-icon :component="LeftOutlined" />
            上一话
          </n-button>
          
          <n-select
            v-model:value="currentChapterId"
            :options="chapterOptions"
            @update:value="goToChapter"
            style="width: 200px;"
          />
          
          <n-button 
            :disabled="!nextChapter" 
            @click="goToChapter(nextChapter?.id)"
            size="large"
          >
            下一话
            <n-icon :component="RightOutlined" />
          </n-button>
        </div>
        
        <!-- 阅读器设置面板 -->
        <div class="reader-settings" :class="{ 'show': showSettings }">
          <div class="settings-content">
            <div class="setting-group">
              <label>阅读模式</label>
              <n-radio-group v-model:value="readingMode" @update:value="onReadingModeChange">
                <n-radio value="vertical">纵向滚动</n-radio>
                <n-radio value="horizontal">横向翻页</n-radio>
                <n-radio value="double">双页模式</n-radio>
              </n-radio-group>
            </div>
            
            <div class="setting-group">
              <label>夜间模式</label>
              <n-switch v-model:value="isDarkMode" @update:value="toggleDarkMode" />
            </div>
            
            <div class="setting-group">
              <label>亮度调节</label>
              <n-slider v-model:value="brightness" :min="0.3" :max="1" :step="0.1" @update:value="updateBrightness" />
            </div>
            
            <div class="setting-group">
              <label>缩放比例</label>
              <n-slider v-model:value="zoomLevel" :min="0.5" :max="3" :step="0.1" @update:value="updateZoom" />
            </div>
            
            <div class="setting-group">
              <label>自动隐藏工具栏</label>
              <n-switch v-model:value="autoHideToolbar" />
            </div>
            
            <div class="setting-group">
              <label>阅读统计</label>
              <div class="reading-stats">
                <p>当前页: {{ currentPageIndex + 1 }} / {{ chapter.pages.length }}</p>
                <p>阅读进度: {{ Math.round(readingProgress) }}%</p>
                <p>阅读时长: {{ Math.round(totalReadingTime / 60) }} 分钟</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 书签对话框 -->
        <n-modal v-model:show="showBookmarkDialog" preset="dialog" title="添加书签">
          <template #default>
            <div class="bookmark-form">
              <p class="bookmark-info">
                章节: {{ chapter.title }}<br>
                页面: 第 {{ currentPageIndex + 1 }} 页
              </p>
              <n-input
                v-model:value="bookmarkNote"
                type="textarea"
                placeholder="添加备注（可选）"
                :rows="3"
                maxlength="200"
                show-count
              />
            </div>
          </template>
          <template #action>
            <n-space>
              <n-button @click="cancelAddBookmark">取消</n-button>
              <n-button type="primary" @click="confirmAddBookmark">确认</n-button>
            </n-space>
          </template>
        </n-modal>
        
        <!-- 书签列表抽屉 -->
        <n-drawer v-model:show="showBookmarkList" :width="320" placement="right">
          <n-drawer-content title="书签列表">
            <div class="bookmark-list">
              <div v-if="bookmarks.length === 0" class="empty-bookmarks">
                <n-icon size="48" color="#ccc">
                  <BookmarkOutlined />
                </n-icon>
                <p>暂无书签</p>
              </div>
              <div v-else>
                <div 
                  v-for="bookmark in bookmarks" 
                  :key="bookmark.id"
                  class="bookmark-item"
                  @click="goToBookmark(bookmark)"
                >
                  <div class="bookmark-info">
                    <h4>{{ bookmark.title }}</h4>
                    <p v-if="bookmark.note" class="bookmark-note">{{ bookmark.note }}</p>
                    <p class="bookmark-time">{{ new Date(bookmark.timestamp).toLocaleString() }}</p>
                  </div>
                  <n-button 
                    size="small" 
                    type="error" 
                    ghost
                    @click.stop="removeBookmark(bookmark.id)"
                  >
                    删除
                  </n-button>
                </div>
              </div>
            </div>
          </n-drawer-content>
        </n-drawer>
        
        <!-- 阅读内容区域 -->
        <div 
          class="reading-area" 
          :class="{
            'vertical-mode': readingMode === 'vertical',
            'horizontal-mode': readingMode === 'horizontal',
            'double-mode': readingMode === 'double',
            'fullscreen': isFullscreen,
            'dark-mode': isDarkMode
          }"
          :style="{ filter: `brightness(${brightness})` }"
          @click="handleAreaClick"
          @touchstart="handleTouchStart"
          @touchmove="handleTouchMove"
          @touchend="handleTouchEnd"
          @wheel="handleWheel"
        >
          <!-- 纵向阅读模式 -->
          <div v-if="readingMode === 'vertical'" class="vertical-reader" ref="verticalReader">
            <div 
              v-for="(page, index) in chapter.pages" 
              :key="index"
              class="page-container"
              :data-page="index + 1"
              @click="handleImageClick($event, index)"
            >
              <div class="image-wrapper" :style="{ transform: `scale(${zoomLevel})` }">
                <img 
                  :src="page.imageUrl" 
                  :alt="`第${index + 1}页`"
                  class="page-image"
                  @load="onImageLoad(index)"
                  @error="onImageError(index)"
                  draggable="false"
                />
              </div>
              <div class="page-number">{{ index + 1 }} / {{ chapter.pages.length }}</div>
            </div>
          </div>
          
          <!-- 横向阅读模式 -->
          <div v-else-if="readingMode === 'horizontal'" class="horizontal-reader">
            <div class="page-slider">
              <div 
                class="page-track"
                :style="{ transform: `translateX(-${currentPageIndex * 100}%)` }"
              >
                <div 
                  v-for="(page, index) in chapter.pages" 
                  :key="index"
                  class="page-slide"
                  @click="handleImageClick($event, index)"
                >
                  <div class="image-wrapper" :style="{ transform: `scale(${zoomLevel})` }">
                    <img 
                      :src="page.imageUrl" 
                      :alt="`第${index + 1}页`"
                      class="page-image"
                      @load="onImageLoad(index)"
                      @error="onImageError(index)"
                      draggable="false"
                    />
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 横向导航控制 -->
            <div class="horizontal-controls" :class="{ 'hidden': isHeaderHidden }">
              <n-button 
                circle 
                size="large" 
                :disabled="currentPageIndex === 0"
                @click="prevPage"
                class="nav-button prev-button"
              >
                <n-icon :component="LeftOutlined" />
              </n-button>
              
              <div class="page-indicator">
                <span>{{ currentPageIndex + 1 }} / {{ chapter.pages.length }}</span>
              </div>
              
              <n-button 
                circle 
                size="large" 
                :disabled="currentPageIndex === chapter.pages.length - 1"
                @click="nextPage"
                class="nav-button next-button"
              >
                <n-icon :component="RightOutlined" />
              </n-button>
            </div>
          </div>
          
          <!-- 双页阅读模式 -->
          <div v-else class="double-reader">
            <div class="double-page-container">
              <div class="double-pages" :style="{ transform: `translateX(-${Math.floor(currentPageIndex / 2) * 100}%)` }">
                <div 
                  v-for="(pageGroup, groupIndex) in doublePageGroups" 
                  :key="groupIndex"
                  class="double-page-group"
                >
                  <div 
                    v-for="(page, pageIndex) in pageGroup" 
                    :key="page.index"
                    class="double-page-item"
                    @click="handleImageClick($event, page.index)"
                  >
                    <div class="image-wrapper" :style="{ transform: `scale(${zoomLevel})` }">
                      <img 
                        :src="page.imageUrl" 
                        :alt="`第${page.index + 1}页`"
                        class="page-image"
                        @load="onImageLoad(page.index)"
                        @error="onImageError(page.index)"
                        draggable="false"
                      />
                    </div>
                    <div class="page-number">{{ page.index + 1 }}</div>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 双页导航控制 -->
            <div class="double-controls" :class="{ 'hidden': isHeaderHidden }">
              <n-button 
                circle 
                size="large" 
                :disabled="currentPageIndex === 0"
                @click="prevDoublePage"
                class="nav-button prev-button"
              >
                <n-icon :component="LeftOutlined" />
              </n-button>
              
              <div class="page-indicator">
                <span>{{ Math.floor(currentPageIndex / 2) + 1 }} / {{ Math.ceil(chapter.pages.length / 2) }}</span>
              </div>
              
              <n-button 
                circle 
                size="large" 
                :disabled="currentPageIndex >= chapter.pages.length - 2"
                @click="nextDoublePage"
                class="nav-button next-button"
              >
                <n-icon :component="RightOutlined" />
              </n-button>
            </div>
          </div>
        </div>
        
        <!-- 底部操作栏 -->
        <div class="reader-footer" :class="{ 'hidden': isHeaderHidden }">
          <div class="footer-left">
            <n-button @click="toggleFavorite">
              <n-icon :component="isFavorited ? HeartFilled : HeartOutlined" />
              {{ isFavorited ? '已收藏' : '收藏' }}
            </n-button>
            <n-button @click="addBookmark">
              <n-icon :component="BookmarkOutlined" />
              书签
            </n-button>
            <n-button @click="showBookmarkList = true">
              <n-icon :component="HistoryOutlined" />
              书签列表
            </n-button>
          </div>
          
          <div class="footer-center">
            <div class="zoom-controls">
              <n-button size="small" @click="zoomOut" :disabled="zoomLevel <= 0.5">
                <n-icon :component="ZoomOutOutlined" />
              </n-button>
              <span class="zoom-text">{{ Math.round(zoomLevel * 100) }}%</span>
              <n-button size="small" @click="zoomIn" :disabled="zoomLevel >= 3">
                <n-icon :component="ZoomInOutlined" />
              </n-button>
            </div>
            <n-progress 
              :percentage="readingProgress" 
              :show-indicator="false"
              style="width: 150px;"
            />
            <span class="progress-text">{{ Math.round(readingProgress) }}%</span>
          </div>
          
          <div class="footer-right">
            <n-button @click="resetZoom">
              <n-icon :component="ExpandOutlined" />
              重置
            </n-button>
            <n-button @click="shareChapter">
              <n-icon :component="ShareAltOutlined" />
              分享
            </n-button>
          </div>
        </div>
        
        <!-- 触摸提示 -->
        <div class="touch-hint" v-if="showTouchHint">
          <div class="hint-content">
            <p>触摸操作提示：</p>
            <p>• 单击：显示/隐藏工具栏</p>
            <p>• 双击：缩放图片</p>
            <p>• 左右滑动：翻页</p>
            <p>• 上下滑动：滚动（纵向模式）</p>
          </div>
        </div>
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
  import { useRoute, useRouter } from 'vue-router';
  import { useMessage } from 'naive-ui';
  import {
    ArrowLeftOutlined,
    LeftOutlined,
    RightOutlined,
    ViewListOutlined,
    ViewModuleOutlined,
    BookOutlined,
    SettingOutlined,
    FullscreenOutlined,
    FullscreenExitOutlined,
    HeartOutlined,
    HeartFilled,
    ShareAltOutlined,
    BookmarkOutlined,
    ZoomInOutlined,
    ZoomOutOutlined,
    ExpandOutlined,
    HistoryOutlined
  } from '@vicons/antd';
  import { readingProgressManager, type ReadingProgress, type Bookmark } from '@/utils/readingProgress';

  const route = useRoute();
  const router = useRouter();
  const message = useMessage();

  // 数据状态
  const loading = ref(false);
  const chapter = ref<any>(null);
  const currentChapterId = ref<number>(0);
  const currentPageIndex = ref(0);
  const readingMode = ref<'vertical' | 'horizontal' | 'double'>('vertical');
  const isFullscreen = ref(false);
  const isFavorited = ref(false);
  const chapterList = ref<any[]>([]);
  
  // 新增状态
  const isHeaderHidden = ref(false);
  const showSettings = ref(false);
  const isDarkMode = ref(false);
  const brightness = ref(1);
  const zoomLevel = ref(1);
  const autoHideToolbar = ref(false);
  const showTouchHint = ref(false);
  const verticalReader = ref<HTMLElement>();
  
  // 触摸相关状态
  const touchStartX = ref(0);
  const touchStartY = ref(0);
  const touchStartTime = ref(0);
  const lastTouchTime = ref(0);
  const hideToolbarTimer = ref<NodeJS.Timeout>();
  
  // 阅读进度相关状态
  const readingStartTime = ref(0);
  const totalReadingTime = ref(0);
  const bookmarks = ref<Bookmark[]>([]);
  const showBookmarkDialog = ref(false);
  const showBookmarkList = ref(false);
  const bookmarkNote = ref('');
  const autoSaveTimer = ref<NodeJS.Timeout>();

  // 计算属性
  const chapterOptions = computed(() => {
    return chapterList.value.map(ch => ({
      label: `第${ch.chapterNumber}话 ${ch.title}`,
      value: ch.id
    }));
  });

  const prevChapter = computed(() => {
    const currentIndex = chapterList.value.findIndex(ch => ch.id === currentChapterId.value);
    return currentIndex > 0 ? chapterList.value[currentIndex - 1] : null;
  });

  const nextChapter = computed(() => {
    const currentIndex = chapterList.value.findIndex(ch => ch.id === currentChapterId.value);
    return currentIndex < chapterList.value.length - 1 ? chapterList.value[currentIndex + 1] : null;
  });

  const readingProgress = computed(() => {
    if (!chapter.value?.pages.length) return 0;
    if (readingMode.value === 'vertical') {
      // 纵向模式基于滚动位置计算
      return Math.min(100, (currentPageIndex.value / chapter.value.pages.length) * 100);
    } else if (readingMode.value === 'double') {
      // 双页模式基于双页组计算
      const totalDoublePages = Math.ceil(chapter.value.pages.length / 2);
      const currentDoublePage = Math.floor(currentPageIndex.value / 2);
      return ((currentDoublePage + 1) / totalDoublePages) * 100;
    } else {
      // 横向模式基于当前页面
      return ((currentPageIndex.value + 1) / chapter.value.pages.length) * 100;
    }
  });

  // 双页模式相关计算属性
  const doublePageGroups = computed(() => {
    if (!chapter.value?.pages?.length) return [];
    const groups = [];
    for (let i = 0; i < chapter.value.pages.length; i += 2) {
      const group = [{ ...chapter.value.pages[i], index: i }];
      if (i + 1 < chapter.value.pages.length) {
        group.push({ ...chapter.value.pages[i + 1], index: i + 1 });
      }
      groups.push(group);
    }
    return groups;
  });

  const currentDoublePageIndex = computed(() => {
    return Math.floor(currentPageIndex.value / 2);
  });

  // 返回上一页
  const goBack = () => {
    router.back();
  };

  // 切换阅读模式
  const toggleReadingMode = () => {
    const modes = ['vertical', 'horizontal', 'double'];
    const currentIndex = modes.indexOf(readingMode.value);
    readingMode.value = modes[(currentIndex + 1) % modes.length] as 'vertical' | 'horizontal' | 'double';
    currentPageIndex.value = 0;
  };

  // 切换全屏
  const toggleFullscreen = () => {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen();
      isFullscreen.value = true;
    } else {
      document.exitFullscreen();
      isFullscreen.value = false;
    }
  };
  
  // 切换设置面板
  const toggleSettings = () => {
    showSettings.value = !showSettings.value;
  };
  
  // 切换夜间模式
  const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value;
  };
  
  // 调节亮度
  const updateBrightness = (value: number) => {
    brightness.value = value;
  };
  
  // 更新缩放
  const updateZoom = (value: number) => {
    zoomLevel.value = value;
  };
  
  // 阅读模式变化处理
  const onReadingModeChange = (mode: 'vertical' | 'horizontal' | 'double') => {
    readingMode.value = mode;
    currentPageIndex.value = 0;
  };
  
  // 区域点击处理
  const handleAreaClick = (e: MouseEvent) => {
    if (autoHideToolbar.value) {
      toggleToolbar();
    }
  };
  
  // 缩放功能
  const zoomIn = () => {
    if (zoomLevel.value < 3) {
      zoomLevel.value = Math.min(3, zoomLevel.value + 0.25);
    }
  };
  
  const zoomOut = () => {
    if (zoomLevel.value > 0.5) {
      zoomLevel.value = Math.max(0.5, zoomLevel.value - 0.25);
    }
  };
  
  const resetZoom = () => {
    zoomLevel.value = 1;
  };
  
  // 添加书签
  const addBookmark = () => {
    showBookmarkDialog.value = true;
  };
  
  // 确认添加书签
  const confirmAddBookmark = () => {
    try {
      const bookmarkId = readingProgressManager.addBookmark({
        comicId: parseInt(route.params.comicId as string),
        chapterId: parseInt(route.params.chapterId as string),
        pageIndex: currentPageIndex.value,
        title: `${chapter.value.title} - 第${currentPageIndex.value + 1}页`,
        note: bookmarkNote.value || undefined
      });
      
      loadBookmarks();
      showBookmarkDialog.value = false;
      bookmarkNote.value = '';
      message.success('书签添加成功');
    } catch (error) {
      message.error('书签添加失败');
    }
  };
  
  // 取消添加书签
  const cancelAddBookmark = () => {
    showBookmarkDialog.value = false;
    bookmarkNote.value = '';
  };
  
  // 删除书签
  const removeBookmark = (bookmarkId: string) => {
    if (readingProgressManager.removeBookmark(bookmarkId)) {
      loadBookmarks();
      message.success('书签删除成功');
    } else {
      message.error('书签删除失败');
    }
  };
  
  // 跳转到书签位置
  const goToBookmark = (bookmark: Bookmark) => {
    if (bookmark.chapterId !== parseInt(route.params.id as string)) {
      // 跳转到其他章节
      router.push(`/anix/chapter/${bookmark.chapterId}`);
    } else {
      // 当前章节，直接跳转页面
      currentPageIndex.value = bookmark.pageIndex;
    }
  };
  
  // 加载书签列表
  const loadBookmarks = () => {
    const comicId = chapter.value?.comicId;
    if (comicId) {
      bookmarks.value = readingProgressManager.getBookmarks(comicId);
    }
  };

  // 监听全屏状态变化
  const handleFullscreenChange = () => {
    isFullscreen.value = !!document.fullscreenElement;
  };

  // 跳转到指定章节
  const goToChapter = (chapterId: number) => {
    if (chapterId && chapterId !== currentChapterId.value) {
      // 保存当前进度
      saveReadingProgress();
      router.push(`/anix/chapter/${chapterId}`);
    }
  };

  // 上一页
  const prevPage = () => {
    if (readingMode.value === 'double') {
      prevDoublePage();
    } else if (currentPageIndex.value > 0) {
      currentPageIndex.value--;
      // 自动保存进度
      saveReadingProgress();
    }
  };
  
  // 下一页
  const nextPage = () => {
    if (readingMode.value === 'double') {
      nextDoublePage();
    } else if (chapter.value?.pages && currentPageIndex.value < chapter.value.pages.length - 1) {
      currentPageIndex.value++;
      // 自动保存进度
      saveReadingProgress();
    }
  };
  
  // 双页模式翻页
  const prevDoublePage = () => {
    if (currentDoublePageIndex.value > 0) {
      currentPageIndex.value = Math.max(0, currentPageIndex.value - 2);
      // 自动保存进度
      saveReadingProgress();
    }
  };
  
  const nextDoublePage = () => {
    if (chapter.value?.pages) {
      const nextIndex = currentPageIndex.value + 2;
      if (nextIndex < chapter.value.pages.length) {
        currentPageIndex.value = nextIndex;
        // 自动保存进度
        saveReadingProgress();
      }
    }
  };

  // 图片加载成功
  const onImageLoad = (index: number) => {
    console.log(`页面 ${index + 1} 加载完成`);
  };

  // 图片加载失败
  const onImageError = (index: number) => {
    console.error(`页面 ${index + 1} 加载失败`);
    message.error(`第${index + 1}页加载失败`);
  };

  // 切换收藏状态
  const toggleFavorite = () => {
    isFavorited.value = !isFavorited.value;
    message.success(isFavorited.value ? '已添加到收藏' : '已取消收藏');
  };

  // 分享章节
  const shareChapter = () => {
    if (navigator.share) {
      navigator.share({
        title: `${chapter.value.comicTitle} - ${chapter.value.title}`,
        url: window.location.href
      });
    } else {
      // 复制链接到剪贴板
      navigator.clipboard.writeText(window.location.href);
      message.success('链接已复制到剪贴板');
    }
  };
  
  // 触摸手势处理
  const handleTouchStart = (e: TouchEvent) => {
    const touch = e.touches[0];
    touchStartX.value = touch.clientX;
    touchStartY.value = touch.clientY;
    touchStartTime.value = Date.now();
  };
  
  const handleTouchMove = (e: TouchEvent) => {
    // 阻止默认滚动行为（在需要时）
    if (readingMode.value === 'horizontal' || readingMode.value === 'double') {
      e.preventDefault();
    }
  };
  
  const handleTouchEnd = (e: TouchEvent) => {
    const touch = e.changedTouches[0];
    const deltaX = touch.clientX - touchStartX.value;
    const deltaY = touch.clientY - touchStartY.value;
    const deltaTime = Date.now() - touchStartTime.value;
    
    // 检测双击
    if (deltaTime < 300 && Math.abs(deltaX) < 10 && Math.abs(deltaY) < 10) {
      const now = Date.now();
      if (now - lastTouchTime.value < 300) {
        // 双击缩放
        if (zoomLevel.value === 1) {
          zoomLevel.value = 2;
        } else {
          zoomLevel.value = 1;
        }
        return;
      }
      lastTouchTime.value = now;
    }
    
    // 检测滑动手势
    if (deltaTime < 500 && Math.abs(deltaX) > 50) {
      if (deltaX > 0) {
        // 右滑 - 上一页
        prevPage();
      } else {
        // 左滑 - 下一页
        nextPage();
      }
    }
  };
  
  // 鼠标滚轮缩放
  const handleWheel = (e: WheelEvent) => {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault();
      if (e.deltaY < 0) {
        zoomIn();
      } else {
        zoomOut();
      }
    }
  };
  
  // 图片点击处理
  const handleImageClick = () => {
    if (autoHideToolbar.value) {
      toggleToolbar();
    }
  };
  
  // 切换工具栏显示/隐藏
  const toggleToolbar = () => {
    isHeaderHidden.value = !isHeaderHidden.value;
    
    if (autoHideToolbar.value && !isHeaderHidden.value) {
      // 自动隐藏工具栏
      clearTimeout(hideToolbarTimer.value);
      hideToolbarTimer.value = setTimeout(() => {
        isHeaderHidden.value = true;
      }, 3000);
    }
  };
  
  // 显示触摸提示
  const showTouchHintDialog = () => {
    showTouchHint.value = true;
    setTimeout(() => {
      showTouchHint.value = false;
    }, 5000);
  };
  
  // 保存阅读进度
  const saveReadingProgress = () => {
    const chapterId = parseInt(route.params.id as string);
    const comicId = chapter.value?.comicId;
    
    if (isNaN(chapterId) || !comicId) return;
    
    const currentTime = Date.now();
    const sessionTime = readingStartTime.value > 0 ? (currentTime - readingStartTime.value) / 1000 : 0;
    totalReadingTime.value += sessionTime;
    
    const progress: ReadingProgress = {
      comicId,
      chapterId,
      pageIndex: currentPageIndex.value,
      readingMode: readingMode.value,
      zoomLevel: zoomLevel.value,
      timestamp: currentTime,
      totalPages: chapter.value.pages.length,
      readingTime: totalReadingTime.value
    };
    
    readingProgressManager.saveProgress(progress);
    readingStartTime.value = currentTime;
  };
  
  // 加载阅读进度
  const loadReadingProgress = () => {
    const chapterId = parseInt(route.params.id as string);
    const comicId = chapter.value?.comicId;
    
    if (isNaN(chapterId) || !comicId) return;
    
    const progress = readingProgressManager.getProgress(comicId, chapterId);
    if (progress) {
      currentPageIndex.value = progress.pageIndex;
      readingMode.value = progress.readingMode;
      zoomLevel.value = progress.zoomLevel;
      totalReadingTime.value = progress.readingTime || 0;
    }
    
    // 加载阅读设置
    const settings = readingProgressManager.getSettings();
    isDarkMode.value = settings.isDarkMode;
    brightness.value = settings.brightness;
    autoHideToolbar.value = settings.autoHideToolbar;
    
    // 如果没有进度记录，使用默认设置
    if (!progress) {
      readingMode.value = settings.defaultReadingMode;
      zoomLevel.value = settings.defaultZoomLevel;
    }
  };
  
  // 保存阅读设置
  const saveReadingSettings = () => {
    readingProgressManager.saveSettings({
      isDarkMode: isDarkMode.value,
      brightness: brightness.value,
      autoHideToolbar: autoHideToolbar.value,
      defaultReadingMode: readingMode.value,
      defaultZoomLevel: zoomLevel.value
    });
  };
  
  // 开始自动保存定时器
  const startAutoSave = () => {
    if (autoSaveTimer.value) {
      clearInterval(autoSaveTimer.value);
    }
    
    // 每30秒自动保存一次进度
    autoSaveTimer.value = setInterval(() => {
      saveReadingProgress();
    }, 30000);
  };
  
  // 停止自动保存定时器
  const stopAutoSave = () => {
    if (autoSaveTimer.value) {
      clearInterval(autoSaveTimer.value);
      autoSaveTimer.value = undefined;
    }
  };

  // 键盘事件处理
  const handleKeydown = (event: KeyboardEvent) => {
    switch (event.key) {
      case 'ArrowLeft':
        if (readingMode.value === 'horizontal') {
          prevPage();
        }
        break;
      case 'ArrowRight':
        if (readingMode.value === 'horizontal') {
          nextPage();
        }
        break;
      case 'Escape':
        if (isFullscreen.value) {
          document.exitFullscreen();
        }
        break;
    }
  };

  // 获取章节详情
  const loadChapterDetail = async () => {
    try {
      loading.value = true;
      const chapterId = Number(route.params.id);
      if (!chapterId) {
        message.error('章节ID无效');
        router.back();
        return;
      }

      currentChapterId.value = chapterId;

      // 模拟章节数据
      chapter.value = {
        id: chapterId,
        title: '第1话 开始的故事',
        comicTitle: '示例漫画',
        comicId: 1,
        chapterNumber: 1,
        pages: [
          { imageUrl: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%201%20story%20beginning&image_size=portrait_4_3' },
          { imageUrl: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%202%20character%20introduction&image_size=portrait_4_3' },
          { imageUrl: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%203%20action%20scene&image_size=portrait_4_3' },
          { imageUrl: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%204%20dialogue%20scene&image_size=portrait_4_3' },
          { imageUrl: 'https://trae-api-sg.mchost.guru/api/ide/v1/text_to_image?prompt=manga%20page%205%20chapter%20ending&image_size=portrait_4_3' }
        ]
      };

      // 模拟章节列表
      chapterList.value = [
        { id: 1, title: '开始的故事', chapterNumber: 1 },
        { id: 2, title: '新的冒险', chapterNumber: 2 },
        { id: 3, title: '意外的相遇', chapterNumber: 3 }
      ];
    } catch (error) {
      console.error('获取章节详情失败:', error);
      message.error('获取章节详情失败');
      router.back();
    } finally {
      loading.value = false;
    }
  };

  // 组件挂载
  onMounted(() => {
    loadChapterDetail();
    document.addEventListener('keydown', handleKeydown);
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    
    // 加载阅读进度和设置
    loadReadingProgress();
    loadBookmarks();
    
    // 开始记录阅读时间
    readingStartTime.value = Date.now();
    
    // 开始自动保存
    startAutoSave();
    
    // 显示触摸提示（移动端）
    if ('ontouchstart' in window) {
      setTimeout(() => {
        showTouchHintDialog();
      }, 1000);
    }
  });

  // 组件卸载
  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown);
    document.removeEventListener('fullscreenchange', handleFullscreenChange);
    
    // 保存最终的阅读进度
    saveReadingProgress();
    saveReadingSettings();
    
    // 清理定时器
    if (hideToolbarTimer.value) {
      clearTimeout(hideToolbarTimer.value);
    }
    if (autoSaveTimer.value) {
      clearInterval(autoSaveTimer.value);
    }
    
    // 停止自动保存
    stopAutoSave();
  });
</script>

<style scoped lang="less">
  .chapter-reader-container {
    min-height: 100vh;
    background: #000;
    color: #fff;
  }

  .chapter-reader {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }

  .chapter-reader.dark-mode {
    background: #1a1a1a;
  }

  .reader-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    background: rgba(0, 0, 0, 0.8);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    position: sticky;
    top: 0;
    z-index: 100;
    transition: transform 0.3s ease;
  }

  .reader-header.hidden {
    transform: translateY(-100%);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .comic-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .comic-title {
    font-size: 14px;
    color: #ccc;
  }

  .chapter-title {
    font-size: 16px;
    font-weight: 600;
    color: #fff;
  }

  .chapter-navigation {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: rgba(0, 0, 0, 0.6);
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .reading-area {
    flex: 1;
    overflow: hidden;
    position: relative;

    &.fullscreen {
      position: fixed;
      top: 0;
      left: 0;
      width: 100vw;
      height: 100vh;
      z-index: 9999;
    }
  }

  .vertical-reader {
    height: 100%;
    overflow-y: auto;
    scroll-behavior: smooth;
  }

  .page-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 8px;
    position: relative;
  }

  .page-image {
    max-width: 100%;
    height: auto;
    display: block;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  }

  .page-number {
    position: absolute;
    bottom: 8px;
    right: 8px;
    background: rgba(0, 0, 0, 0.7);
    color: #fff;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
  }

  .horizontal-reader {
    height: 100%;
    display: flex;
    flex-direction: column;
    position: relative;
  }

  .page-slider {
    flex: 1;
    overflow: hidden;
    position: relative;
  }

  .page-track {
    display: flex;
    height: 100%;
    transition: transform 0.3s ease;
  }

  .page-slide {
    flex: 0 0 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 20px;

    .image-wrapper {
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.3s ease;
    }

    .page-image {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
      user-select: none;
      pointer-events: none;
    }
  }

  /* 双页模式样式 */
  .double-reader {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
  }

  .double-page-container {
    width: 100%;
    height: 100%;
    overflow: hidden;
    position: relative;
  }

  .double-pages {
    display: flex;
    height: 100%;
    transition: transform 0.3s ease;
  }

  .double-page-group {
    min-width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
    padding: 0 20px;
  }

  .double-page-item {
    flex: 1;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    position: relative;

    .image-wrapper {
      display: flex;
      align-items: center;
      justify-content: center;
      transition: transform 0.3s ease;
      max-height: 100%;
    }

    .page-image {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
      user-select: none;
    }
  }

  .double-controls {
    position: absolute;
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    align-items: center;
    gap: 20px;
    background: rgba(0, 0, 0, 0.7);
    padding: 10px 20px;
    border-radius: 25px;
    backdrop-filter: blur(10px);
    transition: opacity 0.3s ease, transform 0.3s ease;

    &.hidden {
      opacity: 0;
      transform: translateX(-50%) translateY(100%);
    }
  }

  .horizontal-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    background: rgba(0, 0, 0, 0.8);
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
  }

  .nav-button {
    background: rgba(255, 255, 255, 0.1);
    border: none;
    color: #fff;

    &:hover {
      background: rgba(255, 255, 255, 0.2);
    }

    &:disabled {
      background: rgba(255, 255, 255, 0.05);
      color: rgba(255, 255, 255, 0.3);
    }
  }

  .page-indicator {
    background: rgba(0, 0, 0, 0.7);
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 14px;
  }

  .reader-footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 80px;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  z-index: 1000;
  transition: transform 0.3s ease;
}

.reader-footer.hidden {
  transform: translateY(100%);
}

.footer-left,
.footer-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.footer-center {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
  justify-content: center;
}

.zoom-controls {
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  padding: 5px 10px;
  border-radius: 15px;
}

.zoom-text {
  color: white;
  font-size: 12px;
  min-width: 40px;
  text-align: center;
}

.progress-text {
  color: white;
  font-size: 12px;
  margin-left: 10px;
}

  .footer-center {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .progress-text {
    font-size: 12px;
    color: #ccc;
    min-width: 40px;
  }

  /* 设置面板样式 */
  .reader-settings {
    position: fixed;
    top: 60px;
    right: 0;
    width: 300px;
    height: calc(100vh - 140px);
    background: rgba(0, 0, 0, 0.9);
    backdrop-filter: blur(20px);
    border-left: 1px solid rgba(255, 255, 255, 0.1);
    padding: 20px;
    z-index: 999;
    transform: translateX(100%);
    transition: transform 0.3s ease;
    overflow-y: auto;
  }

  .reader-settings.show {
    transform: translateX(0);
  }

  .settings-section {
    margin-bottom: 25px;
  }

  .settings-section h4 {
    color: white;
    font-size: 14px;
    margin-bottom: 15px;
    font-weight: 500;
  }

  .settings-item {
    margin-bottom: 15px;
  }

  .settings-item label {
    color: rgba(255, 255, 255, 0.8);
    font-size: 12px;
    display: block;
    margin-bottom: 8px;
  }

  .mode-buttons {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }

  .mode-button {
    flex: 1;
    min-width: 70px;
    height: 32px;
    font-size: 12px;
  }

  .mode-button.active {
    background: #18a058;
    border-color: #18a058;
  }

  /* 触摸提示样式 */
  .touch-hint {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: rgba(0, 0, 0, 0.9);
    backdrop-filter: blur(20px);
    border-radius: 12px;
    padding: 20px;
    z-index: 1001;
    max-width: 280px;
    animation: fadeInOut 5s ease-in-out;
  }

  .hint-content {
    color: white;
    text-align: center;
  }

  .hint-content p {
    margin: 8px 0;
    font-size: 14px;
    line-height: 1.4;
  }

  .hint-content p:first-child {
    font-weight: 500;
    margin-bottom: 12px;
  }

  @keyframes fadeInOut {
    0%, 100% { opacity: 0; transform: translate(-50%, -50%) scale(0.9); }
    10%, 90% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
  }

  /* 设置面板样式 */
  .setting-group {
    margin-bottom: 20px;
  }
  
  .setting-group label {
    display: block;
    margin-bottom: 8px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.9);
  }
  
  .brightness-value {
    margin-left: 8px;
    font-size: 12px;
    color: #666;
  }
  
  .reading-stats {
    background: rgba(255, 255, 255, 0.1);
    padding: 12px;
    border-radius: 6px;
    font-size: 14px;
  }
  
  .reading-stats p {
    margin: 4px 0;
    color: rgba(255, 255, 255, 0.8);
  }
  
  /* 书签相关样式 */
  .bookmark-form {
    padding: 16px 0;
  }
  
  .bookmark-info {
    margin-bottom: 16px;
    padding: 12px;
    background: #f5f5f5;
    border-radius: 6px;
    font-size: 14px;
    line-height: 1.5;
  }
  
  .bookmark-list {
    height: 100%;
  }
  
  .empty-bookmarks {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: #999;
  }
  
  .empty-bookmarks p {
    margin-top: 16px;
    font-size: 14px;
  }
  
  .bookmark-item {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid #f0f0f0;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .bookmark-item:hover {
    background-color: #f5f5f5;
  }
  
  .bookmark-item:last-child {
    border-bottom: none;
  }
  
  .bookmark-item .bookmark-info {
    flex: 1;
    margin: 0;
    padding: 0;
    background: none;
  }
  
  .bookmark-item h4 {
    margin: 0 0 8px 0;
    font-size: 14px;
    font-weight: 500;
    color: #333;
  }
  
  .bookmark-item .bookmark-note {
    margin: 4px 0;
    font-size: 12px;
    color: #666;
    line-height: 1.4;
  }
  
  .bookmark-item .bookmark-time {
    margin: 4px 0 0 0;
    font-size: 11px;
    color: #999;
  }

  /* 响应式布局 */
  @media (max-width: 768px) {
    .reader-header {
      height: 50px;
      padding: 0 15px;
    }
    
    .reader-footer {
      height: 60px;
      padding: 0 15px;
    }
    
    .chapter-nav {
      gap: 10px;
    }
    
    .horizontal-controls {
      bottom: 80px;
    }
    
    .reader-settings {
      width: 100%;
      top: 50px;
      height: calc(100vh - 110px);
    }
    
    .zoom-controls {
      display: none;
    }
    
    .footer-center {
      gap: 10px;
    }
    
    .double-page-group {
      flex-direction: column;
      gap: 5px;
      padding: 10px;
    }
    
    .double-page-item {
      height: 50%;
    }

    .header-left {
      gap: 12px;
    }

    .comic-title {
      font-size: 12px;
    }

    .chapter-title {
      font-size: 14px;
    }

    .chapter-navigation {
      padding: 12px;
      gap: 12px;
    }

    .footer-center {
      .n-progress {
        width: 120px !important;
      }
    }
  }
</style>