package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOrderConfirmsAfterInventoryAndPaymentSucceed(t *testing.T) {
	inventory := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/inventory/pencil" {
			t.Fatalf("inventory path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"product_id":"pencil","available":true,"quantity":100}`))
	}))
	defer inventory.Close()
	payments := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order_id":"received","status":"approved"}`))
	}))
	defer payments.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"product_id":"pencil","quantity":2,"amount":3}`))
	newHandler(inventory.URL, payments.URL, &http.Client{Timeout: time.Second}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201: %s", recorder.Code, recorder.Body.String())
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "confirmed" || body["order_id"] == "" {
		t.Fatalf("body = %#v", body)
	}
}

func TestOrderStopsWhenInventoryIsUnavailable(t *testing.T) {
	inventory := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"product_id":"pencil","available":false,"quantity":0}`))
	}))
	defer inventory.Close()
	paymentCalled := false
	payments := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		paymentCalled = true
		w.WriteHeader(http.StatusOK)
	}))
	defer payments.Close()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"product_id":"pencil","quantity":1,"amount":1}`))
	newHandler(inventory.URL, payments.URL, &http.Client{Timeout: time.Second}).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", recorder.Code)
	}
	if paymentCalled {
		t.Fatal("payment was called for unavailable inventory")
	}
}

func TestOrderConvertsDependencyFailureAndTimeout(t *testing.T) {
	for name, dependency := range map[string]struct {
		dependency http.Handler
		want       int
	}{
		"server error": {http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) }), http.StatusBadGateway},
		"timeout": {http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(50 * time.Millisecond)
			_, _ = w.Write([]byte(`{"available":true}`))
		}), http.StatusGatewayTimeout},
	} {
		t.Run(name, func(t *testing.T) {
			inventory := httptest.NewServer(dependency.dependency)
			defer inventory.Close()
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"product_id":"pencil","quantity":1,"amount":1}`))
			newHandler(inventory.URL, "http://unused", &http.Client{Timeout: 10 * time.Millisecond}).ServeHTTP(recorder, request)
			if recorder.Code != dependency.want {
				t.Fatalf("status = %d, want %d", recorder.Code, dependency.want)
			}
		})
	}
}

func TestOrderRejectsMalformedInput(t *testing.T) {
	recorder := httptest.NewRecorder()
	newHandler("http://unused", "http://unused", &http.Client{Timeout: time.Second}).ServeHTTP(
		recorder,
		httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"product_id":"","quantity":0,"amount":0}`)),
	)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}
