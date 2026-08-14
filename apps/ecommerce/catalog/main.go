package main

import (
	"log"
	"net/http"

	"k8s-practice/internal/web"
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", func(w http.ResponseWriter, _ *http.Request) {
		web.JSON(w, http.StatusOK, map[string]any{"products": []map[string]any{
			{"id": "pencil", "name": "Pencil", "price": 1.5},
			{"id": "notebook", "name": "Notebook", "price": 5.0},
		}})
	})
	return web.Health("catalog", mux)
}

func main() {
	if err := web.Run(":8080", newHandler()); err != nil {
		log.Fatal(err)
	}
}
