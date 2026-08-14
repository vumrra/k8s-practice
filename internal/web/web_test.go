package web

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunReturnsAListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	if err := Run(listener.Addr().String(), http.NotFoundHandler()); err == nil {
		t.Fatal("Run returned nil for an occupied address")
	}
}

func TestJSONWritesStatusContentTypeAndBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	JSON(recorder, http.StatusCreated, map[string]string{"result": "created"})

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}
	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["result"] != "created" {
		t.Fatalf("result = %q, want created", body["result"])
	}
}

func TestHealthServesBothProbeEndpoints(t *testing.T) {
	handler := Health("hello-api", http.NotFoundHandler())

	for _, path := range []string{"/healthz", "/readyz"} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", recorder.Code)
			}
			var body map[string]string
			if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body["status"] != "ok" || body["service"] != "hello-api" {
				t.Fatalf("body = %#v", body)
			}
		})
	}
}

func TestHealthPassesApplicationRequestsToNextHandler(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})
	recorder := httptest.NewRecorder()

	Health("hello-api", next).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
}
