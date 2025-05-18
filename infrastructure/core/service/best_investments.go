package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"stock-api/infrastructure/core/domain"
)

type BestInvestmentsServiceImpl struct{}

func NewBestInvestmentsService() *BestInvestmentsServiceImpl {
	return &BestInvestmentsServiceImpl{}
}

// GetStockRecommendations generates a list of stock recommendations based on their scores.
// It sorts, and limits the provided stock data to produce the top recommendations.
//
// Parameters:
//   - stocks: A slice of Stock objects representing the available stock data.
//   - limit: An integer specifying the maximum number of recommendations to return.
//
// Returns:
//   - A slice of Recommendation objects containing the top stock recommendations,
//     sorted by their calculated scores in descending order.
//
// The function performs the following steps:
//  1. Filters the input stocks using the filterStocks function.
//  2. Sorts the filtered stocks in descending order based on their calculated scores.
//  3. Limits the number of recommendations to the specified limit or the total number of filtered stocks.
//  4. Constructs and returns a slice of Recommendation objects, including the position, ticker, company name,
//     score, and rationale for each recommended stock.
func (s *BestInvestmentsServiceImpl) GetStockRecommendations(stocks []domain.Stock, limit int) []domain.Recommendation {
	// Sort
	sort.Slice(stocks, func(i, j int) bool {
		return calculateScore(stocks[i]) > calculateScore(stocks[j])
	})

	// Limit results
	if limit > len(stocks) {
		limit = len(stocks)
	}

	// Prepare response
	recommendations := make([]domain.Recommendation, limit)
	for i := 0; i < limit; i++ {
		stock := stocks[i]
		recommendations[i] = domain.Recommendation{
			Position:  i + 1,
			Ticker:    stock.Ticker,
			Company:   stock.Company,
			Score:     calculateScore(stock),
			Rationale: getRationale(stock),
		}
	}

	return recommendations
}

// calculateScore calculates the score of a stock based on various factors.
// The score is determined by growth potential, positive classifications, and analyst ratings.
func calculateScore(stock domain.Stock) float64 {
	score := 0.0

	// 1. Growth potential (50% weight)
	upside, err := stock.GetUpside()
	if err != nil {
		fmt.Println("Error:", err)
		panic("Error")
	}

	score += minFloat(upside*2, 100) // Maximum 100 points

	// 2. Positive classifications (30%)
	for _, classification := range stock.Classifications {
		switch classification {
		case "Potential Growth":
			score += 30
		case "Bullish Signal":
			score += 25
		case "New Coverage":
			score += 20
		case "Analyst Positive":
			score += 15
		case "Tech":
			score += 10
		case "Biotech":
			score += 8
		}
	}

	// 3. Analyst ratings (20%)
	switch stock.RatingTo {
	case "Strong-Buy":
		score += 40
	case "Outperform":
		score += 30
	case "Buy":
		score += 20
	}

	return score
}

// getRationale generates a rationale for recommending a stock based on its attributes.
func getRationale(stock domain.Stock) string {
	reasons := []string{}

	upside, err := stock.GetUpside()
	if err != nil {
		fmt.Println("Error:", err)
		panic("Error")
	}

	if upside > 10 {
		reasons = append(reasons,
			"Potential of "+strconv.FormatFloat(upside, 'f', 1, 64)+"%")
	}

	for _, classification := range stock.Classifications {
		switch classification {
		case "Bullish Signal":
			reasons = append(reasons, "Recent upgrade")
		case "New Coverage":
			reasons = append(reasons, "New coverage")
		case "Tech":
			reasons = append(reasons, "Technology sector")
		case "Biotech":
			reasons = append(reasons, "Biotechnology sector")
		}
	}

	if len(reasons) == 0 {
		return "Solid fundamentals"
	}
	return strings.Join(reasons, ", ")
}

// minFloat returns the smaller of two float64 values.
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
