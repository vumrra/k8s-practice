package main

import (
	"log"
	"net/http"

	"k8s-practice/internal/web"
)

func newHandler() http.Handler {
	stock := map[string]int{"pencil": 100, "notebook": 25}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /inventory/{productID}", func(w http.ResponseWriter, r *http.Request) {
		productID := r.PathValue("productID")
		quantity, exists := stock[productID]
		if !exists {
			web.JSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
			return
		}
		web.JSON(w, http.StatusOK, map[string]any{
			"product_id": productID,
			"available":  quantity > 0,
			"quantity":   quantity,
		})
	})
	return web.Health("inventory", mux)
}

func main() {
	if err := web.Run(":8080", newHandler()); err != nil {
		log.Fatal(err)
	}
}
