package service

import (
	"testing"

	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/service"
)

func TestClassificationService_Classify(t *testing.T) {
	service := service.NewClassificationService()

	tests := []struct {
		name           string
		stock          *domain.Stock
		expectedLabels []string
	}{
		{
			name: "Biotech sector classification",
			stock: &domain.Stock{
				Company: "Medical Therapeutics Inc.",
			},
			expectedLabels: []string{"Biotech"},
		},
		{
			name: "Tech sector classification",
			stock: &domain.Stock{
				Company: "Tech Solutions Group",
			},
			expectedLabels: []string{"Tech"},
		},
		{
			name: "Financial sector classification",
			stock: &domain.Stock{
				Company: "Global Financial Advisors",
			},
			expectedLabels: []string{"Financial"},
		},
		{
			name: "Energy sector classification",
			stock: &domain.Stock{
				Company: "Petroleum Gas Resources",
			},
			expectedLabels: []string{"Energy"},
		},
		{
			name: "Other sector classification",
			stock: &domain.Stock{
				Company: "Unknown Company",
			},
			expectedLabels: []string{"Other Sector"},
		},
		{
			name: "High-Risk Speculative classification",
			stock: &domain.Stock{
				TargetFrom: "$100.00",
				TargetTo:   "$70.00",
			},
			expectedLabels: []string{"High-Risk Speculative"},
		},
		{
			name: "Potential Growth classification",
			stock: &domain.Stock{
				TargetFrom: "$100.00",
				TargetTo:   "$120.00",
			},
			expectedLabels: []string{"Potential Growth"},
		},
		{
			name: "Bullish Signal classification",
			stock: &domain.Stock{
				Action: "Upgraded by analysts",
			},
			expectedLabels: []string{"Bullish Signal"},
		},
		{
			name: "Bearish Signal classification",
			stock: &domain.Stock{
				Action: "Downgraded by analysts",
			},
			expectedLabels: []string{"Bearish Signal"},
		},
		{
			name: "New Coverage classification",
			stock: &domain.Stock{
				Action: "Initiated by analysts",
			},
			expectedLabels: []string{"New Coverage"},
		},
		{
			name: "Analyst Positive classification",
			stock: &domain.Stock{
				RatingTo: "Buy",
			},
			expectedLabels: []string{"Analyst Positive"},
		},
		{
			name: "Analyst Negative classification",
			stock: &domain.Stock{
				RatingTo: "Sell",
			},
			expectedLabels: []string{"Analyst Negative"},
		},
		// {
		// 	name: "Neutral classification",
		// 	stock: &domain.Stock{
		// 		Company: "",
		// 	},
		// 	expectedLabels: []string{"Neutral"},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.Classify(tt.stock)
			for _, label := range tt.expectedLabels {
				found := false
				for _, classification := range tt.stock.Classifications {
					if classification == label {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected classification %s not found in %v", label, tt.stock.Classifications)
				}
			}
		})
	}
}
