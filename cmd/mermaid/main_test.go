package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func resetFlags(args ...string) {
	flag.CommandLine = flag.NewFlagSet("mermaid", flag.ContinueOnError)
	os.Args = append([]string{"mermaid"}, args...)
}

func TestRunFileToFile(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "d.mmd")
	if err := os.WriteFile(in, []byte("graph TD\nA --> B"), 0o644); err != nil {
		t.Fatal(err)
	}
	out := filepath.Join(dir, "d.svg")
	resetFlags("-theme", "dark", "-o", out, in)

	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(got), "<svg") {
		t.Errorf("output not SVG: %.20q", got)
	}
}

func TestRunStdinToStdout(t *testing.T) {
	dir := t.TempDir()
	stdin := filepath.Join(dir, "in")
	if err := os.WriteFile(stdin, []byte("graph LR\nA --> B"), 0o644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(stdin)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	oldIn, oldOut := os.Stdin, os.Stdout
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()
	os.Stdin = f

	outFile := filepath.Join(dir, "captured")
	w, err := os.Create(outFile)
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	resetFlags("-")
	runErr := run()
	w.Close()
	os.Stdout = oldOut
	if runErr != nil {
		t.Fatalf("run: %v", runErr)
	}

	got, _ := os.ReadFile(outFile)
	if !strings.Contains(string(got), "<svg") {
		t.Errorf("stdout missing SVG: %.20q", got)
	}
}

func TestRunBadSource(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "bad.mmd")
	if err := os.WriteFile(in, []byte("graph TD\nA[oops"), 0o644); err != nil {
		t.Fatal(err)
	}
	resetFlags(in)
	if err := run(); err == nil {
		t.Error("expected error for bad source, got nil")
	}
}

func TestRunMissingFile(t *testing.T) {
	resetFlags(filepath.Join(t.TempDir(), "nope.mmd"))
	if err := run(); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
