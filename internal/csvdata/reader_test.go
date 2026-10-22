package csvdata

import (
	"math/big"
	"strings"
	"testing"
)

func TestReadEachStreamsValidRecords(t *testing.T) {
	input := "\uFEFFcampaign_id,date,impressions,clicks,spend,conversions\r\n\"CMP,001\",2025-01-01,12000,300,45.50,12\r\nCMP002,2025-01-02,8000,120,28,4\r\n"
	var records []Record
	if err := ReadEach(strings.NewReader(input), func(record Record) error {
		records = append(records, record)
		return nil
	}); err != nil {
		t.Fatalf("ReadEach() error = %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("got %d records, want 2", len(records))
	}
	if records[0].CampaignID != "CMP,001" || records[0].Spend.Cmp(big.NewRat(455, 10)) != 0 || records[0].Conversions != 12 {
		t.Errorf("first record = %#v", records[0])
	}
	if records[1].CampaignID != "CMP002" {
		t.Errorf("second campaign ID = %q", records[1].CampaignID)
	}
}

func TestReadEachRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name, input, wantError string
	}{
		{"empty input", "", "input is empty"},
		{"wrong header", "date,campaign_id,impressions,clicks,spend,conversions\n", "invalid CSV header"},
		{"invalid date", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-02-30,1,1,1,1\n", "parse CSV row 2"},
		{"invalid integer", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1.5,1,1,1\n", "field impressions"},
		{"integer overflow", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,9223372036854775808,1,1,1\n", "field impressions"},
		{"invalid decimal", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,NaN,1\n", "valid decimal number"},
		{"high precision decimal", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,0.1234567890123456789,1\n", ""},
		{"wrong field count", "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,1\n", "read CSV row 2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ReadEach(strings.NewReader(test.input), func(Record) error { return nil })
			if test.wantError == "" {
				if err != nil {
					t.Fatalf("ReadEach() error = %v, want success", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantError) {
				t.Fatalf("ReadEach() error = %v, want containing %q", err, test.wantError)
			}
		})
	}
}

func TestReadEachPropagatesVisitorError(t *testing.T) {
	input := "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,1,1\n"
	err := ReadEach(strings.NewReader(input), func(Record) error { return errVisitor{} })
	if err == nil || !strings.Contains(err.Error(), "process CSV row 2") {
		t.Fatalf("ReadEach() error = %v, want visitor error context", err)
	}
}

func TestReadEachRequiresVisitor(t *testing.T) {
	if err := ReadEach(strings.NewReader(""), nil); err == nil || !strings.Contains(err.Error(), "visitor") {
		t.Fatalf("ReadEach(nil visitor) error = %v", err)
	}
}

func TestReadEachPreservesSpendBeyondCents(t *testing.T) {
	input := "campaign_id,date,impressions,clicks,spend,conversions\nA,2025-01-01,1,1,0.1234567890123456789,1\n"
	var got *big.Rat
	err := ReadEach(strings.NewReader(input), func(record Record) error {
		got = record.Spend
		return nil
	})
	if err != nil {
		t.Fatalf("ReadEach() error = %v", err)
	}
	want, _ := new(big.Rat).SetString("0.1234567890123456789")
	if got == nil || got.Cmp(want) != 0 {
		t.Fatalf("parsed spend = %v, want exact %s", got, want.RatString())
	}
}

type errVisitor struct{}

func (errVisitor) Error() string { return "visitor failed" }
