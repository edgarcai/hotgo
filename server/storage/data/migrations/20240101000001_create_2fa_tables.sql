-- 2FA功能数据库迁移文件
-- 创建时间: 2024-01-01
-- 描述: 创建二次因子验证相关的数据表

-- 创建2FA配置表
CREATE TABLE `hg_admin_2fa_config` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
    `admin_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID',
    `secret_key` varchar(32) NOT NULL COMMENT 'TOTP密钥',
    `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用',
    `enabled_at` datetime DEFAULT NULL COMMENT '启用时间',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_admin_id` (`admin_id`),
    KEY `idx_enabled` (`enabled`),
    CONSTRAINT `fk_admin_2fa_config_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员2FA配置表';

-- 创建备用恢复码表
CREATE TABLE `hg_admin_2fa_backup_code` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '备用码ID',
    `admin_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID',
    `code` varchar(16) NOT NULL COMMENT '备用恢复码',
    `used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已使用',
    `used_at` datetime DEFAULT NULL COMMENT '使用时间',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_admin_id` (`admin_id`),
    KEY `idx_used` (`used`),
    CONSTRAINT `fk_admin_2fa_backup_code_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='2FA备用恢复码表';

-- 创建2FA操作日志表
CREATE TABLE `hg_admin_2fa_log` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `admin_id` bigint(20) unsigned NOT NULL COMMENT '管理员ID',
    `action` varchar(50) NOT NULL COMMENT '操作类型',
    `ip` varchar(45) NOT NULL COMMENT 'IP地址',
    `user_agent` text COMMENT '用户代理',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '操作状态',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注信息',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_admin_id` (`admin_id`),
    KEY `idx_action` (`action`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_admin_2fa_log_admin_id` FOREIGN KEY (`admin_id`) REFERENCES `hg_admin_member` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='2FA操作日志表';

-- 在sys_config表中插入2FA相关配置
INSERT INTO `hg_sys_config` (`group`, `key`, `value`, `type`, `title`, `tip`, `rule`, `extend`, `remark`, `sort`, `status`, `created_at`, `updated_at`) VALUES
('security', '2faEnabled', '0', 'switch', '启用二次验证', '是否启用系统二次验证功能', '', '', '控制整个系统的2FA功能开关', 100, 1, NOW(), NOW()),
('security', '2faEnforced', '0', 'switch', '强制二次验证', '是否强制所有用户启用2FA', '', '', '强制策略，启用后所有用户必须配置2FA', 101, 1, NOW(), NOW()),
('security', '2faGracePeriod', '7', 'input', '2FA宽限期', '强制启用2FA后的宽限期（天）', 'integer|min:1|max:30', '', '用户必须在宽限期内完成2FA配置', 102, 1, NOW(), NOW()),
('security', '2faMaxAttempts', '5', 'input', '最大尝试次数', '2FA验证最大失败次数', 'integer|min:3|max:10', '', '超过次数将锁定账户', 103, 1, NOW(), NOW()),
('security', '2faLockoutDuration', '30', 'input', '锁定时长', '验证失败锁定时长（分钟）', 'integer|min:5|max:1440', '', '账户锁定的持续时间', 104, 1, NOW(), NOW());

-- 创建索引以优化查询性能
CREATE INDEX idx_admin_2fa_config_admin_id ON hg_admin_2fa_config(admin_id);
CREATE INDEX idx_admin_2fa_config_enabled ON hg_admin_2fa_config(enabled);
CREATE INDEX idx_admin_2fa_backup_code_admin_id ON hg_admin_2fa_backup_code(admin_id);
CREATE INDEX idx_admin_2fa_backup_code_used ON hg_admin_2fa_backup_code(used);
CREATE INDEX idx_admin_2fa_log_admin_id ON hg_admin_2fa_log(admin_id);
CREATE INDEX idx_admin_2fa_log_created_at ON hg_admin_2fa_log(created_at DESC);
CREATE INDEX idx_admin_2fa_log_action ON hg_admin_2fa_log(action);