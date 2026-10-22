package app

import (
	"fmt"
	"os"

	"adaggregator/internal/aggregate"
	"adaggregator/internal/csvdata"
	"adaggregator/internal/output"
)

// Run streams the input CSV, aggregates by campaign, and writes both rankings.
func Run(inputPath, outputDir string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input %q: %w", inputPath, err)
	}
	defer input.Close()
	info, err := input.Stat()
	if err != nil {
		return fmt.Errorf("inspect input %q: %w", inputPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("input %q is a directory, expected a CSV file", inputPath)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	aggregator := aggregate.New()
	if err := csvdata.ReadEach(input, func(record csvdata.Record) error {
		return aggregator.Add(record)
	}); err != nil {
		return fmt.Errorf("aggregate input %q: %w", inputPath, err)
	}
	return output.WriteRankings(outputDir, aggregator.TopCTR(10), aggregator.TopCPA(10))
}
