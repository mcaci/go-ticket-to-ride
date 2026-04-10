package main

import "testing"

func TestRun_NoShow(t *testing.T) {
	args := []string{
		"-show=false",
		"-routes=./pkg/data/USA/routes.csv",
		"-cities=./pkg/data/USA/cities.csv",
		"-tickets=./pkg/data/USA/tickets.csv",
		"-mode=tickets",
	}
	if err := Run(args); err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
}
