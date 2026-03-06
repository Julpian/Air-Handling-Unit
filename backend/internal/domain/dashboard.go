package domain

import "time"

type FilterPressureChartRow struct {
	Month time.Time
	Label string
	Value float64
}