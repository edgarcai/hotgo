-- HotGo自动生成2FA管理菜单权限SQL
-- 2FA管理功能菜单配置
-- Version: 2.13.1
-- Date: 2024-01-17 12:00:00
-- Link https://github.com/bufanyun/hotgo

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;

--
-- 数据库： `hotgo`
--

-- --------------------------------------------------------

--
-- 插入表中的数据 `hg_admin_menu`
--

SET @now := now();

-- 查找会员管理菜单的ID（假设已存在）
SET @memberMenuId = (SELECT id FROM `hg_admin_menu` WHERE `name` = 'member' AND `type` = 1 LIMIT 1);

-- 如果会员管理菜单不存在，创建一个
INSERT IGNORE INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, '0', '会员管理', 'member', '/member', 'UserOutlined', '1', '/member/list', '', '', 'ParentLayout', '1', '', '0', '0', '', '0', '0', '0', '1', 'tr_0 ', '300', '会员管理功能', '1', @now, @now);

-- 获取会员管理菜单ID
SET @memberMenuId = (SELECT id FROM `hg_admin_menu` WHERE `name` = 'member' AND `type` = 1 LIMIT 1);

-- 2FA管理主目录
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, @memberMenuId, '2FA管理', 'twoFactorManagement', '/member/2fa', 'SafetyOutlined', '1', '/member/2fa/user', '', '', 'ParentLayout', '1', '', '0', '0', '', '0', '0', '0', '2', CONCAT('tr_', @memberMenuId, ' '), '400', '双因素认证管理', '1', @now, @now);

SET @twoFactorMenuId = LAST_INSERT_ID();

-- 2FA用户管理页面
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, @twoFactorMenuId, '用户管理', 'twoFactorUser', '/member/2fa/user', '', '2', '', '/member/twofa/list,/member/twofa/edit,/member/twofa/delete,/member/twofa/reset', '2FA用户管理', '/admin/member/twofa-user/index', '0', '', '0', '0', '', '1', '0', '0', '3', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' '), '100', '2FA用户管理页面', '1', @now, @now);

SET @twoFactorUserMenuId = LAST_INSERT_ID();

-- 2FA用户管理按钮权限
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES 
(NULL, @twoFactorUserMenuId, '查看', 'twoFactorUserView', '', '', '3', '', '/member/twofa/list', '查看2FA用户', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorUserMenuId, ' '), '100', '', '1', @now, @now),
(NULL, @twoFactorUserMenuId, '启用/禁用2FA', 'twoFactorUserToggle', '', '', '3', '', '/member/twofa/edit', '启用/禁用用户2FA', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorUserMenuId, ' '), '200', '', '1', @now, @now),
(NULL, @twoFactorUserMenuId, '重置备用码', 'twoFactorUserReset', '', '', '3', '', '/member/twofa/reset', '重置用户2FA备用码', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorUserMenuId, ' '), '300', '', '1', @now, @now);

-- 2FA操作日志页面
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, @twoFactorMenuId, '操作日志', 'twoFactorAudit', '/member/2fa/audit', '', '2', '', '/log/twofa/list,/log/twofa/export', '2FA操作日志', '/admin/member/twofa-audit/index', '0', '', '0', '0', '', '1', '0', '0', '3', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' '), '200', '2FA操作审计日志', '1', @now, @now);

SET @twoFactorAuditMenuId = LAST_INSERT_ID();

-- 2FA操作日志按钮权限
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES 
(NULL, @twoFactorAuditMenuId, '查看', 'twoFactorAuditView', '', '', '3', '', '/log/twofa/list', '查看2FA操作日志', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorAuditMenuId, ' '), '100', '', '1', @now, @now),
(NULL, @twoFactorAuditMenuId, '导出', 'twoFactorAuditExport', '', '', '3', '', '/log/twofa/export', '导出2FA操作日志', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorAuditMenuId, ' '), '200', '', '1', @now, @now);

-- 2FA统计页面
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, @twoFactorMenuId, '使用统计', 'twoFactorStats', '/member/2fa/stats', '', '2', '', '/statistics/twofa/usage,/statistics/twofa/trends', '2FA使用统计', '/admin/member/twofa-stats/index', '0', '', '0', '0', '', '1', '0', '0', '3', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' '), '300', '2FA使用情况统计', '1', @now, @now);

SET @twoFactorStatsMenuId = LAST_INSERT_ID();

-- 2FA统计按钮权限
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES 
(NULL, @twoFactorStatsMenuId, '查看统计', 'twoFactorStatsView', '', '', '3', '', '/statistics/twofa/usage', '查看2FA使用统计', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorStatsMenuId, ' '), '100', '', '1', @now, @now),
(NULL, @twoFactorStatsMenuId, '查看趋势', 'twoFactorStatsTrends', '', '', '3', '', '/statistics/twofa/trends', '查看2FA使用趋势', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorStatsMenuId, ' '), '200', '', '1', @now, @now);

-- 2FA策略管理页面
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) 
VALUES (NULL, @twoFactorMenuId, '策略管理', 'twoFactorPolicy', '/member/2fa/policy', '', '2', '', '/config/twofa/policy,/config/twofa/save', '2FA策略管理', '/admin/member/twofa-policy/index', '0', '', '0', '0', '', '1', '0', '0', '3', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' '), '400', '2FA安全策略配置', '1', @now, @now);

SET @twoFactorPolicyMenuId = LAST_INSERT_ID();

-- 2FA策略管理按钮权限
INSERT INTO `hg_admin_menu` (`id`, `pid`, `title`, `name`, `path`, `icon`, `type`, `redirect`, `permissions`, `permission_name`, `component`, `always_show`, `active_menu`, `is_root`, `is_frame`, `frame_src`, `keep_alive`, `hidden`, `affix`, `level`, `tree`, `sort`, `remark`, `status`, `created_at`, `updated_at`) VALUES 
(NULL, @twoFactorPolicyMenuId, '查看策略', 'twoFactorPolicyView', '', '', '3', '', '/config/twofa/policy', '查看2FA策略配置', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorPolicyMenuId, ' '), '100', '', '1', @now, @now),
(NULL, @twoFactorPolicyMenuId, '保存策略', 'twoFactorPolicySave', '', '', '3', '', '/config/twofa/save', '保存2FA策略配置', '', '0', '', '0', '0', '', '0', '0', '0', '4', CONCAT('tr_', @memberMenuId, ' tr_', @twoFactorMenuId, ' tr_', @twoFactorPolicyMenuId, ' '), '200', '', '1', @now, @now);

COMMIT;

-- 显示创建的菜单信息
SELECT 
    id,
    pid,
    title,
    name,
    path,
    type,
    permissions,
    component,
    sort,
    remark
FROM `hg_admin_menu` 
WHERE 
    (name LIKE '%twoFactor%' OR name = 'member')
    AND created_at >= DATE_SUB(@now, INTERVAL 1 MINUTE)
ORDER BY id ASC;

-- 使用说明：
-- 1. 执行此SQL脚本将在数据库中创建2FA管理相关的菜单项
-- 2. 菜单结构：会员管理 -> 2FA管理 -> [用户管理|操作日志|使用统计|策略管理]
-- 3. 每个页面都包含相应的按钮权限配置
-- 4. 执行后需要在角色权限中分配相应的菜单权限给管理员角色
-- 5. 前端路由将自动根据菜单配置生成动态路由