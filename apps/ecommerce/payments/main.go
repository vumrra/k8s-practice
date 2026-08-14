package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"k8s-practice/internal/web"
)

func newHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /payments", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			OrderID string  `json:"order_id"`
			Amount  float64 `json:"amount"`
		}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil || request.OrderID == "" || request.Amount <= 0 {
			web.JSON(w, http.StatusBadRequest, map[string]string{"error": "order_id and positive amount are required"})
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			web.JSON(w, http.StatusBadRequest, map[string]string{"error": "request body must contain one JSON object"})
			return
		}
		web.JSON(w, http.StatusOK, map[string]string{"order_id": request.OrderID, "status": "approved"})
	})
	return web.Health("payments", mux)
}

func main() {
	if err := web.Run(":8080", newHandler()); err != nil {
		log.Fatal(err)
	}
}
