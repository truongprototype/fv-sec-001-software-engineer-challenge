package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunProcessesInputAndWritesRankings(t *testing.T) {
	root := t.TempDir()
	inputPath := filepath.Join(root, "input.csv")
	input := "campaign_id,date,impressions,clicks,spend,conversions\n" +
		"A,2025-01-01,100,20,40.00,2\n" +
		"B,2025-01-01,100,10,5.00,1\n" +
		"A,2025-01-02,100,30,20.00,3\n"
	if err := os.WriteFile(inputPath, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(root, "nested", "results")
	if err := Run(inputPath, outputDir); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	ctr, err := os.ReadFile(filepath.Join(outputDir, "top10_ctr.csv"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ctr), "A,200,50,60.00,5,0.2500,12.00") {
		t.Errorf("CTR output missing aggregated campaign A: %q", ctr)
	}
	cpa, err := os.ReadFile(filepath.Join(outputDir, "top10_cpa.csv"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(cpa), "B,100,10,5.00,1,0.1000,5.00") {
		t.Errorf("CPA output missing campaign B: %q", cpa)
	}
}

func TestRunMalformedInputDoesNotWriteResultFiles(t *testing.T) {
	root := t.TempDir()
	inputPath := filepath.Join(root, "bad.csv")
	if err := os.WriteFile(inputPath, []byte("campaign_id,date,impressions,clicks,spend,conversions\nA,not-a-date,1,1,1,1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	outputDir := filepath.Join(root, "results")
	if err := Run(inputPath, outputDir); err == nil || !strings.Contains(err.Error(), "row 2") {
		t.Fatalf("Run() error = %v, want malformed row context", err)
	}
	for _, name := range []string{"top10_ctr.csv", "top10_cpa.csv"} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); !os.IsNotExist(err) {
			t.Errorf("result %s exists after input failure (stat error %v)", name, err)
		}
	}
}

func TestRunMissingInput(t *testing.T) {
	err := Run(filepath.Join(t.TempDir(), "missing.csv"), t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "open input") {
		t.Fatalf("Run() error = %v, want input-open error", err)
	}
}

func TestRunRejectsDirectoryAsInput(t *testing.T) {
	dir := t.TempDir()
	err := Run(dir, filepath.Join(t.TempDir(), "results"))
	if err == nil || !strings.Contains(err.Error(), "expected a CSV file") {
		t.Fatalf("Run(directory) error = %v", err)
	}
}

func TestRunRejectsOutputPathThatIsAFile(t *testing.T) {
	root := t.TempDir()
	inputPath := filepath.Join(root, "input.csv")
	if err := os.WriteFile(inputPath, []byte("campaign_id,date,impressions,clicks,spend,conversions\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(root, "output")
	if err := os.WriteFile(outputPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Run(inputPath, outputPath); err == nil || !strings.Contains(err.Error(), "create output directory") {
		t.Fatalf("Run() error = %v, want output directory error", err)
	}
}
