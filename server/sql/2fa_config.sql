-- 添加2FA相关的系统配置项

-- 2FA全局开关
INSERT INTO `hg_sys_config` (`group`, `name`, `type`, `key`, `value`, `default_value`, `sort`, `tip`, `is_default`, `status`) VALUES
('basic', '2FA全局开关', 'int', 'basicTwoFASwitch', '0', '0', 100, '是否启用双因子认证功能，0:关闭 1:开启', 1, 1);

-- 2FA强制启用
INSERT INTO `hg_sys_config` (`group`, `name`, `type`, `key`, `value`, `default_value`, `sort`, `tip`, `is_default`, `status`) VALUES
('basic', '强制启用2FA', 'int', 'basicTwoFAForce', '0', '0', 101, '是否强制用户启用2FA，0:可选 1:强制', 1, 1);

-- 2FA发行者名称
INSERT INTO `hg_sys_config` (`group`, `name`, `type`, `key`, `value`, `default_value`, `sort`, `tip`, `is_default`, `status`) VALUES
('basic', '2FA发行者名称', 'string', 'basicTwoFAIssuer', 'HotGo', 'HotGo', 102, '2FA应用中显示的发行者名称', 1, 1);