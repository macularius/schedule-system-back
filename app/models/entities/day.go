package entities

import "time"

// Day type of schedule's day
type Day struct {
	Date       time.Time
	RangeStart string
	RangeEnd   string
}
