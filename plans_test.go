package main

import "testing"

func TestGetPlan(t *testing.T) {
	plan := GetPlan("standard-yanari")
	if plan.ConnectionLimit != 60 {
		t.Fatalf("epxected 60 for standard-yanari, got %v", plan.ConnectionLimit)
	}
}
