package routes

import (
	"net/http"
	"todolist/internal/interfaces/http/handler"
)

func InitHealthRoute(mux *http.ServeMux) {
	mux.Handle("/health", handler.Wrap(handler.GetHealthHandler))
}
