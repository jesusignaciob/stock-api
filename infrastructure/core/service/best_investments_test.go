package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"stock-api/infrastructure/core/domain"
)

func TestGetStockRecommendations(t *testing.T) {
	service := NewBestInvestmentsService()

	mockStocks := []domain.Stock{
		{
			Ticker:          "AAPL",
			Company:         "Apple Inc.",
			Classifications: []string{"Potential Growth", "Bullish Signal"},
			RatingTo:        "Strong-Buy",
			TargetFrom:      "$100.00",
			TargetTo:        "$115.00",
		},
		{
			Ticker:          "TSLA",
			Company:         "Tesla Inc.",
			Classifications: []string{"High-Risk Speculative"},
			RatingTo:        "Buy",
			TargetFrom:      "$200.00",
			TargetTo:        "$240.00",
		},
		{
			Ticker:          "MSFT",
			Company:         "Microsoft Corp.",
			Classifications: []string{"Potential Growth", "New Coverage"},
			RatingTo:        "Outperform",
			TargetFrom:      "$150.00",
			TargetTo:        "$168.00",
		},
	}

	t.Run("should return top recommendations based on score", func(t *testing.T) {
		limit := 2
		recommendations := service.GetStockRecommendations(mockStocks, limit)

		assert.Equal(t, limit, len(recommendations))
		assert.Equal(t, "AAPL", recommendations[0].Ticker)
		assert.Equal(t, "MSFT", recommendations[1].Ticker)
	})

	t.Run("should exclude stocks with problematic classifications", func(t *testing.T) {
		limit := 3
		recommendations := service.GetStockRecommendations(mockStocks, limit)

		for _, rec := range recommendations {
			assert.NotEqual(t, "TSLA", rec.Ticker)
		}
	})

	t.Run("should handle limit greater than available stocks", func(t *testing.T) {
		limit := 10
		recommendations := service.GetStockRecommendations(mockStocks, limit)

		assert.Equal(t, 2, len(recommendations)) // Only 2 valid stocks
	})

	t.Run("should generate correct rationale for recommendations", func(t *testing.T) {
		limit := 1
		recommendations := service.GetStockRecommendations(mockStocks, limit)

		assert.Contains(t, recommendations[0].Rationale, "Potential of 15.0%")
		assert.Contains(t, recommendations[0].Rationale, "Recent upgrade")
	})
}
