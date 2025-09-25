package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Date is a custom type for handling date-only values
type Date struct {
	time.Time
}

// UnmarshalJSON handles parsing of date strings in multiple formats
func (d *Date) UnmarshalJSON(data []byte) error {
	dateStr := strings.Trim(string(data), `"`)
	
	// Handle empty or null values
	if dateStr == "" || dateStr == "null" {
		return nil
	}
	
	// Try different date formats
	formats := []string{
		"2006-01-02",           // YYYY-MM-DD
		"2006-01-02T15:04:05Z", // RFC3339
		"2006-01-02T15:04:05Z07:00", // RFC3339 with timezone
	}
	
	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			d.Time = parsedTime
			return nil
		}
	}
	
	return fmt.Errorf("invalid date format: %s", dateStr)
}

// MarshalJSON formats the date as YYYY-MM-DD
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.Time.Format("2006-01-02"))
}

// String returns the date as YYYY-MM-DD string
func (d Date) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format("2006-01-02")
}