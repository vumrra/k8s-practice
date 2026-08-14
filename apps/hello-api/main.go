package main

import (
	"log"
	"net/http"
	"os"

	"k8s-practice/internal/web"
)

func newHandler(getenv func(string) string) http.Handler {
	value := func(key, fallback string) string {
		if result := getenv(key); result != "" {
			return result
		}
		return fallback
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		web.JSON(w, http.StatusOK, map[string]string{
			"service": "hello-api",
			"version": value("APP_VERSION", "dev"),
			"pod":     value("HOSTNAME", "local"),
		})
	})
	mux.HandleFunc("GET /config", func(w http.ResponseWriter, _ *http.Request) {
		web.JSON(w, http.StatusOK, map[string]string{"message": value("MESSAGE", "hello")})
	})
	return web.Health("hello-api", mux)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := web.Run(":"+port, newHandler(os.Getenv)); err != nil {
		log.Fatal(err)
	}
}
