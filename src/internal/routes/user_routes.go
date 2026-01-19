package routes

import (
	"net/http"
	"todolist/internal/interfaces/http/handler"
	"todolist/internal/interfaces/http/middleware"
)

func InitUserRoute(mux *http.ServeMux) {
	authmiddle := middleware.GetAuthMiddleware()
	mux.Handle("/api/v1/users/login", handler.Wrap(handler.LoginUserHandler))

	// 用户路由
	mux.Handle("/api/v1/users/register", handler.Wrap(handler.RegisterUserHandler))
	mux.Handle("/api/v1/users/password", authmiddle.Authenticate(handler.Wrap(handler.ChangePasswordHandler)))
	mux.Handle("/api/v1/users/email", authmiddle.Authenticate(handler.Wrap(handler.UpdateEmailHandler)))
	mux.Handle("/api/v1/users/avatar", authmiddle.Authenticate(handler.Wrap(handler.UpdateAvatarHandler)))
}
