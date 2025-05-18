package repository

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lib/pq"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"stock-api/infrastructure/core/domain"
)

// In-memory cache for Count results
var (
	countCache sync.Map
	countGroup singleflight.Group
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

// Find retrieves a list of stocks from the database based on the provided pagination
// parameters and filters. It applies filtering, ordering, and pagination to the query
// before executing it.
//
// Parameters:
//   - ctx: The context for managing request deadlines, cancellation signals, and other
//     request-scoped values.
//   - pagination: An instance of domain.PaginationParams containing pagination details
//     such as page number and page size.
//   - filters: A map of field names to filter values, used to apply filtering criteria
//     to the query.
//
// Returns:
//   - []domain.Stock: A slice of domain.Stock objects that match the query criteria.
//   - error: An error object if the query fails, or nil if the operation is successful.
func (r *StockBDRepository) Find(ctx context.Context, pagination domain.PaginationParams, filters domain.Filters) ([]domain.Stock, error) {
	var stocks []domain.Stock
	query := r.db.WithContext(ctx)

	var err error
	for field, filter := range filters {
		query, err = applyFilter(query, field, filter)
		if err != nil {
			return nil, err
		}
	}

	query = applyOrder(query, pagination)
	query = applyPagination(query, pagination)

	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
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

// Count returns the number of stocks in the database that match the provided filters.
// It uses an in-memory cache with the serialized and hashed filters as the key.
// Uses singleflight to avoid duplicate DB queries for the same key under concurrency.
func (r *StockBDRepository) Count(ctx context.Context, filters domain.Filters) (int, error) {
	cacheKey := getCacheKey(filters)

	// Try to get from cache
	if v, ok := countCache.Load(cacheKey); ok {
		if cachedCount, ok := v.(int); ok {
			return cachedCount, nil
		}
	}

	// Use singleflight to avoid duplicate DB queries for the same key
	val, err, _ := countGroup.Do(cacheKey, func() (interface{}, error) {
		var count int64
		var err error
		query := r.db.WithContext(ctx)
		for field, filter := range filters {
			query, err = applyFilter(query, field, filter)
			if err != nil {
				return 0, err
			}
		}
		err = query.Model(&domain.Stock{}).Count(&count).Error
		if err == nil {
			countCache.Store(cacheKey, int(count))
		}
		return int(count), err
	})
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// getCacheKey serializes and hashes the filters to generate a unique cache key.
func getCacheKey(filters domain.Filters) string {
	b, _ := json.Marshal(filters)
	hash := sha256.Sum256(b)
	return fmt.Sprintf("%x", hash)
}

// applyFilter applies a filter to the GORM query based on the field and filter criteria.
// Supports array fields (like classifications) and various match modes for both array and scalar fields.
// Returns the updated query and an error if the filter is invalid.
func applyFilter(query *gorm.DB, field string, filter domain.Filter) (*gorm.DB, error) {
	if field == "classifications" {
		var arr []string
		switch v := filter.Value.(type) {
		case string:
			arr = []string{v}
		case []string:
			arr = v
		case []interface{}:
			arr = make([]string, len(v))
			for i, val := range v {
				s, ok := val.(string)
				if !ok {
					return query, fmt.Errorf("invalid value type in array for field %s: expected string", field)
				}
				arr[i] = s
			}
		default:
			return query, fmt.Errorf("invalid value type for field %s: expected string or []string", field)
		}
		// Handle the classification filter
		switch filter.MatchMode {
		case "equals":
			// Checks if the array is exactly equal (same elements in the same order)
			query = query.Where("classifications = ?", pq.Array(arr))
		case "contains":
			// Checks if the array contains at least the specified value
			query = query.Where("classifications @> ?", pq.Array(arr))
		case "contained":
			// Checks if the array is contained in the specified value
			query = query.Where("classifications <@ ?", pq.Array(arr))
		case "overlap":
			// Checks if there is any common element
			query = query.Where("classifications && ?", pq.Array(arr))
		case "notEquals":
			// Checks if the array is not exactly equal
			query = query.Where("classifications != ?", pq.Array(arr))
		case "notContains":
			// Checks if the array does not contain the specified value
			query = query.Where("NOT classifications @> ?", pq.Array(arr))
		case "notContained":
			// Checks if the array is not contained in the specified value
			query = query.Where("NOT classifications <@ ?", pq.Array(arr))
		case "notOverlap":
			// Checks if there is no common element
			query = query.Where("NOT classifications && ?", pq.Array(arr))
		default:
			// Default mode or error handling
			return nil, fmt.Errorf("unsupported match mode for array field: %s", filter.MatchMode)
		}

		return query, nil
	}
	switch filter.MatchMode {
	case "equals":
		// Checks if the field is equal to the specified value
		query = query.Where(fmt.Sprintf("%s = ?", field), filter.Value)
	case "contains":
		// Checks if the field contains the specified value
		query = query.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", filter.Value))
	case "startsWith":
		// Checks if the field starts with the specified value
		query = query.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%v%%", filter.Value))
	case "endsWith":
		// Checks if the field ends with the specified value
		query = query.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v", filter.Value))
	case "greaterThan":
		// Checks if the field is greater than the specified value
		query = query.Where(fmt.Sprintf("%s > ?", field), filter.Value)
	case "lessThan":
		// Checks if the field is less than the specified value
		query = query.Where(fmt.Sprintf("%s < ?", field), filter.Value)
	default:
		// Default mode or error handling
		return nil, fmt.Errorf("unsupported match mode for array field: %s", filter.MatchMode)
	}

	return query, nil
}

// applyOrder applies sorting to the GORM query based on the pagination parameters.
// Returns the updated query.
func applyOrder(query *gorm.DB, pagination domain.PaginationParams) *gorm.DB {
	if pagination.SortField != "" {
		order := "ASC"
		if pagination.SortOrder == -1 {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", pagination.SortField, order))
	}

	return query
}

// applyPagination applies pagination to the GORM query based on the pagination parameters.
// Returns the updated query.
func applyPagination(query *gorm.DB, pagination domain.PaginationParams) *gorm.DB {
	if pagination.Page > 0 && pagination.PageSize > 0 {
		query = query.Offset((pagination.Page - 1) * pagination.PageSize).Limit(pagination.PageSize)
	}

	return query
}
