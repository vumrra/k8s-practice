package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"k8s-practice/internal/web"
)

func newHandler(catalogURL, ordersURL string) http.Handler {
	catalogTarget, catalogErr := url.Parse(catalogURL)
	ordersTarget, ordersErr := url.Parse(ordersURL)
	if catalogErr != nil || ordersErr != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			web.JSON(w, http.StatusInternalServerError, map[string]string{"error": "invalid service URL"})
		})
	}

	mux := http.NewServeMux()
	mux.Handle("GET /api/products", http.StripPrefix("/api", httputil.NewSingleHostReverseProxy(catalogTarget)))
	mux.Handle("POST /api/orders", http.StripPrefix("/api", httputil.NewSingleHostReverseProxy(ordersTarget)))
	return web.Health("gateway", mux)
}

func main() {
	catalogURL := os.Getenv("CATALOG_URL")
	if catalogURL == "" {
		catalogURL = "http://catalog:8080"
	}
	ordersURL := os.Getenv("ORDERS_URL")
	if ordersURL == "" {
		ordersURL = "http://orders:8080"
	}
	if err := web.Run(":8080", newHandler(catalogURL, ordersURL)); err != nil {
		log.Fatal(err)
	}
}
