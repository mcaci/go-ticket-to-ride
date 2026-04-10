package main_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestMainBinaryRuns_NoShow(t *testing.T) {
	// build binary
	tmpDir := t.TempDir()
	bin := filepath.Join(tmpDir, "ttr_bin")
	build := exec.Command("go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build failed: %v, output: %s", err, string(out))
	}

	// run binary with show=false and explicit data paths; use timeout to avoid hangs
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin,
		"-show=false",
		"-routes=./pkg/data/USA/routes.csv",
		"-cities=./pkg/data/USA/cities.csv",
		"-tickets=./pkg/data/USA/tickets.csv",
		"-mode=tickets",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("running binary failed: %v, output: %s", err, string(out))
	}
}
