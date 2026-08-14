package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayRoutesOnlyDocumentedAPIs(t *testing.T) {
	catalog := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/products" {
			t.Fatalf("catalog path = %s, want /products", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"products":[]}`))
	}))
	defer catalog.Close()
	orders := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders" {
			t.Fatalf("orders path = %s, want /orders", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer orders.Close()
	handler := newHandler(catalog.URL, orders.URL)

	for _, test := range []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodGet, "/api/products", http.StatusOK},
		{http.MethodPost, "/api/orders", http.StatusCreated},
		{http.MethodGet, "/api/unknown", http.StatusNotFound},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(test.method, test.path, nil))
		if recorder.Code != test.want {
			t.Fatalf("%s %s status = %d, want %d", test.method, test.path, recorder.Code, test.want)
		}
	}
}
