package csvdata

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var expectedHeader = []string{"campaign_id", "date", "impressions", "clicks", "spend", "conversions"}
var decimalNumber = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)

// ReadEach streams a header-bearing CSV and calls visit for each valid record.
// It stops at the first malformed row and includes its logical row number in
// the returned error. Blank lines are ignored by encoding/csv.
func ReadEach(input io.Reader, visit func(Record) error) error {
	if visit == nil {
		return fmt.Errorf("record visitor must not be nil")
	}

	reader := csv.NewReader(input)
	reader.FieldsPerRecord = len(expectedHeader)
	reader.ReuseRecord = true

	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("read CSV header: input is empty")
		}
		return fmt.Errorf("read CSV header: %w", err)
	}
	if len(header) > 0 {
		header[0] = strings.TrimPrefix(header[0], "\uFEFF")
	}
	for i, name := range expectedHeader {
		if header[i] != name {
			return fmt.Errorf("invalid CSV header column %d: got %q, want %q", i+1, header[i], name)
		}
	}

	rowNumber := 1
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		rowNumber++
		if err != nil {
			return fmt.Errorf("read CSV row %d: %w", rowNumber, err)
		}

		record, err := parseRecord(fields)
		if err != nil {
			return fmt.Errorf("parse CSV row %d: %w", rowNumber, err)
		}
		if err := visit(record); err != nil {
			return fmt.Errorf("process CSV row %d: %w", rowNumber, err)
		}
	}
}

func parseRecord(fields []string) (Record, error) {
	var record Record
	if len(fields) != len(expectedHeader) {
		return record, fmt.Errorf("got %d fields, want %d", len(fields), len(expectedHeader))
	}

	record.CampaignID = fields[0]
	if _, err := time.Parse("2006-01-02", fields[1]); err != nil {
		return record, fmt.Errorf("field date %q is invalid: %w", fields[1], err)
	}
	record.Date = fields[1]

	var err error
	if record.Impressions, err = parseInteger("impressions", fields[2]); err != nil {
		return record, err
	}
	if record.Clicks, err = parseInteger("clicks", fields[3]); err != nil {
		return record, err
	}
	if !decimalNumber.MatchString(fields[4]) {
		return record, fmt.Errorf("field spend %q is not a valid decimal number", fields[4])
	}
	var ok bool
	record.Spend, ok = new(big.Rat).SetString(fields[4])
	if !ok {
		return record, fmt.Errorf("field spend %q is not a valid decimal number", fields[4])
	}
	if record.Conversions, err = parseInteger("conversions", fields[5]); err != nil {
		return record, err
	}
	return record, nil
}

func parseInteger(field, value string) (int64, error) {
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("field %s %q is not a valid 64-bit integer: %w", field, value, err)
	}
	return parsed, nil
}
