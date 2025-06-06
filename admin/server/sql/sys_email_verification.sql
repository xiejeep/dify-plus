-- 创建邮箱验证码表
CREATE TABLE IF NOT EXISTS `sys_email_verifications` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `email` varchar(255) NOT NULL COMMENT '邮箱地址',
  `code` varchar(6) NOT NULL COMMENT '验证码',
  `type` int(11) NOT NULL DEFAULT '1' COMMENT '验证码类型 1:注册 2:找回密码',
  `expired_at` datetime NOT NULL COMMENT '过期时间',
  `used` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已使用',
  `ip` varchar(45) DEFAULT NULL COMMENT '请求IP',
  PRIMARY KEY (`id`),
  KEY `idx_email` (`email`),
  KEY `idx_sys_email_verifications_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='邮箱验证码表'; 