package port

import (
	"context"

	"stock-api/infrastructure/core/domain"
)

type StockRepository interface {
	Create(ctx context.Context, stock *domain.Stock) error
	Delete(ctx context.Context, stock *domain.Stock, id uint) error
	Find(ctx context.Context, pagination domain.PaginationParams, filters domain.Filters) ([]domain.Stock, error)
	FindAll(ctx context.Context, order string, page, limit int) ([]domain.Stock, error)
	FindByTicker(ctx context.Context, ticker string) (*domain.Stock, error)
	FindByClassification(ctx context.Context, classification string) ([]domain.Stock, error)
	SaveBatch(ctx context.Context, data []*domain.Stock) error
	Count(ctx context.Context, filters domain.Filters) (int, error)
}

type FieldValidator interface {
	IsValidField(field string) bool
	GetAllValidFields() []string
}

type StockService interface {
	RegisterStock(ctx context.Context, stock *domain.Stock) error
	FindStockByTicker(ctx context.Context, ticker string) (*domain.Stock, error)
	DeleteStock(ctx context.Context, stock *domain.Stock, id uint) error
	Find(ctx context.Context, pagination domain.PaginationParams, filters domain.Filters) ([]domain.Stock, int, error)
	FindAllStocks(ctx context.Context, order string, page int, limit int) ([]domain.Stock, error)
}

type ClassificationService interface {
	Classify(stock *domain.Stock)
	ClassifyBatch(batch []*domain.Stock)
}

type BestInvestmentsService interface {
	GetStockRecommendations(batch []domain.Stock, limit int) []domain.Recommendation
}

type APIClient interface {
	FetchStocks(ctx context.Context, jwtToken string, lastTicker string) ([]*domain.Stock, string, error)
}
