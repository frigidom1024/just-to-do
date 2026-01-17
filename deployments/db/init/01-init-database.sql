-- ====================================================================
-- 数据库初始化脚本
-- 项目: Todo
-- 数据库: MySQL 8.0
-- 字符集: utf8mb4
-- ====================================================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ====================================================================
-- 创建 users 表
-- ====================================================================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` VARCHAR(50) NOT NULL COMMENT '用户名',
  `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `avatar_url` VARCHAR(500) DEFAULT '' COMMENT '头像URL',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '用户状态: active/inactive/banned',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ====================================================================
-- 插入测试数据
-- ====================================================================

-- 插入管理员用户（密码: 123456，使用 bcrypt 哈希）
-- 注意：实际使用中应使用强密码
INSERT INTO `users` (`username`, `email`, `password_hash`, `avatar_url`, `status`, `created_at`, `updated_at`) VALUES
('admin', 'admin@todo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '', 'active', NOW(3), NOW(3)),
('test_user', 'test@todo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '', 'active', NOW(3), NOW(3));

SET FOREIGN_KEY_CHECKS = 1;
