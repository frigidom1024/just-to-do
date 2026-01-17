package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"todolist/internal/infrastructure/config"

	"github.com/jmoiron/sqlx"
)

// Client 数据库客户端
// 封装数据库操作，提供简洁的 API
type Client struct {
	db *sqlx.DB
}

// NewClient 创建数据库客户端
func NewClient() (*Client, error) {
	cfg, err := config.GetMySQLConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get mysql config: %w", err)
	}

	db, err := sqlx.Connect("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}

	// 设置连接池
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// 验证连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return &Client{db: db}, nil
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// GetDB 获取底层 *sqlx.DB（用于复杂操作）
func (c *Client) GetDB() *sqlx.DB {
	return c.db
}

// ==================== 查询操作 ====================

// Query 查询多行数据并映射到切片
// dest: 目标切片的指针，如 &[]UserDO{}
// query: SQL 查询语句
// args: 查询参数
func (c *Client) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.SelectContext(ctx, dest, query, args...)
}

// QueryOne 查询单行数据并映射到结构体
// dest: 目标结构体的指针，如 &UserDO{}
// query: SQL 查询语句
// args: 查询参数
func (c *Client) QueryOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.db.GetContext(ctx, dest, query, args...)
}

// QueryRow 查询单行数据，返回 *sqlx.Row（用于自定义扫描）
func (c *Client) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return c.db.QueryRowxContext(ctx, query, args...)
}

// ==================== 执行操作 ====================

// Exec 执行 SQL 语句（INSERT, UPDATE, DELETE）
// 返回 sql.Result 包含 LastInsertId 和 RowsAffected
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (sqlResult, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// ExecWithID 执行 INSERT 并返回插入的 ID
func (c *Client) ExecWithID(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

// ExecWithAffected 执行 SQL 并返回影响的行数
func (c *Client) ExecWithAffected(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return affected, nil
}

// ==================== 事务操作 ====================

// BeginTxs 开启事务
func (c *Client) BeginTxs(ctx context.Context) (*Tx, error) {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &Tx{tx: tx}, nil
}

// BeginTxWithOpts 开启事务（自定义选项）
func (c *Client) BeginTxWithOpts(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := c.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &Tx{tx: tx}, nil
}

// Transaction 执行事务函数（自动提交/回滚）
func (c *Client) Transaction(ctx context.Context, fn func(*Tx) error) error {
	_, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	tx, err := c.BeginTxs(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // 重新抛出 panic
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// ==================== 批量操作 ====================

// BatchExec 批量执行 SQL 语句
// 返回每个语句的执行结果
func (c *Client) BatchExec(ctx context.Context, queries []string, argsList [][]interface{}) ([]sqlResult, error) {
	results := make([]sqlResult, 0, len(queries))

	tx, err := c.BeginTxs(ctx)
	if err != nil {
		return nil, err
	}

	for i, query := range queries {
		args := argsList[i]
		result, err := tx.Exec(ctx, query, args...)
		if err != nil {
			_ = tx.Rollback()
			return nil, fmt.Errorf("batch exec failed at index %d: %w", i, err)
		}
		results = append(results, result)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("batch commit failed: %w", err)
	}

	return results, nil
}

// ==================== 便捷方法 ====================

// Exists 检查数据是否存在
func (c *Client) Exists(ctx context.Context, query string, args ...interface{}) (bool, error) {
	var count int
	if err := c.db.GetContext(ctx, &count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count 统计行数
func (c *Client) Count(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	if err := c.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}

// ==================== 事务封装 ====================

// Tx 事务封装
type Tx struct {
	tx *sqlx.Tx
}

// Commit 提交事务
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback 回滚事务
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// ==================== 事务查询操作 ====================

// Query 事务中查询多行数据
func (t *Tx) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.SelectContext(ctx, dest, query, args...)
}

// QueryOne 事务中查询单行数据
func (t *Tx) QueryOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return t.tx.GetContext(ctx, dest, query, args...)
}

// QueryRow 事务中查询单行数据（自定义扫描）
func (t *Tx) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return t.tx.QueryRowxContext(ctx, query, args...)
}

// ==================== 事务执行操作 ====================

// Exec 事务中执行 SQL
func (t *Tx) Exec(ctx context.Context, query string, args ...interface{}) (sqlResult, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// ExecWithID 事务中执行 INSERT 并返回 ID
func (t *Tx) ExecWithID(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}
	return id, nil
}

// ExecWithAffected 事务中执行 SQL 并返回影响行数
func (t *Tx) ExecWithAffected(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}
	return affected, nil
}

// Exists 事务中检查数据是否存在
func (t *Tx) Exists(ctx context.Context, query string, args ...interface{}) (bool, error) {
	var count int
	if err := t.tx.GetContext(ctx, &count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count 事务中统计行数
func (t *Tx) Count(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	if err := t.tx.GetContext(ctx, &count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}

// sqlResult 接口，用于类型断言
type sqlResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
