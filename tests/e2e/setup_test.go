package e2e

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

const SEED_COUNT = 10_000
const CONCURRENCY = 50

var seededIDs []string

func TestMain(m *testing.M) {
	resp, err := http.Get(baseURL + "/stock?page=1&limit=1")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "API not available at %s. Is the app running?\n", baseURL)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	t.Run("Create", runCreate)
	if len(seededIDs) == 0 {
		t.Fatal("Create phase failed, cannot continue")
	}

	t.Run("GetByID", runGetByID)
	t.Run("GetAllProductStock", runGetAllProductStock)
	t.Run("GetByCategory", runGetByCategory)
	t.Run("Update", runUpdate)
	t.Run("RestockPriorities", runRestockPriorities)
	t.Run("Delete", runDelete)
}
