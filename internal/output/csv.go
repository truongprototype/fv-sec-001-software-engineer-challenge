// Package output writes ranked campaign metrics as CSV files.
package output

import (
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
			cpa = formatDecimal(item.CPA, 2, 18)
		}
		ctr := ""
		if item.HasCTR {
			ctr = strconv.FormatFloat(item.CTR, 'f', 4, 64)
		}
		record := []string{
			campaign.ID,
			strconv.FormatInt(campaign.Impressions, 10),
			strconv.FormatInt(campaign.Clicks, 10),
			formatDecimal(campaign.Spend, 2, 0),
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

// formatDecimal preserves terminating decimals up to maxPlaces fractional
// digits. Values requiring more places are rounded; trailing zeroes are
// removed without going below minPlaces. A non-positive maxPlaces is unbounded.
func formatDecimal(value *big.Rat, minPlaces, maxPlaces int) string {
	if value == nil || value.Sign() == 0 {
		return "0"
	}
	denominator := new(big.Int).Set(value.Denom())
	two, five := big.NewInt(2), big.NewInt(5)
	twoCount, fiveCount := 0, 0
	for new(big.Int).Mod(denominator, two).Sign() == 0 {
		denominator.Div(denominator, two)
		twoCount++
	}
	for new(big.Int).Mod(denominator, five).Sign() == 0 {
		denominator.Div(denominator, five)
		fiveCount++
	}
	places := twoCount
	if fiveCount > places {
		places = fiveCount
	}
	if places < minPlaces {
		places = minPlaces
	}
	if maxPlaces > 0 && (denominator.Cmp(big.NewInt(1)) != 0 || places > maxPlaces) {
		places = maxPlaces
	}
	formatted := value.FloatString(places)
	if decimalPoint := strings.IndexByte(formatted, '.'); decimalPoint >= 0 {
		fractionalPlaces := len(formatted) - decimalPoint - 1
		for fractionalPlaces > minPlaces && strings.HasSuffix(formatted, "0") {
			formatted = strings.TrimSuffix(formatted, "0")
			fractionalPlaces--
		}
		if fractionalPlaces == 0 {
			formatted = strings.TrimSuffix(formatted, ".")
		}
	}
	return formatted
}
