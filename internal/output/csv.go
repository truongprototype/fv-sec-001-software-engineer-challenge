// Package output writes ranked campaign metrics as CSV files.
package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"adaggregator/internal/aggregate"
)

var header = []string{
	"campaign_id",
	"total_impressions",
	"total_clicks",
	"total_spend",
	"total_conversions",
	"CTR",
	"CPA",
}

// WriteRankings writes the two challenge result files into outputDir. Existing
// files with these names are replaced.
func WriteRankings(outputDir string, topCTR, topCPA []aggregate.Metrics) error {
	if err := writeFile(filepath.Join(outputDir, "top10_ctr.csv"), topCTR); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(outputDir, "top10_cpa.csv"), topCPA); err != nil {
		return err
	}
	return nil
}

func writeFile(path string, metrics []aggregate.Metrics) (returnErr error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create output %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); returnErr == nil && err != nil {
			returnErr = fmt.Errorf("close output %q: %w", path, err)
		}
	}()

	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("write output %q header: %w", path, err)
	}
	for _, item := range metrics {
		campaign := item.Campaign
		cpa := ""
		if item.HasCPA {
			cpa = strconv.FormatFloat(item.CPA, 'f', 2, 64)
		}
		ctr := ""
		if item.HasCTR {
			ctr = strconv.FormatFloat(item.CTR, 'f', 4, 64)
		}
		record := []string{
			campaign.ID,
			strconv.FormatInt(campaign.Impressions, 10),
			strconv.FormatInt(campaign.Clicks, 10),
			strconv.FormatFloat(campaign.Spend, 'f', 2, 64),
			strconv.FormatInt(campaign.Conversions, 10),
			ctr,
			cpa,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("write output %q: %w", path, err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flush output %q: %w", path, err)
	}
	return nil
}
