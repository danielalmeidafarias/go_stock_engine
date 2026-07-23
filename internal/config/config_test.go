package config

import "testing"

func TestEnvironmentConfig(t *testing.T) {
	t.Setenv("DATABASE_DRIVER", "POSTGRES")
	t.Setenv("DATABASE_URL", "postgres://postgres:example@localhost:5432/stock?sslmode=disable")
	t.Setenv("HANDLER_TYPE", "HTTP")
	t.Setenv("PAGINATION_DEFAULT_LIMIT", "20")
	t.Setenv("PAGINATION_MAX_LIMIT", "100")
	t.Setenv("PRIORITY_USE_POLICY", "false")
	t.Setenv("PRIORITY_NEGATIVE_STOCK_FACTOR", "1.0")
	t.Setenv("PRIORITY_LEAD_TIME_FACTOR", "1.0")
	t.Setenv("PRIORITY_ZERO_SALES_FACTOR", "1.0")

	configuration := EnvironmentConfig{}

	if got := configuration.Database().Driver; got != "POSTGRES" {
		t.Errorf("Database().Driver = %q, want POSTGRES", got)
	}
	if got := configuration.Pagination().MaxLimit; got != "100" {
		t.Errorf("Pagination().MaxLimit = %q, want 100", got)
	}
	if got := configuration.PriorityPolicy().LeadTimeFactor; got != "1.0" {
		t.Errorf("PriorityPolicy().LeadTimeFactor = %q, want 1.0", got)
	}
	if got := configuration.Database().ConnectionString; got != "postgres://postgres:example@localhost:5432/stock?sslmode=disable" {
		t.Errorf("Database().ConnectionString = %q", got)
	}
}
