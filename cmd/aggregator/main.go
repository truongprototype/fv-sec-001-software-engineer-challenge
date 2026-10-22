package main

import (
	"flag"
	"fmt"
	"os"

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

	return app.Run(*inputPath, *outputDir)
}
