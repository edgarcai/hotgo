-- 创建书签表
CREATE TABLE IF NOT EXISTS `anix_bookmark` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '书签ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `comic_id` bigint(20) NOT NULL COMMENT '漫画ID',
  `chapter_id` bigint(20) NOT NULL COMMENT '章节ID',
  `page_index` int(11) NOT NULL DEFAULT '0' COMMENT '页码索引',
  `note` varchar(500) DEFAULT '' COMMENT '书签备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_comic_id` (`comic_id`),
  KEY `idx_chapter_id` (`chapter_id`),
  KEY `idx_user_comic` (`user_id`, `comic_id`),
  UNIQUE KEY `uk_user_comic_chapter_page` (`user_id`, `comic_id`, `chapter_id`, `page_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='书签表';