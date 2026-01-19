package handler

import (
	"context"
	"errors"

	"todolist/internal/interfaces/http/middleware"
	request "todolist/internal/interfaces/http/request"
	response "todolist/internal/interfaces/http/response"

	dailynoteapp "todolist/internal/application/daily_note"
	dailynote "todolist/internal/domain/daily_note"
	"todolist/internal/infrastructure/persistence/mysql"
)

// CreateDailyNoteHandler 创建每日笔记处理器
func CreateDailyNoteHandler(ctx context.Context, req request.DailyNoteRequest) (response.DailyNoteResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewDailyNoteRepository()
	dailyNoteService := dailynote.NewService(repo)
	dailyNoteAppService := dailynoteapp.NewDailyNoteApplicationService(dailyNoteService)

	// 2. 从上下文中获取用户信息（由认证中间件设置）
	user, ok := middleware.GetDataFromContext(ctx)
	if !ok {
		return response.DailyNoteResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务创建每日笔记
	dailyNoteDTO, err := dailyNoteAppService.CreateDailyNote(ctx, user.UserID, req.Content)
	if err != nil {
		return response.DailyNoteResponse{}, err
	}

	// 4. 转换为HTTP响应
	return response.ToDailyNoteResponse(*dailyNoteDTO), nil
}

// GetTodayDailyNoteHandler 获取今日的每日笔记处理器
func GetTodayDailyNoteHandler(ctx context.Context, req request.EmptyRequest) (response.DailyNoteResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewDailyNoteRepository()
	dailyNoteService := dailynote.NewService(repo)
	dailyNoteAppService := dailynoteapp.NewDailyNoteApplicationService(dailyNoteService)

	// 2. 从上下文中获取用户信息（由认证中间件设置）
	user, ok := middleware.GetDataFromContext(ctx)
	if !ok {
		return response.DailyNoteResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务获取今日笔记
	dailyNoteDTO, err := dailyNoteAppService.GetTodayDailyNote(ctx, user.UserID)
	if err != nil {
		return response.DailyNoteResponse{}, err
	}

	// 4. 转换为HTTP响应
	return response.ToDailyNoteResponse(*dailyNoteDTO), nil
}

// GetDailyNoteListHandler 分页获取每日笔记列表处理器
func GetDailyNoteListHandler(ctx context.Context, req request.EmptyRequest) (response.DailyNoteListResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewDailyNoteRepository()
	dailyNoteService := dailynote.NewService(repo)
	dailyNoteAppService := dailynoteapp.NewDailyNoteApplicationService(dailyNoteService)

	// 2. 从上下文中获取用户信息（由认证中间件设置）
	user, ok := middleware.GetDataFromContext(ctx)
	if !ok {
		return response.DailyNoteListResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 设置默认分页参数
	// 注意：当前实现不支持从查询参数获取page和pageSize，使用默认值
	page := 1
	pageSize := 10

	// 4. 调用应用服务获取笔记列表
	dailyNotePageDTO, err := dailyNoteAppService.GetDailyNoteList(ctx, user.UserID, page, pageSize)
	if err != nil {
		return response.DailyNoteListResponse{}, err
	}

	// 5. 转换为HTTP响应
	return response.ToDailyNoteListResponse(*dailyNotePageDTO), nil
}

// UpdateDailyNoteHandler 更新今日的每日笔记处理器
func UpdateDailyNoteHandler(ctx context.Context, req request.DailyNoteRequest) (response.DailyNoteResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewDailyNoteRepository()
	dailyNoteService := dailynote.NewService(repo)
	dailyNoteAppService := dailynoteapp.NewDailyNoteApplicationService(dailyNoteService)

	// 2. 从上下文中获取用户信息（由认证中间件设置）
	user, ok := middleware.GetDataFromContext(ctx)
	if !ok {
		return response.DailyNoteResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务更新今日笔记
	dailyNoteDTO, err := dailyNoteAppService.UpdateDailyNote(ctx, user.UserID, req.Content)
	if err != nil {
		return response.DailyNoteResponse{}, err
	}

	// 4. 转换为HTTP响应
	return response.ToDailyNoteResponse(*dailyNoteDTO), nil
}

// DeleteDailyNoteHandler 删除今日的每日笔记处理器
func DeleteDailyNoteHandler(ctx context.Context, req request.EmptyRequest) (response.MessageResponse, error) {
	// 1. 初始化服务层
	repo := mysql.NewDailyNoteRepository()
	dailyNoteService := dailynote.NewService(repo)
	dailyNoteAppService := dailynoteapp.NewDailyNoteApplicationService(dailyNoteService)

	// 2. 从上下文中获取用户信息（由认证中间件设置）
	user, ok := middleware.GetDataFromContext(ctx)
	if !ok {
		return response.MessageResponse{}, errors.New("unauthorized: invalid user context")
	}

	// 3. 调用应用服务删除今日笔记
	err := dailyNoteAppService.DeleteDailyNote(ctx, user.UserID)
	if err != nil {
		return response.MessageResponse{}, err
	}

	// 4. 返回成功消息
	return response.MessageResponse{
		Message: "每日笔记删除成功",
	}, nil
}
