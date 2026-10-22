// Package csvdata defines the CSV input and output boundary.
package csvdata

// Record represents one parsed input row.
type Record struct {
	CampaignID  string
	Date        string
	Impressions int64
	Clicks      int64
	Spend       float64
	Conversions int64
}
