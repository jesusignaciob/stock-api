package service

import (
	"strconv"
	"strings"

	"stock-api/infrastructure/core/domain"
)

type ClassificationService struct{}

// NewClassificationService creates a new instance of ClassificationService.
// This service is responsible for classifying stocks based on various financial criteria.
func NewClassificationService() *ClassificationService {
	return &ClassificationService{}
}

// Classify classifies the stock based on various financial criteria.
// The classification process evaluates the stock's sector, target price changes, analyst actions, and ratings.
// It assigns one or more classifications to the stock, which are stored in the Classifications field.
func (s *ClassificationService) Classify(stock *domain.Stock) {
	// Initialize the classifications field as an empty slice
	stock.Classifications = []string{}
	classifications := make(map[string]struct{}) // Use a map to avoid duplicate classifications

	// 1. Classify by Sector (based on company name)
	// The sector classification is inferred from keywords in the company name.
	switch {
	case strings.Contains(stock.Company, "Medical"), strings.Contains(stock.Company, "Therapeutics"), strings.Contains(stock.Company, "Biopharma"), strings.Contains(stock.Company, "Pharma"):
		// Biotech sector includes companies in pharmaceuticals, biotechnology, and medical research.
		classifications["Biotech"] = struct{}{}
	case strings.Contains(stock.Company, "Tech"), strings.Contains(stock.Company, "Software"), strings.Contains(stock.Company, "Group"), strings.Contains(stock.Company, "Systems"), strings.Contains(stock.Company, "Solutions"):
		// Tech sector includes companies in software, hardware, and technology services.
		classifications["Tech"] = struct{}{}
	case strings.Contains(stock.Company, "Financial"), strings.Contains(stock.Company, "Bank"), strings.Contains(stock.Company, "Banc"), strings.Contains(stock.Company, "Capital"), strings.Contains(stock.Company, "Insurance"), strings.Contains(stock.Company, "Investments"), strings.Contains(stock.Company, "Advisors"):
		// Financial sector includes banks, insurance companies, and investment firms.
		classifications["Financial"] = struct{}{}
	case strings.Contains(stock.Company, "Energy"), strings.Contains(stock.Company, "Resources"), strings.Contains(stock.Company, "Petroleum"), strings.Contains(stock.Company, "Gas"):
		// Energy sector includes companies in oil, gas, and renewable energy.
		classifications["Energy"] = struct{}{}
	default:
		// If no specific sector is identified, classify as "Other Sector".
		classifications["Other Sector"] = struct{}{}
	}

	// 2. Classify by Target Price Change
	// This classification evaluates the percentage change between the initial and final target prices.
	priceFrom, errFrom := parsePrice(stock.TargetFrom)
	priceTo, errTo := parsePrice(stock.TargetTo)
	if errFrom == nil && errTo == nil && priceFrom > 0 {
		changePct := ((priceTo - priceFrom) / priceFrom) * 100
		switch {
		case changePct < -20:
			// A significant drop in target price indicates high risk and speculative behavior.
			classifications["High-Risk Speculative"] = struct{}{}
		case changePct > 10:
			// A significant increase in target price suggests potential growth opportunities.
			classifications["Potential Growth"] = struct{}{}
		}
	}

	// 3. Classify by Analyst Action
	// This classification is based on the actions taken by financial analysts.
	actionLower := strings.ToLower(stock.Action)
	switch {
	case strings.Contains(actionLower, "upgraded"):
		// "Upgraded" indicates a bullish signal from analysts.
		classifications["Bullish Signal"] = struct{}{}
	case strings.Contains(actionLower, "downgraded"):
		// "Downgraded" indicates a bearish signal from analysts.
		classifications["Bearish Signal"] = struct{}{}
	case strings.Contains(actionLower, "initiated"):
		// "Initiated" indicates new coverage or interest from analysts.
		classifications["New Coverage"] = struct{}{}
	}

	// 4. Classify by Rating
	// This classification evaluates the stock's rating provided by analysts.
	switch stock.RatingTo {
	case "Buy", "Outperform", "Strong-Buy":
		// Positive ratings indicate strong performance expectations.
		classifications["Analyst Positive"] = struct{}{}
	case "Sell", "Underweight":
		// Negative ratings indicate weak performance expectations.
		classifications["Analyst Negative"] = struct{}{}
	}

	// 5. Default classification if no other classifications exist
	// If no classifications are assigned, default to "Neutral".
	if len(classifications) == 0 {
		classifications["Neutral"] = struct{}{}
	}

	// Convert map keys to a sorted slice and assign to the stock's Classifications field
	for key := range classifications {
		stock.Classifications = append(stock.Classifications, key)
	}
}

// ClassifyBatch applies classification to each stock in the batch.
// This method iterates over a batch of stocks and applies the Classify method to each one.
func (s *ClassificationService) ClassifyBatch(batch []*domain.Stock) {
	for _, stock := range batch {
		s.Classify(stock)
	}
}

// parsePrice converts a price string (e.g., "$13.00") to a float64.
// It removes the "$" symbol and parses the remaining string as a float.
func parsePrice(priceStr string) (float64, error) {
	priceStr = strings.ReplaceAll(priceStr, "$", "")
	return strconv.ParseFloat(priceStr, 64)
}
