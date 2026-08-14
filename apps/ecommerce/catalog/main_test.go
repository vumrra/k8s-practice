package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProductsReturnsTheSampleCatalog(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/products", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var body struct {
		Products []struct {
			ID string `json:"id"`
		} `json:"products"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body.Products) != 2 || body.Products[0].ID != "pencil" || body.Products[1].ID != "notebook" {
		t.Fatalf("products = %#v", body.Products)
	}
}

func TestProductsRejectsPost(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/products", nil))
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", recorder.Code)
	}
}
