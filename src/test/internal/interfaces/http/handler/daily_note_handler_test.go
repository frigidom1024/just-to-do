package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/request"
)

// TestCreateDailyNoteHandler 测试创建每日笔记接口
func TestCreateDailyNoteHandler(t *testing.T) {
	// 测试用例1：无效上下文 - 没有用户信息
	t.Run("invalid context - no user", func(t *testing.T) {
		// 创建请求
		req := request.DailyNoteRequest{
			Content: "测试内容",
		}

		resp, err := handler.CreateDailyNoteHandler(context.Background(), req)
		// 由于没有用户信息，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
	})
}

// TestGetTodayDailyNoteHandler 测试获取今日每日笔记接口
func TestGetTodayDailyNoteHandler(t *testing.T) {
	// 测试用例1：无效上下文 - 没有用户信息
	t.Run("invalid context - no user", func(t *testing.T) {
		resp, err := handler.GetTodayDailyNoteHandler(context.Background(), request.EmptyRequest{})
		// 由于没有用户信息，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
	})
}

// TestGetDailyNoteListHandler 测试分页获取每日笔记列表接口
func TestGetDailyNoteListHandler(t *testing.T) {
	// 测试用例1：无效上下文 - 没有用户信息
	t.Run("invalid context - no user", func(t *testing.T) {
		resp, err := handler.GetDailyNoteListHandler(context.Background(), request.EmptyRequest{})
		// 由于没有用户信息，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
	})
}

// TestUpdateDailyNoteHandler 测试更新今日每日笔记接口
func TestUpdateDailyNoteHandler(t *testing.T) {
	// 测试用例1：无效上下文 - 没有用户信息
	t.Run("invalid context - no user", func(t *testing.T) {
		// 创建请求
		req := request.DailyNoteRequest{
			Content: "更新后的内容",
		}

		resp, err := handler.UpdateDailyNoteHandler(context.Background(), req)
		// 由于没有用户信息，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
	})
}

// TestDeleteDailyNoteHandler 测试删除今日每日笔记接口
func TestDeleteDailyNoteHandler(t *testing.T) {
	// 测试用例1：无效上下文 - 没有用户信息
	t.Run("invalid context - no user", func(t *testing.T) {
		resp, err := handler.DeleteDailyNoteHandler(context.Background(), request.EmptyRequest{})
		// 由于没有用户信息，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
	})
}
