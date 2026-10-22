// Package aggregate contains campaign totals and ranking logic.
package aggregate

import (
	"fmt"
	"math"
	"sort"

	"adaggregator/internal/csvdata"
)

// Campaign holds totals accumulated from input rows for one campaign.
type Campaign struct {
	ID          string
	Impressions int64
	Clicks      int64
	Spend       float64
	Conversions int64
}

// Metrics contains ratios derived from a campaign's aggregate totals. A false
// validity flag means the corresponding denominator was zero.
type Metrics struct {
	Campaign Campaign
	CTR      float64
	CPA      float64
	HasCTR   bool
	HasCPA   bool
}

// Aggregator holds one aggregate per campaign ID.
type Aggregator struct {
	campaigns map[string]*Campaign
}

// New creates an empty campaign aggregator.
func New() *Aggregator {
	return &Aggregator{campaigns: make(map[string]*Campaign)}
}

// Add adds one input record to its campaign totals.
func (a *Aggregator) Add(record csvdata.Record) error {
	current, exists := a.campaigns[record.CampaignID]
	if !exists {
		current = &Campaign{ID: record.CampaignID}
	}

	impressions, ok := addInt64(current.Impressions, record.Impressions)
	if !ok {
		return fmt.Errorf("campaign %q impressions total overflows int64", record.CampaignID)
	}
	clicks, ok := addInt64(current.Clicks, record.Clicks)
	if !ok {
		return fmt.Errorf("campaign %q clicks total overflows int64", record.CampaignID)
	}
	conversions, ok := addInt64(current.Conversions, record.Conversions)
	if !ok {
		return fmt.Errorf("campaign %q conversions total overflows int64", record.CampaignID)
	}
	spend := current.Spend + record.Spend
	if math.IsInf(spend, 0) || math.IsNaN(spend) {
		return fmt.Errorf("campaign %q spend total is not finite", record.CampaignID)
	}

	current.Impressions = impressions
	current.Clicks = clicks
	current.Conversions = conversions
	current.Spend = spend
	a.campaigns[record.CampaignID] = current
	return nil
}

// Campaigns returns a copy of the aggregate values.
func (a *Aggregator) Campaigns() []Campaign {
	result := make([]Campaign, 0, len(a.campaigns))
	for _, campaign := range a.campaigns {
		result = append(result, *campaign)
	}
	return result
}

// CalculateMetrics derives CTR and CPA from aggregate totals without dividing
// by zero. Campaign totals are copied into the result.
func CalculateMetrics(campaign Campaign) Metrics {
	metrics := Metrics{Campaign: campaign}
	if campaign.Impressions != 0 {
		metrics.CTR = float64(campaign.Clicks) / float64(campaign.Impressions)
		metrics.HasCTR = true
	}
	if campaign.Conversions != 0 {
		metrics.CPA = campaign.Spend / float64(campaign.Conversions)
		metrics.HasCPA = true
	}
	return metrics
}

// TopCTR returns up to limit campaigns ranked by descending CTR. Campaigns
// with zero impressions are omitted because their CTR is undefined.
func (a *Aggregator) TopCTR(limit int) []Metrics {
	if limit <= 0 {
		return nil
	}
	metrics := make([]Metrics, 0, len(a.campaigns))
	for _, campaign := range a.campaigns {
		item := CalculateMetrics(*campaign)
		if item.HasCTR {
			metrics = append(metrics, item)
		}
	}
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].CTR != metrics[j].CTR {
			return metrics[i].CTR > metrics[j].CTR
		}
		return metrics[i].Campaign.ID < metrics[j].Campaign.ID
	})
	return firstMetrics(metrics, limit)
}

// TopCPA returns up to limit campaigns ranked by ascending CPA. Campaigns
// with zero conversions are excluded because their CPA is undefined.
func (a *Aggregator) TopCPA(limit int) []Metrics {
	if limit <= 0 {
		return nil
	}
	metrics := make([]Metrics, 0, len(a.campaigns))
	for _, campaign := range a.campaigns {
		item := CalculateMetrics(*campaign)
		if item.HasCPA {
			metrics = append(metrics, item)
		}
	}
	sort.Slice(metrics, func(i, j int) bool {
		if metrics[i].CPA != metrics[j].CPA {
			return metrics[i].CPA < metrics[j].CPA
		}
		return metrics[i].Campaign.ID < metrics[j].Campaign.ID
	})
	return firstMetrics(metrics, limit)
}

func firstMetrics(metrics []Metrics, limit int) []Metrics {
	if limit > len(metrics) {
		limit = len(metrics)
	}
	return metrics[:limit]
}

func addInt64(a, b int64) (int64, bool) {
	if (b > 0 && a > math.MaxInt64-b) || (b < 0 && a < math.MinInt64-b) {
		return 0, false
	}
	return a + b, true
}
