package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:8080"

var httpClient = &http.Client{Timeout: 10 * time.Second}

type createRequest struct {
	Name              string  `json:"name"`
	Category          string  `json:"category"`
	CurrentStock      int     `json:"current_stock"`
	MinimumStock      int     `json:"minimum_stock"`
	AverageDailySales int     `json:"average_daily_sales"`
	LeadTimeDays      int     `json:"lead_time_days"`
	UnitCost          float64 `json:"unit_cost"`
	CriticalityLevel  int     `json:"criticality_level"`
}

type createResponse struct {
	ID string `json:"id"`
}

type updateRequest struct {
	CurrentStock     *int `json:"current_stock,omitempty"`
	CriticalityLevel *int `json:"criticality_level,omitempty"`
}

func BenchmarkCreateProductStock(b *testing.B) {
	requireAPI(b)
	ids := make([]string, 0, b.N)
	b.Cleanup(func() {
		for _, id := range ids {
			deleteProduct(b, id)
		}
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ids = append(ids, createProduct(b, i))
	}
}

func BenchmarkGetProductStockByID(b *testing.B) {
	requireAPI(b)
	id := createProduct(b, 0)
	b.Cleanup(func() { deleteProduct(b, id) })

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request(b, http.MethodGet, "/stock/"+id, nil, http.StatusOK)
	}
}

func BenchmarkGetAllProductStock(b *testing.B) {
	requireAPI(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request(b, http.MethodGet, "/stock?page=1&limit=100", nil, http.StatusOK)
	}
}

func BenchmarkGetProductStockByCategory(b *testing.B) {
	requireAPI(b)
	id := createProduct(b, 0)
	b.Cleanup(func() { deleteProduct(b, id) })

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request(b, http.MethodGet, "/stock/category/benchmark?page=1&limit=100", nil, http.StatusOK)
	}
}

func BenchmarkUpdateProductStock(b *testing.B) {
	requireAPI(b)
	id := createProduct(b, 0)
	b.Cleanup(func() { deleteProduct(b, id) })

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		currentStock := i
		payload, err := json.Marshal(updateRequest{CurrentStock: &currentStock})
		if err != nil {
			b.Fatalf("encode update request: %v", err)
		}
		request(b, http.MethodPut, "/stock/"+id, bytes.NewReader(payload), http.StatusNoContent)
	}
}

func BenchmarkDeleteProductStock(b *testing.B) {
	requireAPI(b)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		id := createProduct(b, i)
		b.StartTimer()

		deleteProduct(b, id)
	}
}

func BenchmarkRestockPriorities(b *testing.B) {
	requireAPI(b)
	id := createProduct(b, 0)
	b.Cleanup(func() { deleteProduct(b, id) })

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request(b, http.MethodGet, "/restock/priorities?page=1&limit=100", nil, http.StatusOK)
	}
}

func requireAPI(b *testing.B) {
	b.Helper()
	request(b, http.MethodGet, "/stock?page=1&limit=1", nil, http.StatusOK)
}

func createProduct(b *testing.B, index int) string {
	b.Helper()
	payload, err := json.Marshal(createRequest{
		Name:              fmt.Sprintf("Benchmark Product %d", index),
		Category:          "benchmark",
		CurrentStock:      10,
		MinimumStock:      20,
		AverageDailySales: 4,
		LeadTimeDays:      5,
		UnitCost:          18.50,
		CriticalityLevel:  3,
	})
	if err != nil {
		b.Fatalf("encode create request: %v", err)
	}

	body := request(b, http.MethodPost, "/stock", bytes.NewReader(payload), http.StatusCreated)
	var response createResponse
	if err := json.Unmarshal(body, &response); err != nil {
		b.Fatalf("decode create response: %v", err)
	}
	if response.ID == "" {
		b.Fatal("create response has an empty id")
	}
	return response.ID
}

func deleteProduct(b *testing.B, id string) {
	b.Helper()
	request(b, http.MethodDelete, "/stock/"+id, nil, http.StatusNoContent)
}

func request(b *testing.B, method, path string, payload io.Reader, expectedStatus int) []byte {
	b.Helper()
	req, err := http.NewRequest(method, baseURL+path, payload)
	if err != nil {
		b.Fatalf("create %s request: %v", method, err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := httpClient.Do(req)
	if err != nil {
		b.Fatalf("%s %s: %v", method, path, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		b.Fatalf("read %s %s response: %v", method, path, err)
	}
	if response.StatusCode != expectedStatus {
		b.Fatalf("%s %s: got %d, want %d: %s", method, path, response.StatusCode, expectedStatus, body)
	}
	return body
}
