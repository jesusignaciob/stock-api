package port

import (
	"context"

	"stock-api/infrastructure/core/domain"
)

type StockRepository interface {
	Create(ctx context.Context, stock *domain.Stock) error
	Delete(ctx context.Context, stock *domain.Stock, id uint) error
	FindAll(ctx context.Context, order string, page, limit int) ([]domain.Stock, error)
	FindByTicker(ctx context.Context, ticker string) (*domain.Stock, error)
	FindByClassification(ctx context.Context, classification string) ([]domain.Stock, error)
	SaveBatch(ctx context.Context, data []*domain.Stock) error
}

type StockService interface {
	FindAllStocks(ctx context.Context, order string, page int, limit int) ([]domain.Stock, error)
}

type ClassificationService interface {
	Classify(stock *domain.Stock)
	ClassifyBatch(batch []*domain.Stock)
}

type APIClient interface {
	FetchStocks(ctx context.Context, jwtToken string, lastTicker string) ([]*domain.Stock, string, error)
}
