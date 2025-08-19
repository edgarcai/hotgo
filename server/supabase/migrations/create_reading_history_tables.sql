-- 创建阅读历史记录表
CREATE TABLE IF NOT EXISTS anix_reading_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    comic_id BIGINT NOT NULL COMMENT '漫画ID',
    chapter_id BIGINT NOT NULL COMMENT '章节ID',
    page_index INTEGER DEFAULT 0 COMMENT '当前阅读页面索引',
    progress DECIMAL(5,2) DEFAULT 0.00 COMMENT '阅读进度百分比',
    reading_time BIGINT DEFAULT 0 COMMENT '累计阅读时长(秒)',
    last_read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后阅读时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE(user_id, comic_id, chapter_id)
);

-- 创建书签表
CREATE TABLE IF NOT EXISTS anix_bookmarks (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    comic_id BIGINT NOT NULL COMMENT '漫画ID',
    chapter_id BIGINT NOT NULL COMMENT '章节ID',
    page_index INTEGER DEFAULT 0 COMMENT '书签页面索引',
    note TEXT COMMENT '书签备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE(user_id, comic_id, chapter_id, page_index)
);

-- 创建阅读统计表
CREATE TABLE IF NOT EXISTS anix_reading_stats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    comic_id BIGINT NOT NULL COMMENT '漫画ID',
    total_chapters INTEGER DEFAULT 0 COMMENT '总章节数',
    read_chapters INTEGER DEFAULT 0 COMMENT '已读章节数',
    total_reading_time BIGINT DEFAULT 0 COMMENT '总阅读时长(秒)',
    first_read_at TIMESTAMP COMMENT '首次阅读时间',
    last_read_at TIMESTAMP COMMENT '最后阅读时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE(user_id, comic_id)
);

-- 创建索引
CREATE INDEX idx_reading_history_user_id ON anix_reading_history(user_id);
CREATE INDEX idx_reading_history_comic_id ON anix_reading_history(comic_id);
CREATE INDEX idx_reading_history_last_read ON anix_reading_history(last_read_at DESC);

CREATE INDEX idx_bookmarks_user_id ON anix_bookmarks(user_id);
CREATE INDEX idx_bookmarks_comic_id ON anix_bookmarks(comic_id);

CREATE INDEX idx_reading_stats_user_id ON anix_reading_stats(user_id);
CREATE INDEX idx_reading_stats_last_read ON anix_reading_stats(last_read_at DESC);

-- MySQL不支持RLS策略，权限控制在应用层实现