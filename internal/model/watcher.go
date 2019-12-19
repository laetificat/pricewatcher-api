package model

import "time"

// Watcher contains metadata and a list of prices.
type Watcher struct {
	ID           int
	Name         string
	URL          string
	Domain       string
	LastChecked  time.Time
	IsChecking   bool
	PriceHistory []Price
}
