package repository

import (
	"context"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"stock-api/infrastructure/core/domain"
)

// StockBDRepository is the repository responsible for interacting with the database
// for operations related to the Stock model.
type StockBDRepository struct {
	db *gorm.DB
}

// NewStockBDRepository creates a new instance of StockBDRepository.
// It takes a GORM database instance as a parameter.
func NewStockBDRepository(db *gorm.DB) *StockBDRepository {
	repository := &StockBDRepository{db: db}
	return repository
}

// Create inserts a new stock record into the database.
// It takes a context and a pointer to a Stock object as parameters.
func (r *StockBDRepository) Create(ctx context.Context, stock *domain.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

// Delete removes a stock record from the database by its ID.
// It takes a context, a pointer to a Stock object, and the ID of the stock to delete.
func (r *StockBDRepository) Delete(ctx context.Context, stock *domain.Stock, id uint) error {
	return r.db.WithContext(ctx).Delete(stock, id).Error
}

// FindAll retrieves a paginated list of stocks from the database.
// It takes a context, the order of sorting, the page number, and the limit of records per page.
// Returns a slice of Stock objects and an error if any.
func (r *StockBDRepository) FindAll(ctx context.Context, order string, page, limit int) ([]domain.Stock, error) {
	stocks := []domain.Stock{}
	if err := r.db.WithContext(ctx).Order(order).Offset((page - 1) * limit).Limit(limit).Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}

// FindByTicker retrieves a stock record from the database by its ticker.
// It takes a context and the ticker string as parameters.
// Returns a pointer to a Stock object and an error if any.
func (r *StockBDRepository) FindByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	var stock domain.Stock
	if err := r.db.WithContext(ctx).Where("ticker = ?", ticker).First(&stock).Error; err != nil {
		return nil, err
	}
	return &stock, nil
}

// FindByClassification retrieves all stocks that match a specific classification.
// It takes a context and the classification string as parameters.
// Returns a slice of Stock objects and an error if any.
func (r *StockBDRepository) FindByClassification(ctx context.Context, classification string) ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.WithContext(ctx).
		Where("classifications @> ?", pq.StringArray{classification}).
		Find(&stocks).Error
	return stocks, err
}

// SaveBatch inserts multiple stock records into the database in batches.
// It takes a context and a slice of pointers to Stock objects as parameters.
func (r *StockBDRepository) SaveBatch(ctx context.Context, data []*domain.Stock) error {
	return r.db.WithContext(ctx).CreateInBatches(data, len(data)).Error
}
