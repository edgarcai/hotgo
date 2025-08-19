-- AniX 漫画阅读应用数据库表结构
-- 创建时间: 2024-01-20

-- 作者表
CREATE TABLE IF NOT EXISTS `anix_authors` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '作者ID',
    `name` VARCHAR(100) NOT NULL COMMENT '作者姓名',
    `bio` TEXT COMMENT '作者简介',
    `avatar_url` VARCHAR(255) DEFAULT NULL COMMENT '作者头像',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='作者表';

-- 漫画表
CREATE TABLE IF NOT EXISTS `anix_comics` (
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
    `rating` DECIMAL(3,2) DEFAULT 0.00 COMMENT '评分（0.00-10.00）',
    `rating_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '评分人数',
    `chapter_count` INT UNSIGNED DEFAULT 0 COMMENT '章节总数',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_title` (`title`),
    INDEX `idx_author_id` (`author_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_views` (`views` DESC),
    INDEX `idx_favorites` (`favorites` DESC),
    INDEX `idx_rating` (`rating` DESC),
    INDEX `idx_created_at` (`created_at` DESC),
    INDEX `idx_updated_at` (`updated_at` DESC),
    FOREIGN KEY (`author_id`) REFERENCES `anix_authors`(`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='漫画表';

-- 章节表
CREATE TABLE IF NOT EXISTS `anix_chapters` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '章节ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '所属漫画ID',
    `title` VARCHAR(200) NOT NULL COMMENT '章节标题',
    `chapter_number` INT UNSIGNED NOT NULL COMMENT '章节序号',
    `views` BIGINT UNSIGNED DEFAULT 0 COMMENT '章节浏览量',
    `page_count` INT UNSIGNED DEFAULT 0 COMMENT '页面总数',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uq_comic_chapter` (`comic_id`, `chapter_number`),
    INDEX `idx_comic_id` (`comic_id`),
    INDEX `idx_chapter_number` (`chapter_number`),
    INDEX `idx_views` (`views` DESC),
    INDEX `idx_created_at` (`created_at` DESC),
    FOREIGN KEY (`comic_id`) REFERENCES `anix_comics`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节表';

-- 章节页面表
CREATE TABLE IF NOT EXISTS `anix_chapter_pages` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '页面ID',
    `chapter_id` BIGINT UNSIGNED NOT NULL COMMENT '所属章节ID',
    `page_number` INT UNSIGNED NOT NULL COMMENT '页面序号',
    `image_url` VARCHAR(255) NOT NULL COMMENT '页面图片URL',
    `width` INT UNSIGNED DEFAULT 0 COMMENT '图片宽度',
    `height` INT UNSIGNED DEFAULT 0 COMMENT '图片高度',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_chapter_page` (`chapter_id`, `page_number`),
    INDEX `idx_chapter_id` (`chapter_id`),
    INDEX `idx_page_number` (`page_number`),
    FOREIGN KEY (`chapter_id`) REFERENCES `anix_chapters`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='章节页面表';

-- 用户行为表（收藏、点赞等）
CREATE TABLE IF NOT EXISTS `anix_user_actions` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '行为ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `action_type` ENUM('favorite','like','follow') NOT NULL COMMENT '行为类型：favorite-收藏，like-点赞，follow-关注',
    `target_type` ENUM('comic','chapter','author') NOT NULL COMMENT '目标类型：comic-漫画，chapter-章节，author-作者',
    `target_id` BIGINT UNSIGNED NOT NULL COMMENT '目标ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_user_action` (`user_id`, `action_type`, `target_type`, `target_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_target` (`target_type`, `target_id`),
    INDEX `idx_created_at` (`created_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户行为表';

-- 阅读历史表
CREATE TABLE IF NOT EXISTS `anix_reading_history` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '历史记录ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '漫画ID',
    `chapter_id` BIGINT UNSIGNED NOT NULL COMMENT '章节ID',
    `last_page` INT UNSIGNED DEFAULT 1 COMMENT '上次阅读到的页面',
    `progress` DECIMAL(5,2) DEFAULT 0.00 COMMENT '阅读进度百分比（0.00-100.00）',
    `reading_time` INT UNSIGNED DEFAULT 0 COMMENT '阅读时长（秒）',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次阅读时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后阅读时间',
    UNIQUE KEY `uq_user_comic` (`user_id`, `comic_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_comic_id` (`comic_id`),
    INDEX `idx_chapter_id` (`chapter_id`),
    INDEX `idx_updated_at` (`updated_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读历史表';

-- 漫画分类表
CREATE TABLE IF NOT EXISTS `anix_categories` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '分类ID',
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `description` VARCHAR(255) COMMENT '分类描述',
    `icon` VARCHAR(100) COMMENT '分类图标',
    `sort_order` INT UNSIGNED DEFAULT 0 COMMENT '排序权重',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用：1-启用，0-禁用',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uq_name` (`name`),
    INDEX `idx_sort_order` (`sort_order`),
    INDEX `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='漫画分类表';

-- 漫画分类关联表
CREATE TABLE IF NOT EXISTS `anix_comic_categories` (
    `id` BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '关联ID',
    `comic_id` BIGINT UNSIGNED NOT NULL COMMENT '漫画ID',
    `category_id` BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE KEY `uq_comic_category` (`comic_id`, `category_id`),
    INDEX `idx_comic_id` (`comic_id`),
    INDEX `idx_category_id` (`category_id`),
    FOREIGN KEY (`comic_id`) REFERENCES `anix_comics`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`category_id`) REFERENCES `anix_categories`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='漫画分类关联表';

-- 插入初始数据

-- 插入作者数据
INSERT INTO `anix_authors` (`name`, `bio`, `avatar_url`) VALUES 
('尾田荣一郎', '《海贼王》作者，日本著名漫画家，以其丰富的想象力和精彩的故事情节闻名于世。', '/images/authors/oda.jpg'),
('鸟山明', '《龙珠》作者，日本国民级漫画家，创造了众多经典角色，影响了一代又一代的读者。', '/images/authors/toriyama.jpg'),
('岸本齐史', '《火影忍者》作者，以忍者题材漫画著称，作品充满了友情、努力和胜利的主题。', '/images/authors/kishimoto.jpg'),
('青山刚昌', '《名侦探柯南》作者，推理漫画大师，创造了深受读者喜爱的侦探故事。', '/images/authors/aoyama.jpg'),
('富樫义博', '《幽游白书》《全职猎人》作者，以其独特的画风和深刻的故事内容著称。', '/images/authors/togashi.jpg');

-- 插入分类数据
INSERT INTO `anix_categories` (`name`, `description`, `icon`, `sort_order`) VALUES 
('热血', '充满激情和战斗的漫画作品', 'fire', 1),
('恋爱', '以爱情为主题的温馨漫画', 'heart', 2),
('冒险', '探索未知世界的冒险故事', 'compass', 3),
('搞笑', '轻松幽默的喜剧漫画', 'smile', 4),
('校园', '以学校生活为背景的青春漫画', 'school', 5),
('奇幻', '充满魔法和幻想元素的作品', 'magic-wand', 6),
('悬疑', '推理解谜类的烧脑漫画', 'search', 7),
('科幻', '未来科技题材的想象作品', 'rocket', 8),
('历史', '以历史背景为题材的漫画', 'book-open', 9),
('运动', '体育竞技类的励志漫画', 'trophy', 10);

-- 插入漫画数据
INSERT INTO `anix_comics` (`title`, `description`, `cover_url`, `author_id`, `status`, `tags`, `views`, `likes`, `favorites`, `rating`, `rating_count`, `chapter_count`) VALUES 
('海贼王', '路飞为了成为海贼王而展开的伟大航路冒险故事。这是一个关于友情、梦想和冒险的传奇故事，充满了激动人心的战斗和感人至深的情节。', '/images/covers/onepiece.jpg', 1, 'ongoing', '["热血", "冒险", "友情"]', 1500000, 89000, 156000, 9.8, 45000, 1000),
('龙珠', '孙悟空从小时候开始寻找龙珠的冒险之旅，后来发展成为保护地球的战士故事。经典的战斗漫画，影响了无数后来的作品。', '/images/covers/dragonball.jpg', 2, 'completed', '["热血", "战斗", "冒险"]', 1200000, 76000, 134000, 9.7, 38000, 519),
('火影忍者', '漩涡鸣人从吊车尾忍者成长为火影的励志故事。充满了忍术战斗、友情羁绊和成长历程的精彩漫画。', '/images/covers/naruto.jpg', 3, 'completed', '["热血", "忍者", "友情"]', 1100000, 72000, 128000, 9.6, 42000, 700),
('名侦探柯南', '高中生侦探工藤新一被神秘组织毒害后身体缩小，化名江户川柯南继续破案的推理故事。', '/images/covers/conan.jpg', 4, 'ongoing', '["悬疑", "推理", "校园"]', 980000, 65000, 112000, 9.4, 35000, 1100),
('全职猎人', '小杰为了寻找父亲而成为猎人，在这个过程中结识伙伴、面对挑战的冒险故事。', '/images/covers/hunterxhunter.jpg', 5, 'hiatus', '["冒险", "奇幻", "友情"]', 850000, 58000, 95000, 9.5, 28000, 390);

-- 插入章节数据（示例）
INSERT INTO `anix_chapters` (`comic_id`, `title`, `chapter_number`, `views`, `page_count`) VALUES 
(1, '我是路飞！要成为海贼王的男人！', 1, 45000, 19),
(1, '那个男人"小丑"巴基', 2, 42000, 20),
(1, '海贼猎人罗罗诺亚·索隆', 3, 40000, 18),
(1, '海军上尉"斧手"蒙卡', 4, 38000, 21),
(1, '海贼王和大剑客', 5, 36000, 19),
(2, '寻找龙珠', 1, 38000, 16),
(2, '龟仙人', 2, 35000, 18),
(2, '天下第一武道会', 3, 33000, 20),
(3, '漩涡鸣人登场！', 1, 42000, 22),
(3, '我的名字叫木叶丸！', 2, 39000, 20);

-- 插入漫画分类关联数据
INSERT INTO `anix_comic_categories` (`comic_id`, `category_id`) VALUES 
(1, 1), (1, 3), -- 海贼王：热血、冒险
(2, 1), (2, 3), -- 龙珠：热血、冒险
(3, 1), (3, 5), -- 火影忍者：热血、校园
(4, 7), (4, 5), -- 名侦探柯南：悬疑、校园
(5, 3), (5, 6); -- 全职猎人：冒险、奇幻