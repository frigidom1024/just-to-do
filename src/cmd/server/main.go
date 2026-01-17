package main

import (
	"fmt"
	"net/http"
	"os"
	"todolist/internal/interfaces/http/handler"
)

func main() {
	fmt.Println("Hello, world!")
	http.Handle("/health", handler.Wrap(handler.GetHealthHandler))
	http.ListenAndServe(":8080", nil)
	os.Exit(0)
}
