// Command count streams an input CSV and reports its data-row and campaign counts.
package main

import (
	"flag"
	"fmt"
	"os"

	"adaggregator/internal/csvdata"
)

func main() {
	input := flag.String("input", "", "CSV file to count")
	flag.Parse()
	if *input == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: go run ./benchmarks/count.go --input FILE")
		os.Exit(2)
	}
	file, err := os.Open(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	campaigns := make(map[string]struct{})
	var rows uint64
	err = csvdata.ReadEach(file, func(record csvdata.Record) error {
		rows++
		campaigns[record.CampaignID] = struct{}{}
		return nil
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	info, err := os.Stat(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("input_bytes=%d\ndata_rows=%d\ndistinct_campaigns=%d\n", info.Size(), rows, len(campaigns))
}
