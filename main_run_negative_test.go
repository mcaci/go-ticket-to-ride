package main

import "testing"

func TestRun_MissingRoutes_ReturnsError(t *testing.T) {
	args := []string{
		"-show=false",
		"-mode=tickets",
		"-routes=./nonexistent/routes.csv",
		"-cities=./pkg/data/USA/cities.csv",
		"-tickets=./pkg/data/USA/tickets.csv",
	}
	if err := Run(args); err == nil {
		t.Fatalf("expected error for missing routes file, got nil")
	}
}

func TestRun_UnknownMode_ReturnsError(t *testing.T) {
	args := []string{
		"-show=false",
		"-mode=notamode",
		"-routes=./pkg/data/USA/routes.csv",
		"-cities=./pkg/data/USA/cities.csv",
		"-tickets=./pkg/data/USA/tickets.csv",
	}
	if err := Run(args); err == nil {
		t.Fatalf("expected error for unknown mode, got nil")
	}
}

func TestRun_BadFlagParsing_ReturnsError(t *testing.T) {
	args := []string{
		"-show=false",
		"-frame-delay=notanint",
	}
	if err := Run(args); err == nil {
		t.Fatalf("expected flag parse error, got nil")
	}
}
