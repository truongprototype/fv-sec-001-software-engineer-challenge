package output

import (
	"os"
	"path/filepath"
	"testing"

	"adaggregator/internal/aggregate"
)

func TestWriteRankingsWritesHeadersEscapingAndUndefinedMetric(t *testing.T) {
	dir := t.TempDir()
	topCTR := []aggregate.Metrics{{
		Campaign: aggregate.Campaign{ID: "CMP,001", Impressions: 100, Clicks: 25, Spend: 12.345, Conversions: 0},
		CTR:      0.25,
		HasCTR:   true,
	}}
	if err := WriteRankings(dir, topCTR, nil); err != nil {
		t.Fatalf("WriteRankings() error = %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dir, "top10_ctr.csv"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	want := "campaign_id,total_impressions,total_clicks,total_spend,total_conversions,CTR,CPA\n\"CMP,001\",100,25,12.35,0,0.2500,\n"
	if string(got) != want {
		t.Errorf("top10_ctr.csv = %q, want %q", got, want)
	}

	gotCPA, err := os.ReadFile(filepath.Join(dir, "top10_cpa.csv"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(gotCPA) != "campaign_id,total_impressions,total_clicks,total_spend,total_conversions,CTR,CPA\n" {
		t.Errorf("empty top10_cpa.csv = %q", gotCPA)
	}
}

func TestWriteRankingsReplacesExistingFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "top10_ctr.csv")
	if err := os.WriteFile(path, []byte("old contents"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteRankings(dir, nil, nil); err != nil {
		t.Fatalf("WriteRankings() error = %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) == "old contents" {
		t.Fatal("existing file was not replaced")
	}
}

func TestWriteRankingsReturnsCreateError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "output")
	if err := WriteRankings(path, nil, nil); err == nil {
		t.Fatal("WriteRankings() succeeded for missing output directory")
	}
}
