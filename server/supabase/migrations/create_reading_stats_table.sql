-- 创建阅读统计表
CREATE TABLE IF NOT EXISTS `anix_reading_stats` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `comic_id` bigint(20) unsigned NOT NULL COMMENT '漫画ID',
  `total_reading_time` int(11) NOT NULL DEFAULT '0' COMMENT '总阅读时长(秒)',
  `chapters_read` int(11) NOT NULL DEFAULT '0' COMMENT '已读章节数',
  `pages_read` int(11) NOT NULL DEFAULT '0' COMMENT '已读页数',
  `last_read_at` datetime DEFAULT NULL COMMENT '最后阅读时间',
  `completion_rate` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '完成度百分比',
  `reading_streak` int(11) NOT NULL DEFAULT '0' COMMENT '连续阅读天数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_comic` (`user_id`, `comic_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_comic_id` (`comic_id`),
  KEY `idx_last_read_at` (`last_read_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='阅读统计表';