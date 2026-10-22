package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"adaggregator/internal/app"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "aggregator:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("aggregator", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	inputPath := flags.String("input", "", "path to input CSV file")
	outputDir := flags.String("output", "", "directory for generated CSV files")
	cpuProfilePath := flags.String("cpuprofile", "", "write a CPU profile to this file")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments: %v", flags.Args())
	}
	if *inputPath == "" {
		return fmt.Errorf("--input is required")
	}
	if *outputDir == "" {
		return fmt.Errorf("--output is required")
	}
	if *cpuProfilePath != "" {
		if err := validateCPUProfilePath(*cpuProfilePath, *inputPath, *outputDir); err != nil {
			return err
		}
	}

	return withCPUProfile(*cpuProfilePath, func() error {
		return app.Run(*inputPath, *outputDir)
	})
}

func validateCPUProfilePath(profilePath, inputPath, outputDir string) error {
	protectedPaths := []string{
		inputPath,
		filepath.Join(outputDir, "top10_ctr.csv"),
		filepath.Join(outputDir, "top10_cpa.csv"),
	}
	for _, protectedPath := range protectedPaths {
		same, err := samePath(profilePath, protectedPath)
		if err != nil {
			return fmt.Errorf("compare CPU profile path with %q: %w", protectedPath, err)
		}
		if same {
			return fmt.Errorf("CPU profile path %q conflicts with input/output path %q", profilePath, protectedPath)
		}
	}
	return nil
}

func samePath(first, second string) (bool, error) {
	firstAbs, err := filepath.Abs(first)
	if err != nil {
		return false, err
	}
	secondAbs, err := filepath.Abs(second)
	if err != nil {
		return false, err
	}
	if filepath.Clean(firstAbs) == filepath.Clean(secondAbs) {
		return true, nil
	}
	firstInfo, firstErr := os.Stat(first)
	secondInfo, secondErr := os.Stat(second)
	if firstErr == nil && secondErr == nil {
		return os.SameFile(firstInfo, secondInfo), nil
	}
	return false, nil
}
