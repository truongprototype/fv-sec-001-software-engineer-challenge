// Package aggregate contains campaign totals and ranking logic.
package aggregate

// Campaign holds totals accumulated from input rows for one campaign.
type Campaign struct {
	ID          string
	Impressions int64
	Clicks      int64
	Spend       float64
	Conversions int64
}

// New creates an empty campaign index keyed by campaign ID.
func New() map[string]*Campaign {
	return make(map[string]*Campaign)
}
