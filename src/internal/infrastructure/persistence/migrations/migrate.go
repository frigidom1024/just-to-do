package mysql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Migrator 数据库迁移器
type Migrator struct {
	db *sqlx.DB
}

// NewMigrator 创建迁移器
func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{db: db}
}

// Migration 迁移记录
type Migration struct {
	Version   int64  `db:"version"`
	Name      string `db:"name"`
	AppliedAt string `db:"applied_at"`
}

// migrations 所有迁移脚本
var migrations = []struct {
	version int64
	name    string
	up      func(db *sqlx.DB) error
	down    func(db *sqlx.DB) error
}{
	{
		version: 20240117000001,
		name:    "create_users_table",
		up:      createUsersTable,
		down:    dropUsersTable,
	},
	// 添加新的迁移脚本
}

// Up 执行所有未执行的迁移
func (m *Migrator) Up(ctx context.Context) error {
	// 创建迁移记录表
	if err := m.createMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// 获取已执行的迁移
	appliedVersions, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// 执行未执行的迁移
	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.version]; applied {
			fmt.Printf("Migration %d (%s) already applied, skipping\n", migration.version, migration.name)
			continue
		}

		fmt.Printf("Applying migration %d (%s)...\n", migration.version, migration.name)
		if err := migration.up(m.db); err != nil {
			return fmt.Errorf("failed to apply migration %d (%s): %w", migration.version, migration.name, err)
		}

		// 记录迁移
		if err := m.recordMigration(ctx, migration.version, migration.name); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.version, err)
		}

		fmt.Printf("Migration %d (%s) applied successfully\n", migration.version, migration.name)
	}

	fmt.Println("All migrations applied successfully")
	return nil
}

// Down 回滚最后一个迁移
func (m *Migrator) Down(ctx context.Context) error {
	appliedVersions, err := m.getAppliedVersions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// 找到最后一个迁移
	var lastMigration *struct {
		version int64
		name    string
		up      func(db *sqlx.DB) error
		down    func(db *sqlx.DB) error
	}

	for i := len(migrations) - 1; i >= 0; i-- {
		if _, applied := appliedVersions[migrations[i].version]; applied {
			lastMigration = &migrations[i]
			break
		}
	}

	if lastMigration == nil {
		fmt.Println("No migration to rollback")
		return nil
	}

	fmt.Printf("Rolling back migration %d (%s)...\n", lastMigration.version, lastMigration.name)
	if err := lastMigration.down(m.db); err != nil {
		return fmt.Errorf("failed to rollback migration %d (%s): %w", lastMigration.version, lastMigration.name, err)
	}

	// 删除迁移记录
	if err := m.deleteMigration(ctx, lastMigration.version); err != nil {
		return fmt.Errorf("failed to delete migration record %d: %w", lastMigration.version, err)
	}

	fmt.Printf("Migration %d (%s) rolled back successfully\n", lastMigration.version, lastMigration.name)
	return nil
}

// createMigrationTable 创建迁移记录表
func (m *Migrator) createMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`
	_, err := m.db.ExecContext(ctx, query)
	return err
}

// getAppliedVersions 获取已执行的迁移版本
func (m *Migrator) getAppliedVersions(ctx context.Context) (map[int64]bool, error) {
	query := `SELECT version, name FROM schema_migrations ORDER BY version`
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[int64]bool)
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name); err != nil {
			return nil, err
		}
		versions[m.Version] = true
	}

	return versions, rows.Err()
}

// recordMigration 记录迁移
func (m *Migrator) recordMigration(ctx context.Context, version int64, name string) error {
	query := `INSERT INTO schema_migrations (version, name) VALUES (?, ?)`
	_, err := m.db.ExecContext(ctx, query, version, name)
	return err
}

// deleteMigration 删除迁移记录
func (m *Migrator) deleteMigration(ctx context.Context, version int64) error {
	query := `DELETE FROM schema_migrations WHERE version = ?`
	_, err := m.db.ExecContext(ctx, query, version)
	return err
}

// ==================== 迁移脚本 ====================

// createUsersTable 创建用户表
func createUsersTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
			username VARCHAR(50) NOT NULL COMMENT '用户名',
			email VARCHAR(100) NOT NULL COMMENT '邮箱',
			password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
			avatar_url VARCHAR(500) DEFAULT '' COMMENT '头像URL',
			status VARCHAR(20) NOT NULL DEFAULT 'active' COMMENT '用户状态: active/inactive/banned',
			created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
			updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
			deleted_at DATETIME(3) DEFAULT NULL COMMENT '删除时间（软删除）',
			PRIMARY KEY (id),
			UNIQUE KEY uk_username (username),
			UNIQUE KEY uk_email (email),
			KEY idx_status (status),
			KEY idx_created_at (created_at),
			KEY idx_deleted_at (deleted_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表'
	`
	_, err := db.Exec(query)
	return err
}

// dropUsersTable 删除用户表
func dropUsersTable(db *sqlx.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}
