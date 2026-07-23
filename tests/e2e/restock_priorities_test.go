package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const defaultBaseURL = "http://localhost:8081"

type productFixture struct {
	name              string
	category          string
	currentStock      int
	minimumStock      int
	averageDailySales int
	leadTimeDays      int
	unitCost          float64
	criticalityLevel  int
}

type priorityExpectation struct {
	productFixture
	projectedStock int
	urgencyScore   int
}

type createProductResponse struct {
	ID *string `json:"id"`
}

type restockPrioritiesResponse struct {
	Priorities *[]restockPriorityResponse `json:"priorities"`
}

type restockPriorityResponse struct {
	PartID         *string `json:"partId"`
	Name           *string `json:"name"`
	CurrentStock   *int    `json:"currentStock"`
	ProjectedStock *int    `json:"projectedStock"`
	MinimumStock   *int    `json:"minimumStock"`
	UrgencyScore   *int    `json:"urgencyScore"`
}

func TestMain(m *testing.M) {
	baseURL := e2eBaseURL()
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		response, err := http.Get(baseURL + "/stock?page=1&limit=1")
		if err == nil {
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				os.Exit(m.Run())
			}
		}
		time.Sleep(time.Second)
	}

	fmt.Fprintf(os.Stderr, "API not available at %s\n", baseURL)
	os.Exit(1)
}

func TestRestockPrioritiesContract(t *testing.T) {
	priorities := []priorityExpectation{
		// (20 - (10 - (5 * 2))) * 5 = 100
		{productFixture: productFixture{"Alpha Filter", "e2e", 10, 20, 5, 2, 18.50, 5}, projectedStock: 0, urgencyScore: 100},
		// Same score, criticality and sales as Alpha; name is the final tie-breaker.
		{productFixture: productFixture{"Bravo Filter", "e2e", 10, 20, 5, 2, 18.50, 5}, projectedStock: 0, urgencyScore: 100},
		// (10 - (10 - (4 * 5))) * 5 = 100
		{productFixture: productFixture{"Charlie Belt", "e2e", 10, 10, 4, 5, 25.00, 5}, projectedStock: -10, urgencyScore: 100},
		// (20 - (10 - (5 * 3))) * 4 = 100
		{productFixture: productFixture{"Delta Bearing", "e2e", 10, 20, 5, 3, 40.00, 4}, projectedStock: -5, urgencyScore: 100},
		// (20 - (5 - (3 * 5))) * 3 = 90
		{productFixture: productFixture{"Echo Pad", "e2e", 5, 20, 3, 5, 42.00, 3}, projectedStock: -10, urgencyScore: 90},
		// (10 - (5 - (2 * 2))) * 3 = 27
		{productFixture: productFixture{"Foxtrot Lamp", "e2e", 5, 10, 2, 2, 12.00, 3}, projectedStock: 1, urgencyScore: 27},
	}
	notNeeded := []productFixture{
		{"Golf Hose", "e2e", 30, 10, 2, 2, 15.00, 2},
		{"Hotel Sensor", "e2e", 20, 20, 0, 5, 28.00, 3},
		{"India Relay", "e2e", 10, 0, 1, 5, 35.00, 2},
		{"Juliet Kit", "e2e", 40, 20, 3, 3, 80.00, 4},
	}

	createdIDs := make(map[string]string, len(priorities)+len(notNeeded))
	for _, priority := range priorities {
		createdIDs[priority.name] = createProduct(t, priority.productFixture)
	}
	for _, product := range notNeeded {
		createdIDs[product.name] = createProduct(t, product)
	}
	t.Cleanup(func() {
		for _, id := range createdIDs {
			deleteProduct(t, id)
		}
	})

	response := getPriorities(t)
	if response.Priorities == nil {
		t.Fatal("response is missing priorities")
	}
	if got := len(*response.Priorities); got != len(priorities) {
		t.Fatalf("priorities count: got %d, want %d", got, len(priorities))
	}

	for index, expected := range priorities {
		actual := (*response.Priorities)[index]
		assertPriority(t, actual, expected, createdIDs[expected.name])
	}
}

func createProduct(t *testing.T, product productFixture) string {
	t.Helper()
	payload, err := json.Marshal(struct {
		Name              string  `json:"name"`
		Category          string  `json:"category"`
		CurrentStock      int     `json:"current_stock"`
		MinimumStock      int     `json:"minimum_stock"`
		AverageDailySales int     `json:"average_daily_sales"`
		LeadTimeDays      int     `json:"lead_time_days"`
		UnitCost          float64 `json:"unit_cost"`
		CriticalityLevel  int     `json:"criticality_level"`
	}{
		Name:              product.name,
		Category:          product.category,
		CurrentStock:      product.currentStock,
		MinimumStock:      product.minimumStock,
		AverageDailySales: product.averageDailySales,
		LeadTimeDays:      product.leadTimeDays,
		UnitCost:          product.unitCost,
		CriticalityLevel:  product.criticalityLevel,
	})
	if err != nil {
		t.Fatalf("encode product: %v", err)
	}

	response, body := request(t, http.MethodPost, "/stock", bytes.NewReader(payload))
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("create %q: got %d, want 201: %s", product.name, response.StatusCode, body)
	}

	var created createProductResponse
	if err := json.Unmarshal(body, &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == nil || *created.ID == "" {
		t.Fatalf("create %q: response is missing id", product.name)
	}
	return *created.ID
}

func getPriorities(t *testing.T) restockPrioritiesResponse {
	t.Helper()
	response, body := request(t, http.MethodGet, "/restock/priorities?page=1&limit=100", nil)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("get priorities: got %d, want 200: %s", response.StatusCode, body)
	}

	var priorities restockPrioritiesResponse
	if err := json.Unmarshal(body, &priorities); err != nil {
		t.Fatalf("decode priorities response: %v", err)
	}
	return priorities
}

func assertPriority(t *testing.T, actual restockPriorityResponse, expected priorityExpectation, expectedID string) {
	t.Helper()
	if actual.PartID == nil || actual.Name == nil || actual.CurrentStock == nil || actual.ProjectedStock == nil || actual.MinimumStock == nil || actual.UrgencyScore == nil {
		t.Fatalf("priority response has missing fields: %#v", actual)
	}
	if *actual.PartID != expectedID {
		t.Errorf("partId: got %q, want %q", *actual.PartID, expectedID)
	}
	if *actual.Name != expected.name {
		t.Errorf("name: got %q, want %q", *actual.Name, expected.name)
	}
	if *actual.CurrentStock != expected.currentStock {
		t.Errorf("currentStock: got %d, want %d", *actual.CurrentStock, expected.currentStock)
	}
	if *actual.ProjectedStock != expected.projectedStock {
		t.Errorf("projectedStock: got %d, want %d", *actual.ProjectedStock, expected.projectedStock)
	}
	if *actual.MinimumStock != expected.minimumStock {
		t.Errorf("minimumStock: got %d, want %d", *actual.MinimumStock, expected.minimumStock)
	}
	if *actual.UrgencyScore != expected.urgencyScore {
		t.Errorf("urgencyScore: got %d, want %d", *actual.UrgencyScore, expected.urgencyScore)
	}
}

func deleteProduct(t *testing.T, id string) {
	t.Helper()
	response, body := request(t, http.MethodDelete, "/stock/"+id, nil)
	if response.StatusCode != http.StatusNoContent {
		t.Errorf("delete %q: got %d, want 204: %s", id, response.StatusCode, body)
	}
}

func request(t *testing.T, method, path string, payload io.Reader) (*http.Response, []byte) {
	t.Helper()
	req, err := http.NewRequest(method, e2eBaseURL()+path, payload)
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read %s %s response: %v", method, path, err)
	}
	return response, body
}

func e2eBaseURL() string {
	if baseURL := os.Getenv("E2E_BASE_URL"); baseURL != "" {
		return baseURL
	}
	return defaultBaseURL
}
