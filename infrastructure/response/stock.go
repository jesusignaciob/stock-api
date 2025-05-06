package response

import (
	"time"

	"stock-api/infrastructure/core/domain"
)

// StockResponse representa la estructura esperada por el frontend
type StockResponse struct {
	Items        []StockItem `json:"items"`
	Page         int         `json:"page"`
	TotalRecords int         `json:"totalRecords,omitempty"`
	OrderBy      string      `json:"order_by"`
}

// StockItem es la representación Go de tu interfaz TypeScript
type StockItem struct {
	Ticker          string   `json:"ticker"`
	TargetFrom      string   `json:"target_from"`
	TargetTo        string   `json:"target_to"`
	Company         string   `json:"company"`
	Action          string   `json:"action"`
	Brokerage       string   `json:"brokerage"`
	RatingFrom      string   `json:"rating_from"`
	RatingTo        string   `json:"rating_to"`
	Time            string   `json:"time"`
	Classifications []string `json:"classifications"`
}

func ToStockResponse(
	stocks []domain.Stock,
	page int,
	totalRecords int,
	orderBy string,
) StockResponse {
	items := make([]StockItem, len(stocks))

	for i := range stocks {
		stock := &stocks[i]
		items[i] = StockItem{
			Ticker:          stock.Ticker,
			TargetFrom:      stock.TargetFrom,
			TargetTo:        stock.TargetTo,
			Company:         stock.Company,
			Action:          stock.Action,
			Brokerage:       stock.Brokerage,
			RatingFrom:      stock.RatingFrom,
			RatingTo:        stock.RatingTo,
			Time:            stock.Time.Format(time.RFC3339), // Formato estándar
			Classifications: stock.Classifications,
		}
	}

	return StockResponse{
		Items:        items,
		Page:         page,
		TotalRecords: totalRecords,
		OrderBy:      orderBy,
	}
}
