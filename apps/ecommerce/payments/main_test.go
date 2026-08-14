package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPaymentApprovesPositiveAmount(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(`{"order_id":"order-1","amount":12.5}`))
	newHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", recorder.Code, recorder.Body.String())
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "approved" || body["order_id"] != "order-1" {
		t.Fatalf("body = %#v", body)
	}
}

func TestPaymentRejectsInvalidInput(t *testing.T) {
	for name, body := range map[string]string{
		"zero amount": `{"order_id":"order-1","amount":0}`,
		"empty order": `{"order_id":"","amount":10}`,
		"bad json":    `{`,
	} {
		t.Run(name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			newHandler().ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/payments", strings.NewReader(body)))
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", recorder.Code)
			}
		})
	}
}
