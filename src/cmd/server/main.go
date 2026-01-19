package main

import (
	"fmt"
	"net/http"
	"os"

	"todolist/internal/routes"
)

func main() {
	fmt.Println("Starting Todo List Server on :8080...")

	// Initialize HTTP server
	mux := http.NewServeMux()
	routes.InitUserRoute(mux)
	routes.InitHealthRoute(mux)
	// Setup routes and middleware

	// Start server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
