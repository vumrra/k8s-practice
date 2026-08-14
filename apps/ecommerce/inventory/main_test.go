package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInventoryReturnsKnownStock(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/inventory/pencil", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var body struct {
		ProductID string `json:"product_id"`
		Available bool   `json:"available"`
		Quantity  int    `json:"quantity"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.ProductID != "pencil" || !body.Available || body.Quantity != 100 {
		t.Fatalf("inventory = %#v", body)
	}
}

func TestInventoryReturnsNotFoundForUnknownProduct(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/inventory/unknown", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", recorder.Code)
	}
}
