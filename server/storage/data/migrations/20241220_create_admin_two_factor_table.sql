-- 创建管理员双因子认证表
CREATE TABLE IF NOT EXISTS `hg_admin_two_factor` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `member_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID',
  `secret_key` text NOT NULL COMMENT '加密后的TOTP密钥',
  `backup_codes` json DEFAULT NULL COMMENT '备用恢复码（哈希后）',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用（0:未启用 1:已启用）',
  `enabled_at` datetime DEFAULT NULL COMMENT '启用时间',
  `last_used_at` datetime DEFAULT NULL COMMENT '最后使用时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_member_id` (`member_id`),
  KEY `idx_member_id` (`member_id`),
  KEY `idx_is_enabled` (`is_enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员双因子认证表';