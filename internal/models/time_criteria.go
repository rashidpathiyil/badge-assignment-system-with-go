package models

// TimePeriodCriteria represents criteria for counting unique time periods
type TimePeriodCriteria struct {
	PeriodType      string                 `json:"periodType"` // "day", "week", "month"
	PeriodCount     map[string]interface{} `json:"periodCount,omitempty"`
	ExcludeWeekends bool                   `json:"excludeWeekends,omitempty"`
	ExcludeHolidays bool                   `json:"excludeHolidays,omitempty"`
	Holidays        []string               `json:"holidays,omitempty"`
}

// PatternCriteria represents criteria for detecting patterns in event frequency
type PatternCriteria struct {
	Pattern        string  `json:"pattern"`                  // "consistent", "increasing", "decreasing"
	PeriodType     string  `json:"periodType"`               // "day", "week", "month"
	MinPeriods     int     `json:"minPeriods"`               // Minimum number of periods to analyze
	MinIncreasePct float64 `json:"minIncreasePct,omitempty"` // For increasing pattern
	MaxDecreasePct float64 `json:"maxDecreasePct,omitempty"` // For decreasing pattern
	MaxDeviation   float64 `json:"maxDeviation,omitempty"`   // For consistent pattern
}

// SequenceCriteria represents criteria for verifying event sequences
type SequenceCriteria struct {
	Sequence      []string `json:"sequence"`                // Ordered list of event types
	MaxGapSeconds int      `json:"maxGapSeconds,omitempty"` // Maximum gap between events
	RequireStrict bool     `json:"requireStrict,omitempty"` // If true, no other events can be between sequence events
}

// GapCriteria represents criteria for detecting gaps between events
type GapCriteria struct {
	MaxGapHours       float64                `json:"maxGapHours"`
	MinGapHours       float64                `json:"minGapHours,omitempty"`
	PeriodType        string                 `json:"periodType,omitempty"` // "all", "business-days"
	ExcludeConditions map[string]interface{} `json:"excludeConditions,omitempty"`
}

// DurationCriteria represents criteria for time duration calculations
type DurationCriteria struct {
	StartEvent map[string]interface{} `json:"startEvent"`
	EndEvent   map[string]interface{} `json:"endEvent"`
	Duration   map[string]interface{} `json:"duration,omitempty"`
	Unit       string                 `json:"unit,omitempty"` // "hours", "days", "minutes", defaults to "hours"
}

// AggregationCriteria represents criteria for advanced aggregation functions
type AggregationCriteria struct {
	Type       string                 `json:"type"`  // "min", "max", "avg"
	Field      string                 `json:"field"` // Field path in event payload
	Value      map[string]interface{} `json:"value,omitempty"`
	TimeWindow map[string]interface{} `json:"timeWindow,omitempty"` // Optional time window
}

// TimeWindowCriteria defines a time window for filtering events
type TimeWindowCriteria struct {
	Start            string `json:"start,omitempty"`
	End              string `json:"end,omitempty"`
	Last             string `json:"last,omitempty"` // e.g., "30d", "2w", "1m"
	BusinessDaysOnly bool   `json:"businessDaysOnly,omitempty"`
}
