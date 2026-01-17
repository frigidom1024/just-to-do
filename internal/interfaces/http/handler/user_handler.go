package handler

import (
	"encoding/json"
	"net/http"

	"todolist/internal/domain/user"
	"todolist/internal/interfaces/http/request"
	"todolist/internal/interfaces/http/response"
	"todolist/internal/pkg/auth"
	"todolist/internal/pkg/logger"
)

// UserHandler 用户 HTTP 处理器。
//
// 负责处理用户相关的 HTTP 请求，包括注册、登录、
// 密码修改、邮箱更新等操作。
type UserHandler struct {
	userService user.UserService
	tokenTool   auth.TokenTool
}

// NewUserHandler 创建用户处理器。
//
// 参数：
//   userService - 用户领域服务
//   tokenTool - Token 工具
//
// 返回：
//   *UserHandler - 用户处理器实例
func NewUserHandler(userService user.UserService, tokenTool auth.TokenTool) *UserHandler {
	return &UserHandler{
		userService: userService,
		tokenTool:   tokenTool,
	}
}

// RegisterUser 处理用户注册请求。
//
// 方法：POST /api/users/register
//
// 请求体：request.RegisterUserRequest
// 成功响应：201 Created, response.UserResponse
// 错误响应：
//   - 400 Bad Request 参数验证失败
//   - 409 Conflict 用户名或邮箱已存在
//   - 500 Internal Server Error 服务器错误
func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 解析请求
	var req request.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorContext(ctx, "解析注册请求失败",
			logger.Err(err))
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 创建值对象（包含验证）
	username, err := user.NewUsername(req.Username)
	if err != nil {
		logger.WarnContext(ctx, "用户名格式无效",
			logger.String("username", req.Username),
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	email, err := user.NewEmail(req.Email)
	if err != nil {
		logger.WarnContext(ctx, "邮箱格式无效",
			logger.String("email", req.Email),
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		logger.WarnContext(ctx, "密码格式无效",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 调用领域服务
	userEntity, err := h.userService.RegisterUser(ctx, username, email, password)
	if err != nil {
		logger.ErrorContext(ctx, "注册用户失败",
			logger.String("username", req.Username),
			logger.Err(err))

		// 根据错误类型返回不同的状态码
		switch {
		case err == user.ErrUsernameTaken:
			respondWithError(w, http.StatusConflict, "Username already taken")
		case err == user.ErrEmailAlreadyExists:
			respondWithError(w, http.StatusConflict, "Email already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to register user")
		}
		return
	}

	logger.InfoContext(ctx, "用户注册成功",
		logger.Int64("user_id", userEntity.GetID()),
		logger.String("username", userEntity.GetUsername()))

	// 返回响应
	respondWithJSON(w, http.StatusCreated, response.ToUserResponse(userEntity))
}

// LoginUser 处理用户登录请求。
//
// 方法：POST /api/users/login
//
// 请求体：request.LoginUserRequest
// 成功响应：200 OK, response.LoginResponse
// 错误响应：
//   - 400 Bad Request 参数验证失败
//   - 401 Unauthorized 认证失败
//   - 500 Internal Server Error 服务器错误
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// 解析请求
	var req request.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorContext(ctx, "解析登录请求失败",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// 注意：不记录密码到日志中
	logger.InfoContext(ctx, "用户登录尝试",
		logger.String("email", req.Email))

	// 创建值对象
	email, err := user.NewEmail(req.Email)
	if err != nil {
		logger.WarnContext(ctx, "邮箱格式无效",
			logger.String("email", req.Email),
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	password, err := user.NewPassword(req.Password)
	if err != nil {
		logger.WarnContext(ctx, "密码格式无效",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 认证用户
	userEntity, err := h.userService.AuthenticateUser(ctx, email, password)
	if err != nil {
		logger.WarnContext(ctx, "用户认证失败",
			logger.String("email", req.Email),
			// 不记录密码
			logger.Err(err))
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// 生成 Token
	token, err := h.tokenTool.GenerateToken(
		userEntity.GetID(),
		userEntity.GetUsername(),
		"user", // TODO: 从用户实体获取角色
	)
	if err != nil {
		logger.ErrorContext(ctx, "生成 Token 失败",
			logger.Int64("user_id", userEntity.GetID()),
			logger.Err(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	logger.InfoContext(ctx, "用户登录成功",
		logger.Int64("user_id", userEntity.GetID()),
		logger.String("username", userEntity.GetUsername()))

	// 返回响应
	respondWithJSON(w, http.StatusOK, response.LoginResponse{
		Token: token,
		User:  response.ToUserResponse(userEntity),
	})
}

// ChangePassword 处理修改密码请求。
//
// 方法：POST /api/users/change-password
//
// 请求头：Authorization: Bearer <token>
// 请求体：request.ChangePasswordRequest
// 成功响应：200 OK
// 错误响应：
//   - 400 Bad Request 参数验证失败
//   - 401 Unauthorized 未认证
//   - 403 Forbidden 旧密码错误
//   - 500 Internal Server Error 服务器错误
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: 从 Token 中获取 userID
	userID := int64(1) // 临时硬编码

	// 解析请求
	var req request.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.ErrorContext(ctx, "解析修改密码请求失败",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// 创建值对象
	oldPassword, err := user.NewPassword(req.OldPassword)
	if err != nil {
		logger.WarnContext(ctx, "旧密码格式无效",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newPassword, err := user.NewPassword(req.NewPassword)
	if err != nil {
		logger.WarnContext(ctx, "新密码格式无效",
			logger.Err(err))
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 调用领域服务
	err = h.userService.ChangePassword(ctx, userID, oldPassword, newPassword)
	if err != nil {
		logger.ErrorContext(ctx, "修改密码失败",
			logger.Int64("user_id", userID),
			logger.Err(err))

		switch {
		case err == user.ErrOldPasswordIncorrect:
			respondWithError(w, http.StatusForbidden, "Old password is incorrect")
		default:
			respondWithError(w, http.StatusInternalServerError, "Failed to change password")
		}
		return
	}

	logger.InfoContext(ctx, "密码修改成功",
		logger.Int64("user_id", userID))

	// 返回响应
	respondWithJSON(w, http.StatusOK, response.MessageResponse{
		Message: "Password changed successfully",
	})
}

// respondWithJSON 返回 JSON 响应。
//
// 参数：
//   w - HTTP 响应写入器
//   status - HTTP 状态码
//   payload - 响应数据
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError 返回错误响应。
//
// 参数：
//   w - HTTP 响应写入器
//   status - HTTP 状态码
//   message - 错误消息
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, response.ErrorResponse{
		Message: message,
	})
}
