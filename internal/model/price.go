package model

import "time"

// Price is a single price object that has links a value with a timestamp.
type Price struct {
	Value     float32
	Timestamp time.Time
}
