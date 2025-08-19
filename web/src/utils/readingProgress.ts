/**
 * 阅读进度管理工具类
 * 负责保存和获取用户的阅读进度、书签等信息
 */

export interface ReadingProgress {
  comicId: number;
  chapterId: number;
  pageIndex: number;
  readingMode: 'vertical' | 'horizontal' | 'double';
  zoomLevel: number;
  timestamp: number;
  totalPages: number;
  readingTime: number; // 阅读时长（秒）
}

export interface Bookmark {
  id: string;
  comicId: number;
  chapterId: number;
  pageIndex: number;
  title: string;
  note?: string;
  timestamp: number;
  thumbnail?: string;
}

export interface ReadingSettings {
  isDarkMode: boolean;
  brightness: number;
  autoHideToolbar: boolean;
  defaultReadingMode: 'vertical' | 'horizontal' | 'double';
  defaultZoomLevel: number;
}

class ReadingProgressManager {
  private readonly PROGRESS_KEY = 'anix_reading_progress';
  private readonly BOOKMARKS_KEY = 'anix_bookmarks';
  private readonly SETTINGS_KEY = 'anix_reading_settings';
  private readonly HISTORY_KEY = 'anix_reading_history';

  /**
   * 保存阅读进度
   */
  saveProgress(progress: ReadingProgress): void {
    try {
      const progressMap = this.getAllProgress();
      const key = `${progress.comicId}_${progress.chapterId}`;
      progressMap[key] = progress;
      
      localStorage.setItem(this.PROGRESS_KEY, JSON.stringify(progressMap));
      
      // 同时保存到阅读历史
      this.addToHistory(progress);
    } catch (error) {
      console.error('保存阅读进度失败:', error);
    }
  }

  /**
   * 获取指定章节的阅读进度
   */
  getProgress(comicId: number, chapterId: number): ReadingProgress | null {
    try {
      const progressMap = this.getAllProgress();
      const key = `${comicId}_${chapterId}`;
      return progressMap[key] || null;
    } catch (error) {
      console.error('获取阅读进度失败:', error);
      return null;
    }
  }

  /**
   * 获取漫画的最新阅读进度
   */
  getLatestProgress(comicId: number): ReadingProgress | null {
    try {
      const progressMap = this.getAllProgress();
      const comicProgresses = Object.values(progressMap)
        .filter(p => p.comicId === comicId)
        .sort((a, b) => b.timestamp - a.timestamp);
      
      return comicProgresses[0] || null;
    } catch (error) {
      console.error('获取最新阅读进度失败:', error);
      return null;
    }
  }

  /**
   * 获取所有阅读进度
   */
  private getAllProgress(): Record<string, ReadingProgress> {
    try {
      const data = localStorage.getItem(this.PROGRESS_KEY);
      return data ? JSON.parse(data) : {};
    } catch (error) {
      console.error('解析阅读进度数据失败:', error);
      return {};
    }
  }

  /**
   * 添加书签
   */
  addBookmark(bookmark: Omit<Bookmark, 'id' | 'timestamp'>): string {
    try {
      const bookmarks = this.getAllBookmarks();
      const id = `bookmark_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      const newBookmark: Bookmark = {
        ...bookmark,
        id,
        timestamp: Date.now()
      };
      
      bookmarks.push(newBookmark);
      localStorage.setItem(this.BOOKMARKS_KEY, JSON.stringify(bookmarks));
      
      return id;
    } catch (error) {
      console.error('添加书签失败:', error);
      throw error;
    }
  }

  /**
   * 删除书签
   */
  removeBookmark(bookmarkId: string): boolean {
    try {
      const bookmarks = this.getAllBookmarks();
      const filteredBookmarks = bookmarks.filter(b => b.id !== bookmarkId);
      
      if (filteredBookmarks.length === bookmarks.length) {
        return false; // 书签不存在
      }
      
      localStorage.setItem(this.BOOKMARKS_KEY, JSON.stringify(filteredBookmarks));
      return true;
    } catch (error) {
      console.error('删除书签失败:', error);
      return false;
    }
  }

  /**
   * 获取指定漫画的所有书签
   */
  getBookmarks(comicId: number): Bookmark[] {
    try {
      const bookmarks = this.getAllBookmarks();
      return bookmarks
        .filter(b => b.comicId === comicId)
        .sort((a, b) => b.timestamp - a.timestamp);
    } catch (error) {
      console.error('获取书签失败:', error);
      return [];
    }
  }

  /**
   * 获取所有书签
   */
  private getAllBookmarks(): Bookmark[] {
    try {
      const data = localStorage.getItem(this.BOOKMARKS_KEY);
      return data ? JSON.parse(data) : [];
    } catch (error) {
      console.error('解析书签数据失败:', error);
      return [];
    }
  }

  /**
   * 保存阅读设置
   */
  saveSettings(settings: ReadingSettings): void {
    try {
      localStorage.setItem(this.SETTINGS_KEY, JSON.stringify(settings));
    } catch (error) {
      console.error('保存阅读设置失败:', error);
    }
  }

  /**
   * 获取阅读设置
   */
  getSettings(): ReadingSettings {
    try {
      const data = localStorage.getItem(this.SETTINGS_KEY);
      const defaultSettings: ReadingSettings = {
        isDarkMode: false,
        brightness: 1,
        autoHideToolbar: false,
        defaultReadingMode: 'vertical',
        defaultZoomLevel: 1
      };
      
      return data ? { ...defaultSettings, ...JSON.parse(data) } : defaultSettings;
    } catch (error) {
      console.error('获取阅读设置失败:', error);
      return {
        isDarkMode: false,
        brightness: 1,
        autoHideToolbar: false,
        defaultReadingMode: 'vertical',
        defaultZoomLevel: 1
      };
    }
  }

  /**
   * 添加到阅读历史
   */
  private addToHistory(progress: ReadingProgress): void {
    try {
      const history = this.getReadingHistory();
      const existingIndex = history.findIndex(
        h => h.comicId === progress.comicId && h.chapterId === progress.chapterId
      );
      
      if (existingIndex >= 0) {
        // 更新现有记录
        history[existingIndex] = progress;
      } else {
        // 添加新记录
        history.unshift(progress);
      }
      
      // 只保留最近100条记录
      const limitedHistory = history.slice(0, 100);
      localStorage.setItem(this.HISTORY_KEY, JSON.stringify(limitedHistory));
    } catch (error) {
      console.error('添加阅读历史失败:', error);
    }
  }

  /**
   * 获取阅读历史
   */
  getReadingHistory(): ReadingProgress[] {
    try {
      const data = localStorage.getItem(this.HISTORY_KEY);
      return data ? JSON.parse(data) : [];
    } catch (error) {
      console.error('获取阅读历史失败:', error);
      return [];
    }
  }

  /**
   * 清除指定漫画的阅读数据
   */
  clearComicData(comicId: number): void {
    try {
      // 清除进度
      const progressMap = this.getAllProgress();
      Object.keys(progressMap).forEach(key => {
        if (progressMap[key].comicId === comicId) {
          delete progressMap[key];
        }
      });
      localStorage.setItem(this.PROGRESS_KEY, JSON.stringify(progressMap));
      
      // 清除书签
      const bookmarks = this.getAllBookmarks();
      const filteredBookmarks = bookmarks.filter(b => b.comicId !== comicId);
      localStorage.setItem(this.BOOKMARKS_KEY, JSON.stringify(filteredBookmarks));
      
      // 清除历史
      const history = this.getReadingHistory();
      const filteredHistory = history.filter(h => h.comicId !== comicId);
      localStorage.setItem(this.HISTORY_KEY, JSON.stringify(filteredHistory));
    } catch (error) {
      console.error('清除漫画数据失败:', error);
    }
  }

  /**
   * 清除所有阅读数据
   */
  clearAllData(): void {
    try {
      localStorage.removeItem(this.PROGRESS_KEY);
      localStorage.removeItem(this.BOOKMARKS_KEY);
      localStorage.removeItem(this.HISTORY_KEY);
    } catch (error) {
      console.error('清除所有数据失败:', error);
    }
  }
}

// 导出单例实例
export const readingProgressManager = new ReadingProgressManager();
export default readingProgressManager;