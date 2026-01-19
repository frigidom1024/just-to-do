package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/request"
	"todolist/internal/interfaces/http/response"
)

// TestRegisterUserHandler 测试用户注册接口
func TestRegisterUserHandler(t *testing.T) {
	// 测试用例1：无效注册请求 - 密码太短
	t.Run("invalid password - too short", func(t *testing.T) {
		// 创建注册请求
		req := request.RegisterUserRequest{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "short1",
		}

		resp, err := handler.RegisterUserHandler(context.Background(), req)
		// 由于密码太短，应该返回错误
		assert.Error(t, err)
		assert.Equal(t, response.UserResponse{}, resp)
	})
}

// TestChangePasswordHandler 测试修改密码接口
func TestChangePasswordHandler(t *testing.T) {
	// 测试用例：无效的上下文（没有用户信息）
	t.Run("invalid context - no user", func(t *testing.T) {
		req := request.ChangePasswordRequest{
			OldPassword: "oldpass",
			NewPassword: "NewPass123!",
		}

		resp, err := handler.ChangePasswordHandler(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
		assert.Equal(t, response.MessageResponse{}, resp)
	})
}

// TestUpdateEmailHandler 测试更新邮箱接口
func TestUpdateEmailHandler(t *testing.T) {
	// 测试用例：无效的上下文（没有用户信息）
	t.Run("invalid context - no user", func(t *testing.T) {
		req := request.UpdateEmailRequest{
			NewEmail: "newemail@example.com",
		}

		resp, err := handler.UpdateEmailHandler(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
		assert.Equal(t, response.MessageResponse{}, resp)
	})
}

// TestUpdateAvatarHandler 测试更新头像接口
func TestUpdateAvatarHandler(t *testing.T) {
	// 测试用例：无效的上下文（没有用户信息）
	t.Run("invalid context - no user", func(t *testing.T) {
		req := request.UpdateAvatarRequest{
			AvatarURL: "https://example.com/avatar.jpg",
		}

		resp, err := handler.UpdateAvatarHandler(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "unauthorized: invalid user context", err.Error())
		assert.Equal(t, response.MessageResponse{}, resp)
	})
}
