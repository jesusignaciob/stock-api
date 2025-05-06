package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Stock represents the stock entity in the system.
// It contains information about the stock's ticker, company, classifications, and other attributes.
type Stock struct {
	gorm.Model
	Ticker          string      `gorm:"size:10;not null;index"`             // Stock ticker (e.g., "AAPL")
	TargetFrom      string      `gorm:"size:20" json:"target_from"`         // Initial target price
	TargetTo        string      `gorm:"size:20" json:"target_to"`           // Final target price
	Company         string      `gorm:"size:255;not null"`                  // Company name
	Action          string      `gorm:"size:100"`                           // Analyst action (e.g., "upgraded by")
	Brokerage       string      `gorm:"size:255;not null"`                  // Brokerage firm
	RatingFrom      string      `gorm:"size:50" json:"rating_from"`         // Initial rating
	RatingTo        string      `gorm:"size:50" json:"rating_to"`           // Final rating
	Time            time.Time   `gorm:"not null;index"`                     // Timestamp of the stock event
	Classifications StringArray `gorm:"type:text[]" json:"classifications"` // Classifications for the stock
}

// StringArray wraps pq.StringArray to provide better JSON handling and database integration.
type StringArray pq.StringArray

// Scan implements the Scanner interface for database deserialization.
// It converts the database value into a StringArray.
func (sa *StringArray) Scan(value interface{}) error {
	return (*pq.StringArray)(sa).Scan(value)
}

// Value implements the driver Valuer interface for database serialization.
// It converts the StringArray into a database-compatible format.
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return pq.StringArray{"Neutral"}.Value() // Default value if empty
	}
	return pq.StringArray(sa).Value()
}

// MarshalJSON provides custom JSON marshaling for StringArray.
// If the array is nil, it defaults to ["Neutral"].
func (sa StringArray) MarshalJSON() ([]byte, error) {
	if sa == nil {
		return json.Marshal([]string{"Neutral"})
	}
	return json.Marshal([]string(sa))
}

// UnmarshalJSON provides custom JSON unmarshaling for StringArray.
// It converts a JSON array into a StringArray.
func (sa *StringArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*sa = StringArray(arr)
	return nil
}

// BeforeCreate is a GORM hook that ensures the classifications field is never empty.
// If the classifications field is empty, it defaults to ["Neutral"] before creating the record.
func (s *Stock) BeforeCreate(_ *gorm.DB) error {
	if len(s.Classifications) == 0 {
		s.Classifications = []string{"Neutral"}
	}
	return nil
}

// Validate performs custom validations for the Stock model.
// It ensures the ticker format is valid and the time is not in the future.
func (s *Stock) Validate() error {
	// Validate ticker format (only uppercase letters and numbers)
	matched, _ := regexp.MatchString(`^[A-Z0-9]+$`, s.Ticker)
	if !matched {
		return fmt.Errorf("ticker must contain only uppercase letters and numbers")
	}

	// Validate that the time is not in the future
	if s.Time.After(time.Now()) {
		return fmt.Errorf("time cannot be in the future")
	}

	return nil
}
