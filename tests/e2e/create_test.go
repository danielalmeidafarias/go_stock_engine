package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

func runCreate(t *testing.T) {
	ids := make([]string, SEED_COUNT)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, CONCURRENCY)

	t.Logf("Creating %d products with concurrency %d...", SEED_COUNT, CONCURRENCY)
	start := time.Now()

	for i := range SEED_COUNT {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			body := createRequest{
				Name:              fmt.Sprintf("Product-%06d", idx),
				Category:          categories[rand.Intn(len(categories))],
				CurrentStock:      rand.Intn(500) - 50,
				MinimumStock:      rand.Intn(100) + 10,
				AverageDailySales: rand.Intn(30),
				LeadTimeDays:      rand.Intn(60),
				UnitCost:          float64(rand.Intn(1000)+1) / 10.0,
				CriticalityLevel:  rand.Intn(5) + 1,
			}

			data, _ := json.Marshal(body)
			resp, err := http.Post(baseURL+"/stock", "application/json", bytes.NewReader(data))
			if err != nil {
				t.Errorf("create[%d]: POST error: %v", idx, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("create[%d]: expected 201, got %d", idx, resp.StatusCode)
				return
			}

			var cr createResponse
			if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
				t.Errorf("create[%d]: unmarshal: %v", idx, err)
				return
			}

			mu.Lock()
			ids[idx] = cr.ID
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("CREATE: %d requests in %v (%.0f products/sec)", SEED_COUNT, elapsed, float64(SEED_COUNT)/elapsed.Seconds())

	var validIDs []string
	for _, id := range ids {
		if id != "" {
			validIDs = append(validIDs, id)
		}
	}

	t.Logf("Successfully created %d/%d products", len(validIDs), SEED_COUNT)

	if len(validIDs) < SEED_COUNT {
		t.Fatalf("Some create failures: only %d/%d succeeded", len(validIDs), SEED_COUNT)
	}

	seededIDs = validIDs
}
