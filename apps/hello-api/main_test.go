package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootDescribesTheRunningPod(t *testing.T) {
	env := map[string]string{
		"APP_VERSION": "dev",
		"HOSTNAME":    "test-pod",
		"MESSAGE":     "hello",
	}
	recorder := httptest.NewRecorder()

	newHandler(func(key string) string { return env[key] }).ServeHTTP(
		recorder,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"service": "hello-api", "version": "dev", "pod": "test-pod"}
	for key, value := range want {
		if body[key] != value {
			t.Fatalf("%s = %q, want %q", key, body[key], value)
		}
	}
}

func TestConfigReturnsOnlyPublicMessage(t *testing.T) {
	env := map[string]string{"MESSAGE": "practice", "SECRET": "do-not-leak"}
	recorder := httptest.NewRecorder()

	newHandler(func(key string) string { return env[key] }).ServeHTTP(
		recorder,
		httptest.NewRequest(http.MethodGet, "/config", nil),
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	if got := recorder.Body.String(); !strings.Contains(got, "practice") || strings.Contains(got, "do-not-leak") {
		t.Fatalf("unexpected body: %s", got)
	}
}
