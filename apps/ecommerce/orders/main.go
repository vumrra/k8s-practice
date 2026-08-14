package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"k8s-practice/internal/web"
)

var orderSequence atomic.Uint64

func newHandler(inventoryURL, paymentsURL string, client *http.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ProductID string  `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Amount    float64 `json:"amount"`
		}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil || request.ProductID == "" || request.Quantity <= 0 || request.Amount <= 0 {
			web.JSON(w, http.StatusBadRequest, map[string]string{"error": "product_id, positive quantity, and positive amount are required"})
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			web.JSON(w, http.StatusBadRequest, map[string]string{"error": "request body must contain one JSON object"})
			return
		}

		inventoryResponse, err := client.Get(inventoryURL + "/inventory/" + url.PathEscape(request.ProductID))
		if err != nil {
			dependencyError(w, err)
			return
		}
		defer inventoryResponse.Body.Close()
		if inventoryResponse.StatusCode != http.StatusOK {
			web.JSON(w, http.StatusBadGateway, map[string]string{"error": "inventory dependency failed"})
			return
		}
		var inventory struct {
			Available bool `json:"available"`
			Quantity  int  `json:"quantity"`
		}
		if err := json.NewDecoder(inventoryResponse.Body).Decode(&inventory); err != nil {
			web.JSON(w, http.StatusBadGateway, map[string]string{"error": "invalid inventory response"})
			return
		}
		if !inventory.Available || inventory.Quantity < request.Quantity {
			web.JSON(w, http.StatusConflict, map[string]string{"error": "insufficient inventory"})
			return
		}

		orderID := fmt.Sprintf("order-%d", orderSequence.Add(1))
		paymentBody, _ := json.Marshal(map[string]any{"order_id": orderID, "amount": request.Amount})
		paymentRequest, _ := http.NewRequestWithContext(r.Context(), http.MethodPost, paymentsURL+"/payments", bytes.NewReader(paymentBody))
		paymentRequest.Header.Set("Content-Type", "application/json")
		paymentResponse, err := client.Do(paymentRequest)
		if err != nil {
			dependencyError(w, err)
			return
		}
		defer paymentResponse.Body.Close()
		if paymentResponse.StatusCode != http.StatusOK {
			web.JSON(w, http.StatusBadGateway, map[string]string{"error": "payments dependency failed"})
			return
		}
		var payment struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(paymentResponse.Body).Decode(&payment); err != nil || payment.Status != "approved" {
			web.JSON(w, http.StatusBadGateway, map[string]string{"error": "payment was not approved"})
			return
		}
		web.JSON(w, http.StatusCreated, map[string]string{"order_id": orderID, "status": "confirmed"})
	})
	return web.Health("orders", mux)
}

func dependencyError(w http.ResponseWriter, err error) {
	var networkError net.Error
	if errors.As(err, &networkError) && networkError.Timeout() {
		web.JSON(w, http.StatusGatewayTimeout, map[string]string{"error": "dependency timed out"})
		return
	}
	web.JSON(w, http.StatusBadGateway, map[string]string{"error": "dependency unavailable"})
}

func main() {
	inventoryURL := os.Getenv("INVENTORY_URL")
	if inventoryURL == "" {
		inventoryURL = "http://inventory:8080"
	}
	paymentsURL := os.Getenv("PAYMENTS_URL")
	if paymentsURL == "" {
		paymentsURL = "http://payments:8080"
	}
	if err := web.Run(":8080", newHandler(inventoryURL, paymentsURL, &http.Client{Timeout: 2 * time.Second})); err != nil {
		log.Fatal(err)
	}
}
