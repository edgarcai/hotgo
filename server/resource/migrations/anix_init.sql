-- AniX 漫画阅读应用数据库初始化脚本
-- 创建时间: 2024-01-20
-- 描述: 创建AniX应用所需的所有数据表

-- 创建作者表
CREATE TABLE `anix_authors` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '作者ID',
    `name` VARCHAR(100) NOT NULL COMMENT '作者姓名',
    `bio` TEXT COMMENT '作者简介',
    `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '作者头像',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_authors_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作者表';

-- 创建漫画表
CREATE TABLE `anix_comics` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '漫画ID',
    `title` VARCHAR(200) NOT NULL COMMENT '漫画标题',
    `description` TEXT COMMENT '漫画简介',
    `cover_url` VARCHAR(255) COMMENT '封面图片URL',
    `author_id` BIGINT UNSIGNED COMMENT '作者ID',
    `status` ENUM('ongoing','completed','hiatus') DEFAULT 'ongoing' COMMENT '连载状态：ongoing-连载中，completed-已完结，hiatus-暂停',
    `tags` JSON COMMENT '标签数组，如["热血", "校园", "魔幻"]',
    `views` BIGINT UNSIGNED DEFAULT 0 COMMENT '总浏览量',
    `likes` BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `favorites` BIGINT UNSIGNED DEFAULT 0 COMMENT '收藏数',
    `rating` DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分(0-10)',
    `rating_count` INT UNSIGNED DEFAULT 0 COMMENT '评分人数',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_comics_title` (`title`),
    INDEX `idx_comics_author_id` (`author_id`),
    INDEX `idx_comics_status` (`status`),
    INDEX `idx_comics_views` (`views` DESC),
    INDEX `idx_comics_favorites` (`favorites` DESC),
    INDEX `idx_comics_rating` (`rating` DESC),
    INDEX `idx_comics_created_at` (`created_at` DESC),
    FOREIGN KEY (`author_id`) REFERENCES `anix_authors`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='漫画表';

-- 创建章节表
CREATE TABLE `anix_chapters` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '章节ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '所属漫画ID',
    `title` VARCHAR(200) NOT NULL COMMENT '章节标题',
    `chapter_number` INT UNSIGNED NOT NULL COMMENT '章节序号',
    `views` BIGINT UNSIGNED DEFAULT 0 COMMENT '章节浏览量',
    `page_count` INT UNSIGNED DEFAULT 0 COMMENT '页数',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uq_comic_chapter` (`comic_id`, `chapter_number`),
    INDEX `idx_chapters_comic_id` (`comic_id`),
    INDEX `idx_chapters_number` (`chapter_number`),
    INDEX `idx_chapters_views` (`views` DESC),
    FOREIGN KEY (`comic_id`) REFERENCES `anix_comics`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节表';

-- 创建章节页面表
CREATE TABLE `anix_chapter_pages` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '页面ID',
    `chapter_id` BIGINT UNSIGNED NOT NULL COMMENT '所属章节ID',
    `page_number` INT UNSIGNED NOT NULL COMMENT '页码',
    `image_url` VARCHAR(255) NOT NULL COMMENT '图片URL',
    `image_width` INT UNSIGNED DEFAULT 0 COMMENT '图片宽度',
    `image_height` INT UNSIGNED DEFAULT 0 COMMENT '图片高度',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_chapter_page` (`chapter_id`, `page_number`),
    INDEX `idx_pages_chapter_id` (`chapter_id`),
    INDEX `idx_pages_number` (`page_number`),
    FOREIGN KEY (`chapter_id`) REFERENCES `anix_chapters`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节页面表';

-- 创建用户表
CREATE TABLE `anix_users` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    `email` VARCHAR(100) UNIQUE COMMENT '邮箱地址',
    `phone` VARCHAR(20) UNIQUE COMMENT '手机号码',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
    `salt` VARCHAR(32) NOT NULL COMMENT '密码盐',
    `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '头像URL',
    `nickname` VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    `gender` ENUM('male','female','secret') DEFAULT 'secret' COMMENT '性别',
    `birthday` DATE DEFAULT NULL COMMENT '生日',
    `bio` VARCHAR(255) DEFAULT NULL COMMENT '个人简介',
    `level` INT UNSIGNED DEFAULT 1 COMMENT '用户等级',
    `exp` INT UNSIGNED DEFAULT 0 COMMENT '经验值',
    `coins` INT UNSIGNED DEFAULT 0 COMMENT '金币数量',
    `vip_level` TINYINT UNSIGNED DEFAULT 0 COMMENT 'VIP等级',
    `vip_expire_at` TIMESTAMP NULL COMMENT 'VIP过期时间',
    `last_login_at` TIMESTAMP NULL COMMENT '最后登录时间',
    `last_login_ip` VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_users_username` (`username`),
    INDEX `idx_users_email` (`email`),
    INDEX `idx_users_phone` (`phone`),
    INDEX `idx_users_level` (`level` DESC),
    INDEX `idx_users_created_at` (`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 创建用户行为表（收藏、点赞、关注）
CREATE TABLE `anix_user_actions` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '行为ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `action_type` ENUM('favorite','like','follow') NOT NULL COMMENT '行为类型：favorite-收藏，like-点赞，follow-关注',
    `target_type` ENUM('comic','chapter','author','user') NOT NULL COMMENT '目标类型：comic-漫画，chapter-章节，author-作者，user-用户',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标对象ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_user_action` (`user_id`, `action_type`, `target_type`, `target_id`),
    INDEX `idx_user_actions_user_id` (`user_id`),
    INDEX `idx_user_actions_target` (`target_type`, `target_id`),
    INDEX `idx_user_actions_type` (`action_type`),
    FOREIGN KEY (`user_id`) REFERENCES `anix_users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户行为表';

-- 创建阅读历史表
CREATE TABLE `anix_reading_history` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '历史记录ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '漫画ID',
    `chapter_id` BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    `last_page` INT UNSIGNED DEFAULT 1 COMMENT '上次阅读到的页码',
    `reading_time` INT UNSIGNED DEFAULT 0 COMMENT '阅读时长（秒）',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次阅读时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后阅读时间',
    UNIQUE KEY `uq_user_comic` (`user_id`, `comic_id`),
    INDEX `idx_reading_history_user_id` (`user_id`),
    INDEX `idx_reading_history_comic_id` (`comic_id`),
    INDEX `idx_reading_history_updated_at` (`updated_at` DESC),
    FOREIGN KEY (`user_id`) REFERENCES `anix_users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`comic_id`) REFERENCES `anix_comics`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`chapter_id`) REFERENCES `anix_chapters`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读历史表';

-- 创建评论表
CREATE TABLE `anix_comments` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '评论ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `target_type` ENUM('comic','chapter') NOT NULL COMMENT '评论目标类型：comic-漫画，chapter-章节',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标对象ID',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父评论ID，0表示顶级评论',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `likes` INT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `dislikes` INT UNSIGNED DEFAULT 0 COMMENT '踩数',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-隐藏，-1-删除',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_comments_user_id` (`user_id`),
    INDEX `idx_comments_target` (`target_type`, `target_id`),
    INDEX `idx_comments_parent_id` (`parent_id`),
    INDEX `idx_comments_created_at` (`created_at` DESC),
    INDEX `idx_comments_likes` (`likes` DESC),
    FOREIGN KEY (`user_id`) REFERENCES `anix_users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- 创建弹幕表
CREATE TABLE `anix_danmu` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '弹幕ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `chapter_id` BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    `page_number` INT UNSIGNED NOT NULL COMMENT '页码',
    `content` VARCHAR(200) NOT NULL COMMENT '弹幕内容',
    `color` VARCHAR(7) DEFAULT '#FFFFFF' COMMENT '弹幕颜色',
    `position` ENUM('top','middle','bottom') DEFAULT 'middle' COMMENT '弹幕位置',
    `time_offset` INT UNSIGNED DEFAULT 0 COMMENT '时间偏移（毫秒）',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-正常，0-隐藏，-1-删除',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_danmu_user_id` (`user_id`),
    INDEX `idx_danmu_chapter_page` (`chapter_id`, `page_number`),
    INDEX `idx_danmu_created_at` (`created_at` DESC),
    FOREIGN KEY (`user_id`) REFERENCES `anix_users`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`chapter_id`) REFERENCES `anix_chapters`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='弹幕表';

-- 创建分类表
CREATE TABLE `anix_categories` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `slug` VARCHAR(50) NOT NULL UNIQUE COMMENT '分类标识',
    `description` VARCHAR(255) DEFAULT NULL COMMENT '分类描述',
    `icon` VARCHAR(100) DEFAULT NULL COMMENT '分类图标',
    `color` VARCHAR(7) DEFAULT '#1890ff' COMMENT '分类颜色',
    `sort_order` INT DEFAULT 0 COMMENT '排序权重',
    `status` TINYINT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_categories_slug` (`slug`),
    INDEX `idx_categories_sort` (`sort_order` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分类表';

-- 创建漫画分类关联表
CREATE TABLE `anix_comic_categories` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '漫画ID',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_comic_category` (`comic_id`, `category_id`),
    INDEX `idx_comic_categories_comic_id` (`comic_id`),
    INDEX `idx_comic_categories_category_id` (`category_id`),
    FOREIGN KEY (`comic_id`) REFERENCES `anix_comics`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`category_id`) REFERENCES `anix_categories`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='漫画分类关联表';

-- 创建系统配置表
CREATE TABLE `anix_configs` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '配置ID',
    `key` VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键',
    `value` TEXT COMMENT '配置值',
    `description` VARCHAR(255) DEFAULT NULL COMMENT '配置描述',
    `type` ENUM('string','int','bool','json') DEFAULT 'string' COMMENT '配置类型',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_configs_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 插入初始数据

-- 插入作者数据
INSERT INTO `anix_authors` (`name`, `bio`) VALUES 
('尾田荣一郎', '《海贼王》作者，日本著名漫画家，以其丰富的想象力和精彩的故事情节著称'),
('鸟山明', '《龙珠》作者，日本国民级漫画家，对全球动漫文化产生了深远影响'),
('岸本齐史', '《火影忍者》作者，以忍者题材漫画闻名世界'),
('久保带人', '《死神》作者，以独特的画风和战斗场面著称'),
('富樫义博', '《幽游白书》《全职猎人》作者，剧情构思精妙');

-- 插入分类数据
INSERT INTO `anix_categories` (`name`, `slug`, `description`, `icon`, `color`, `sort_order`) VALUES 
('热血', 'action', '充满激情与战斗的漫画作品', 'fire', '#ff4d4f', 100),
('恋爱', 'romance', '以爱情为主题的温馨漫画', 'heart', '#ff85c0', 90),
('冒险', 'adventure', '探索未知世界的冒险故事', 'compass', '#52c41a', 80),
('搞笑', 'comedy', '轻松幽默的搞笑漫画', 'smile', '#faad14', 70),
('校园', 'school', '以校园生活为背景的青春漫画', 'book', '#1890ff', 60),
('奇幻', 'fantasy', '充满魔法与奇幻元素的作品', 'magic-wand', '#722ed1', 50),
('科幻', 'sci-fi', '未来科技与科学幻想题材', 'rocket', '#13c2c2', 40),
('悬疑', 'mystery', '推理解谜类悬疑漫画', 'search', '#595959', 30);

-- 插入漫画数据
INSERT INTO `anix_comics` (`title`, `description`, `author_id`, `status`, `tags`, `views`, `likes`, `favorites`, `rating`, `rating_count`) VALUES 
('海贼王', '路飞为了成为海贼王而展开的伟大冒险故事，充满友情、梦想与热血的经典之作', 1, 'ongoing', '["热血", "冒险", "友情", "梦想"]', 1500000, 89000, 156000, 9.8, 25000),
('龙珠', '孙悟空从小到大的成长历程，寻找龙珠的冒险之旅，经典的战斗漫画', 2, 'completed', '["热血", "战斗", "冒险", "成长"]', 1200000, 76000, 134000, 9.6, 20000),
('火影忍者', '鸣人从吊车尾到火影的成长之路，忍者世界的精彩故事', 3, 'completed', '["热血", "忍者", "友情", "成长"]', 1100000, 72000, 128000, 9.5, 18000),
('死神', '黑崎一护获得死神力量后的冒险故事，灵魂世界的精彩战斗', 4, 'completed', '["热血", "超自然", "战斗"]', 800000, 58000, 95000, 9.2, 15000),
('全职猎人', '小杰寻找父亲的冒险旅程，猎人世界的奇妙故事', 5, 'hiatus', '["冒险", "奇幻", "成长"]', 750000, 62000, 98000, 9.7, 12000);

-- 插入章节数据（为每部漫画创建几个示例章节）
INSERT INTO `anix_chapters` (`comic_id`, `title`, `chapter_number`, `views`, `page_count`) VALUES 
-- 海贼王章节
(1, '我是路飞！要成为海贼王的男人！', 1, 45000, 20),
(1, '戴草帽的路飞', 2, 42000, 18),
(1, '索隆登场', 3, 40000, 22),
(1, '海军上校摩根', 4, 38000, 19),
-- 龙珠章节
(2, '寻找龙珠', 1, 35000, 16),
(2, '龟仙人', 2, 33000, 18),
(3, '天下第一武道会', 3, 32000, 20),
-- 火影忍者章节
(3, '漩涡鸣人登场', 1, 38000, 21),
(3, '我的名字叫佐助', 2, 36000, 19),
(3, '你已经是一个优秀的忍者了', 3, 35000, 23);

-- 插入漫画分类关联数据
INSERT INTO `anix_comic_categories` (`comic_id`, `category_id`) VALUES 
-- 海贼王：热血、冒险
(1, 1), (1, 3),
-- 龙珠：热血、冒险
(2, 1), (2, 3),
-- 火影忍者：热血
(3, 1),
-- 死神：热血、奇幻
(4, 1), (4, 6),
-- 全职猎人：冒险、奇幻
(5, 3), (5, 6);

-- 插入系统配置数据
INSERT INTO `anix_configs` (`key`, `value`, `description`, `type`) VALUES 
('site_name', 'AniX漫画', '网站名称', 'string'),
('site_description', '专业的二次元漫画阅读平台', '网站描述', 'string'),
('max_upload_size', '10485760', '最大上传文件大小（字节）', 'int'),
('allow_register', 'true', '是否允许用户注册', 'bool'),
('default_avatar', '/static/images/default_avatar.png', '默认头像路径', 'string'),
('reading_settings', '{"default_mode":"vertical","auto_next":false,"night_mode":false}', '默认阅读设置', 'json');