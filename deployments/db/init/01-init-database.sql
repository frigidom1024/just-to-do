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
  `email` VARCHAR(255) NOT NULL COMMENT '邮箱',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `avatar_url` VARCHAR(500) DEFAULT '' COMMENT '头像URL',
  `status` VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '用户状态: active/inactive/suspended',
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
-- 创建 daily_notes 表（每日笔记）
-- ====================================================================
DROP TABLE IF EXISTS `daily_notes`;
CREATE TABLE `daily_notes` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '笔记ID',
  `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
  `note_date` DATE NOT NULL COMMENT '笔记日期',
  `content` TEXT NOT NULL COMMENT '笔记内容',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_date` (`user_id`, `note_date`) COMMENT '用户+日期唯一索引',
  KEY `idx_user_id` (`user_id`),
  KEY `idx_note_date` (`note_date`),
  CONSTRAINT `fk_daily_notes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='每日笔记表';

-- ====================================================================
-- 创建 todos 表（待办事项）
-- ====================================================================
DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '待办ID',
  `title` VARCHAR(255) NOT NULL COMMENT '标题',
  `description` TEXT COMMENT '描述内容',
  `status` VARCHAR(20) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/in_progress/completed',
  `priority` VARCHAR(20) NOT NULL DEFAULT 'medium' COMMENT '优先级: low/medium/high',
  `estimated_start_time` DATETIME(3) DEFAULT NULL COMMENT '预计开始时间',
  `estimated_end_time` DATETIME(3) DEFAULT NULL COMMENT '预计结束时间',
  `actual_start_time` DATETIME(3) DEFAULT NULL COMMENT '实际开始时间',
  `actual_end_time` DATETIME(3) DEFAULT NULL COMMENT '实际结束时间',
  `daily_note_id` BIGINT(20) UNSIGNED DEFAULT NULL COMMENT '关联的每日笔记ID',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_priority` (`priority`),
  KEY `idx_daily_note_id` (`daily_note_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_todos_daily_note` FOREIGN KEY (`daily_note_id`) REFERENCES `daily_notes` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='待办事项表';

-- ====================================================================
-- 创建 notes 表（备注）
-- ====================================================================
DROP TABLE IF EXISTS `notes`;
CREATE TABLE `notes` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '备注ID',
  `todo_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '关联待办事项ID',
  `content` TEXT NOT NULL COMMENT '备注内容',
  `created_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_todo_id` (`todo_id`),
  CONSTRAINT `fk_notes_todo` FOREIGN KEY (`todo_id`) REFERENCES `todos` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='备注表';

-- ====================================================================
-- 插入测试数据
-- ====================================================================

-- 插入测试用户（密码: 123456，使用 bcrypt 哈希）
-- 注意：实际使用中应使用强密码
INSERT INTO `users` (`username`, `email`, `password_hash`, `avatar_url`, `status`, `created_at`, `updated_at`) VALUES
('admin', 'admin@todo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '', 'active', NOW(3), NOW(3)),
('test_user', 'test@todo.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '', 'active', NOW(3), NOW(3));

-- 插入测试每日笔记
INSERT INTO `daily_notes` (`user_id`, `note_date`, `content`, `created_at`, `updated_at`) VALUES
(1, CURDATE(), '今天是项目的第一天，开始搭建基础架构。', NOW(3), NOW(3)),
(1, DATE_SUB(CURDATE(), INTERVAL 1 DAY), '昨天完成了需求分析。', NOW(3), NOW(3)),
(2, CURDATE(), '测试用户日记：准备开始测试功能。', NOW(3), NOW(3));

-- 插入测试待办事项
INSERT INTO `todos` (`title`, `description`, `status`, `priority`, `daily_note_id`, `created_at`, `updated_at`) VALUES
('完成数据库设计', '设计并创建所有必要的数据库表', 'completed', 'high', 1, NOW(3), NOW(3)),
('实现用户认证', '实现用户注册、登录和JWT认证功能', 'in_progress', 'high', 1, NOW(3), NOW(3)),
('开发待办事项功能', '实现待办事项的CRUD操作', 'pending', 'medium', 1, NOW(3), NOW(3)),
('编写单元测试', '为核心功能编写单元测试', 'pending', 'low', NULL, NOW(3), NOW(3));

-- 插入测试备注
INSERT INTO `notes` (`todo_id`, `content`, `created_at`, `updated_at`) VALUES
(1, '已按照ER图完成表结构设计，包含users、daily_notes、todos和notes表。', NOW(3), NOW(3)),
(2, '需要注意JWT密钥的安全性，使用环境变量配置。', NOW(3), NOW(3)),
(4, '使用testify框架编写测试用例。', NOW(3), NOW(3));

SET FOREIGN_KEY_CHECKS = 1;
