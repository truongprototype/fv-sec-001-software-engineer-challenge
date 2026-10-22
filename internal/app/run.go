package app

import (
	"fmt"
	"os"

	"adaggregator/internal/aggregate"
	"adaggregator/internal/csvdata"
)

// Run validates the input path and prepares the output directory. The data
// processing workflow is intentionally a scaffold for the next implementation step.
func Run(inputPath, outputDir string) error {
	info, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("inspect input %q: %w", inputPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("input %q is a directory, expected a CSV file", inputPath)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	_ = aggregate.New()
	_ = csvdata.Record{}
	return nil
}
