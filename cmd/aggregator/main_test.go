package main

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunValidatesArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing input", args: []string{"--output", "results"}, want: "--input is required"},
		{name: "missing output", args: []string{"--input", "input.csv"}, want: "--output is required"},
		{name: "positional argument", args: []string{"--input", "input.csv", "--output", "results", "extra"}, want: "unexpected positional arguments"},
		{name: "unknown flag", args: []string{"--unknown"}, want: "flag provided but not defined"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := run(test.args)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("run() error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestRunCreatesAndClosesCPUProfile(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.csv")
	writeTestFile(t, input, "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,0.123456789,1\n")
	profile := filepath.Join(root, "cpu.pprof")
	if err := run([]string{"--input", input, "--output", filepath.Join(root, "results"), "--cpuprofile", profile}); err != nil {
		t.Fatalf("run() error = %v", err)
	}
	assertGzipProfile(t, profile)
}

func TestRunClosesCPUProfileWhenApplicationFails(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "invalid.csv")
	writeTestFile(t, input, "campaign_id,date,impressions,clicks,spend,conversions\nA,invalid-date,1,1,0.1,1\n")
	profile := filepath.Join(root, "cpu.pprof")
	err := run([]string{"--input", input, "--output", filepath.Join(root, "results"), "--cpuprofile", profile})
	if err == nil || !strings.Contains(err.Error(), "row 2") {
		t.Fatalf("run() error = %v, want input processing error", err)
	}
	assertGzipProfile(t, profile)
}

func TestRunRejectsProfilePathThatWouldOverwriteInput(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.csv")
	content := "campaign_id,date,impressions,clicks,spend,conversions\n"
	writeTestFile(t, input, content)
	err := run([]string{"--input", input, "--output", filepath.Join(root, "results"), "--cpuprofile", input})
	if err == nil || !strings.Contains(err.Error(), "conflicts with input/output") {
		t.Fatalf("run() error = %v, want profile-path conflict", err)
	}
	got, readErr := os.ReadFile(input)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(got) != content {
		t.Fatal("input file changed after profile-path conflict")
	}
}

func TestRunReportsCPUProfileCreateError(t *testing.T) {
	root := t.TempDir()
	input := filepath.Join(root, "input.csv")
	writeTestFile(t, input, "campaign_id,date,impressions,clicks,spend,conversions\n")
	profileDir := filepath.Join(root, "profile-dir")
	if err := os.Mkdir(profileDir, 0o700); err != nil {
		t.Fatal(err)
	}
	err := run([]string{"--input", input, "--output", filepath.Join(root, "results"), "--cpuprofile", profileDir})
	if err == nil || !strings.Contains(err.Error(), "create CPU profile") {
		t.Fatalf("run() error = %v, want profile create error", err)
	}
	if _, err := os.Stat(filepath.Join(root, "results")); !os.IsNotExist(err) {
		t.Errorf("application ran after profile creation failed (stat error %v)", err)
	}
}

func assertGzipProfile(t *testing.T, path string) {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open profile: %v", err)
	}
	defer file.Close()
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("open gzip CPU profile: %v", err)
	}
	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatalf("read CPU profile: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close gzip CPU profile: %v", err)
	}
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
