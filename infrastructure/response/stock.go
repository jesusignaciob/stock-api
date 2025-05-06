package response

import (
	"fmt"
	"regexp"
	"time"
)

type Stock struct {
	Ticker     string    `json:"ticker" validate:"required,uppercase,min=2,max=10"`
	TargetFrom string    `json:"target_from" validate:"omitempty,startswith=$"`
	TargetTo   string    `json:"target_to" validate:"omitempty,startswith=$"`
	Company    string    `json:"company" validate:"required"`
	Action     string    `json:"action" validate:"omitempty,oneof=upgraded downgraded maintained initiated"`
	Brokerage  string    `json:"brokerage" validate:"required"`
	RatingFrom string    `json:"rating_from" validate:"omitempty,oneof='Strong Buy' Buy Neutral Sell 'Strong Sell'"`
	RatingTo   string    `json:"rating_to" validate:"omitempty,oneof='Strong Buy' Buy Neutral Sell 'Strong Sell'"`
	Time       time.Time `json:"time" validate:"required"`
}

// Validate realiza validaciones personalizadas
func (s *Stock) Validate() error {
	// Validar formato de ticker (solo letras mayúsculas y números)
	matched, _ := regexp.MatchString(`^[A-Z0-9]+$`, s.Ticker)
	if !matched {
		return fmt.Errorf("ticker must contain only uppercase letters and numbers")
	}

	// Validar que el tiempo no sea futuro
	if s.Time.After(time.Now()) {
		return fmt.Errorf("time cannot be in the future")
	}

	return nil
}
