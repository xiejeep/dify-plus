-- 积分管理菜单数据
-- 注意：运行此脚本前请确保没有重复的菜单ID

-- 添加积分管理主菜单 (ID: 41)
INSERT INTO `sys_base_menus` 
(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) 
VALUES 
(41, NOW(), NOW(), NULL, 0, 0, 'points', 'pointsManagement', 0, 'view/gaia/points/index.vue', 3, '', 0, 0, '积分管理', 'coin', 0);

-- 添加用户积分管理子菜单 (ID: 42)
INSERT INTO `sys_base_menus` 
(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) 
VALUES 
(42, NOW(), NOW(), NULL, 0, 41, 'users', 'pointsUsers', 0, 'view/gaia/points/users.vue', 1, '', 1, 0, '用户积分管理', 'user', 0);

-- 添加签到记录子菜单 (ID: 43)  
INSERT INTO `sys_base_menus` 
(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) 
VALUES 
(43, NOW(), NOW(), NULL, 0, 41, 'records', 'pointsRecords', 0, 'view/routerHolder.vue', 2, '', 0, 0, '签到记录管理', 'document', 0);

-- 添加积分流水子菜单 (ID: 44)
INSERT INTO `sys_base_menus` 
(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) 
VALUES 
(44, NOW(), NOW(), NULL, 0, 41, 'transactions', 'pointsTransactions', 0, 'view/routerHolder.vue', 3, '', 0, 0, '积分流水管理', 'money', 0);

-- 添加积分配置子菜单 (ID: 45)
INSERT INTO `sys_base_menus` 
(`id`, `created_at`, `updated_at`, `deleted_at`, `menu_level`, `parent_id`, `path`, `name`, `hidden`, `component`, `sort`, `active_name`, `keep_alive`, `default_menu`, `title`, `icon`, `close_tab`) 
VALUES 
(45, NOW(), NOW(), NULL, 0, 41, 'config', 'pointsConfig', 0, 'view/gaia/points/config.vue', 4, '', 1, 0, '积分配置管理', 'setting', 0);

-- 为超级管理员角色(authority_id=888)添加积分管理菜单权限
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('888', '41');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('888', '42');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('888', '43');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('888', '44');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('888', '45');

-- 为管理员角色(authority_id=9528)添加积分管理菜单权限
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('9528', '41');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('9528', '42');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('9528', '43');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('9528', '44');
INSERT INTO `sys_authority_menus` (`sys_authority_authority_id`, `sys_base_menu_id`) VALUES ('9528', '45');

-- 查询验证菜单是否插入成功
SELECT 
  m.id, 
  m.title, 
  m.path, 
  m.name, 
  m.parent_id, 
  m.sort,
  m.icon,
  CASE WHEN m.parent_id = 0 THEN '主菜单' ELSE '子菜单' END as menu_type
FROM sys_base_menus m 
WHERE m.id BETWEEN 41 AND 45 
OR m.name LIKE 'points%'
ORDER BY m.id; 