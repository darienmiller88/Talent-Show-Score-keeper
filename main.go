package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	
	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "Hello, World!")
	})

	http.ListenAndServe(":8080", router)
}