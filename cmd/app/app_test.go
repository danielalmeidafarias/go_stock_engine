package main

import "testing"

func TestNewPriorityPolicy(t *testing.T) {
	t.Run("disabled allows omitted factors", func(t *testing.T) {
		if policy := NewPriorityPolicy(false, "", "", ""); policy.UsePolicy {
			t.Fatal("expected disabled policy")
		}
	})

	t.Run("enabled parses factors", func(t *testing.T) {
		policy := NewPriorityPolicy(true, "1.5", "1.2", "0.5")
		if !policy.UsePolicy || policy.NegativeStockFactor != 1.5 || policy.LeadTimeFactor != 1.2 || policy.ZeroSalesFactor != 0.5 {
			t.Fatalf("unexpected policy: %#v", policy)
		}
	})
}
