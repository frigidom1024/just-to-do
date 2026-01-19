package routes

import (
	"net/http"

	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/middleware"
)

// InitDailyNoteRoute 初始化每日笔记路由
func InitDailyNoteRoute(mux *http.ServeMux) {
	// 获取认证中间件
	authmiddle := middleware.GetAuthMiddleware()

	// 每日笔记路由，所有路由都需要认证
	// 创建每日笔记
	mux.Handle("/api/v1/daily-notes", authmiddle.Authenticate(handler.Wrap(handler.CreateDailyNoteHandler)))
	// 获取今日每日笔记
	mux.Handle("/api/v1/daily-notes/today", authmiddle.Authenticate(handler.Wrap(handler.GetTodayDailyNoteHandler)))
	// 分页获取每日笔记列表
	mux.Handle("/api/v1/daily-notes/list", authmiddle.Authenticate(handler.Wrap(handler.GetDailyNoteListHandler)))
	// 更新今日每日笔记
	mux.Handle("/api/v1/daily-notes/today/update", authmiddle.Authenticate(handler.Wrap(handler.UpdateDailyNoteHandler)))
	// 删除今日每日笔记
	mux.Handle("/api/v1/daily-notes/today/delete", authmiddle.Authenticate(handler.Wrap(handler.DeleteDailyNoteHandler)))
}
