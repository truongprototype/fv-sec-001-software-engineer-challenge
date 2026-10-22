package aggregate

import (
	"math"
	"math/big"
	"testing"

	"adaggregator/internal/csvdata"
)

func TestAggregatorCombinesCampaignRows(t *testing.T) {
	aggregator := New()
	rows := []csvdata.Record{
		{CampaignID: "A", Impressions: 10, Clicks: 2, Spend: rat("1.25"), Conversions: 1},
		{CampaignID: "B", Impressions: 5, Clicks: 1, Spend: rat("0.5"), Conversions: 0},
		{CampaignID: "A", Impressions: 20, Clicks: 4, Spend: rat("2.75"), Conversions: 3},
	}
	for _, row := range rows {
		if err := aggregator.Add(row); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
	}

	campaigns := make(map[string]Campaign)
	for _, campaign := range aggregator.Campaigns() {
		campaigns[campaign.ID] = campaign
	}
	if len(campaigns) != 2 {
		t.Fatalf("got %d campaigns, want 2", len(campaigns))
	}
	if got := campaigns["A"]; got.Impressions != 30 || got.Clicks != 6 || got.Spend.Cmp(rat("4")) != 0 || got.Conversions != 4 {
		t.Errorf("campaign A totals = %#v", got)
	}
	if got := campaigns["B"]; got.Impressions != 5 || got.Clicks != 1 || got.Spend.Cmp(rat("0.5")) != 0 || got.Conversions != 0 {
		t.Errorf("campaign B totals = %#v", got)
	}
}

func TestAggregatorOverflowLeavesTotalsUnchanged(t *testing.T) {
	aggregator := New()
	if err := aggregator.Add(csvdata.Record{CampaignID: "A", Impressions: math.MaxInt64, Spend: rat("1")}); err != nil {
		t.Fatalf("initial Add() error = %v", err)
	}
	if err := aggregator.Add(csvdata.Record{CampaignID: "A", Impressions: 1, Spend: rat("2")}); err == nil {
		t.Fatal("Add() succeeded; want overflow error")
	}
	campaign := aggregator.Campaigns()[0]
	if campaign.Impressions != math.MaxInt64 || campaign.Spend.Cmp(rat("1")) != 0 {
		t.Errorf("campaign after failed Add() = %#v; want original totals", campaign)
	}
}

func TestAggregatorAddsHighPrecisionSpendExactly(t *testing.T) {
	aggregator := New()
	if err := aggregator.Add(csvdata.Record{CampaignID: "A", Spend: rat("0.1")}); err != nil {
		t.Fatalf("initial Add() error = %v", err)
	}
	if err := aggregator.Add(csvdata.Record{CampaignID: "A", Spend: rat("0.2")}); err != nil {
		t.Fatalf("second Add() error = %v", err)
	}
	if got := aggregator.Campaigns()[0].Spend; got.Cmp(rat("0.3")) != 0 {
		t.Errorf("spend total = %s, want exact 0.3", got.RatString())
	}
}

func TestCalculateMetricsUsesAggregateTotalsAndHandlesZeroDenominators(t *testing.T) {
	tests := []struct {
		name     string
		campaign Campaign
		ctr, cpa float64
		hasCTR   bool
		hasCPA   bool
	}{
		{
			name:     "both defined",
			campaign: Campaign{ID: "A", Impressions: 100, Clicks: 25, Spend: rat("30"), Conversions: 3},
			ctr:      0.25,
			cpa:      10,
			hasCTR:   true,
			hasCPA:   true,
		},
		{
			name:     "zero impressions",
			campaign: Campaign{ID: "B", Clicks: 5, Spend: rat("2"), Conversions: 1},
			cpa:      2,
			hasCPA:   true,
		},
		{
			name:     "zero conversions",
			campaign: Campaign{ID: "C", Impressions: 10, Clicks: 2, Spend: rat("8")},
			ctr:      0.2,
			hasCTR:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CalculateMetrics(test.campaign)
			wantCPA := new(big.Rat).SetFloat64(test.cpa)
			if got.CTR != test.ctr || got.HasCPA != test.hasCPA || got.HasCTR != test.hasCTR || (got.HasCPA && got.CPA.Cmp(wantCPA) != 0) {
				t.Errorf("CalculateMetrics() = %#v", got)
			}
		})
	}
}

func TestTopCTRSortsDescendingWithDeterministicTies(t *testing.T) {
	aggregator := aggregatorWithCampaigns(t, []Campaign{
		{ID: "C", Impressions: 100, Clicks: 20},
		{ID: "A", Impressions: 100, Clicks: 20},
		{ID: "B", Impressions: 100, Clicks: 30},
		{ID: "zero", Impressions: 0, Clicks: 100},
	})

	got := aggregator.TopCTR(10)
	assertCampaignIDs(t, got, []string{"B", "A", "C"})
	if got[0].CTR != 0.3 || got[1].CTR != 0.2 {
		t.Errorf("CTR values = %v, %v", got[0].CTR, got[1].CTR)
	}
}

func TestTopCPASortsAscendingExcludesZeroConversionsAndTiesByID(t *testing.T) {
	aggregator := aggregatorWithCampaigns(t, []Campaign{
		{ID: "C", Spend: rat("20"), Conversions: 2},
		{ID: "A", Spend: rat("10"), Conversions: 1},
		{ID: "B", Spend: rat("5"), Conversions: 1},
		{ID: "zero", Spend: rat("0"), Conversions: 0},
	})

	got := aggregator.TopCPA(10)
	assertCampaignIDs(t, got, []string{"B", "A", "C"})
	if got[0].CPA.Cmp(rat("5")) != 0 || got[1].CPA.Cmp(rat("10")) != 0 || got[2].CPA.Cmp(rat("10")) != 0 {
		t.Errorf("CPA values = %v, %v, %v", got[0].CPA, got[1].CPA, got[2].CPA)
	}
}

func TestTopCPATreatsEquivalentDecimalSumsAsExactTies(t *testing.T) {
	aggregator := New()
	for _, row := range []csvdata.Record{
		{CampaignID: "A", Spend: rat("0.1"), Conversions: 1},
		{CampaignID: "A", Spend: rat("0.2")},
		{CampaignID: "B", Spend: rat("0.3"), Conversions: 1},
	} {
		if err := aggregator.Add(row); err != nil {
			t.Fatalf("Add() error = %v", err)
		}
	}
	assertCampaignIDs(t, aggregator.TopCPA(10), []string{"A", "B"})
}

func TestTopRankingsRespectLimitAndEmptyCandidates(t *testing.T) {
	aggregator := aggregatorWithCampaigns(t, []Campaign{
		{ID: "A", Impressions: 10, Clicks: 1},
		{ID: "B", Impressions: 10, Clicks: 2},
	})
	if got := aggregator.TopCTR(1); len(got) != 1 || got[0].Campaign.ID != "B" {
		t.Errorf("TopCTR(1) = %#v", got)
	}
	if got := aggregator.TopCTR(0); len(got) != 0 {
		t.Errorf("TopCTR(0) = %#v, want empty", got)
	}
	if got := aggregator.TopCPA(10); len(got) != 0 {
		t.Errorf("TopCPA without conversions = %#v, want empty", got)
	}
	if got := New().TopCTR(10); len(got) != 0 {
		t.Errorf("empty TopCTR = %#v, want empty", got)
	}
}

func aggregatorWithCampaigns(t *testing.T, campaigns []Campaign) *Aggregator {
	t.Helper()
	aggregator := New()
	for _, campaign := range campaigns {
		if err := aggregator.Add(csvdata.Record{
			CampaignID:  campaign.ID,
			Impressions: campaign.Impressions,
			Clicks:      campaign.Clicks,
			Spend:       campaign.Spend,
			Conversions: campaign.Conversions,
		}); err != nil {
			t.Fatalf("Add(%q) error = %v", campaign.ID, err)
		}
	}
	return aggregator
}

func rat(value string) *big.Rat {
	parsed, ok := new(big.Rat).SetString(value)
	if !ok {
		panic("invalid test rational: " + value)
	}
	return parsed
}

func assertCampaignIDs(t *testing.T, metrics []Metrics, want []string) {
	t.Helper()
	if len(metrics) != len(want) {
		t.Fatalf("got %d ranked campaigns, want %d", len(metrics), len(want))
	}
	for i, metric := range metrics {
		if metric.Campaign.ID != want[i] {
			t.Errorf("rank %d campaign = %q, want %q", i+1, metric.Campaign.ID, want[i])
		}
	}
}
