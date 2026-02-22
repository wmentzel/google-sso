package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"google-sso/auth"
	"google-sso/handlers"
	"google-sso/middleware"
)

func main() {
	// Load .env if present (ignored in production where env vars are set externally)
	_ = godotenv.Load()

	auth.Init()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/callback", handlers.Callback)
	mux.HandleFunc("/logout", handlers.Logout)

	// Root redirect → /hello
	mux.Handle("/", http.RedirectHandler("/hello", http.StatusFound))

	// Protected routes
	mux.Handle("/hello", middleware.RequireAuth(http.HandlerFunc(handlers.Hello)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
