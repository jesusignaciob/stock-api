package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/service"
)

// Mock implementations
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) FindByClassification(ctx context.Context, classification string) ([]domain.Stock, error) {
	args := m.Called(ctx, classification)
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) Find(ctx context.Context, pagination domain.PaginationParams, filters domain.Filters) ([]domain.Stock, error) {
	args := m.Called(ctx, pagination, filters)
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockRepository) Count(ctx context.Context, filters domain.Filters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockStockRepository) FindAll(ctx context.Context, order string, page, limit int) ([]domain.Stock, error) {
	args := m.Called(ctx, order, page, limit)
	return args.Get(0).([]domain.Stock), args.Error(1)
}

func (m *MockStockRepository) FindByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*domain.Stock), args.Error(1)
}

func (m *MockStockRepository) Delete(ctx context.Context, stock *domain.Stock, id uint) error {
	args := m.Called(ctx, stock, id)
	return args.Error(0)
}

func (m *MockStockRepository) SaveBatch(ctx context.Context, stocks []*domain.Stock) error {
	args := m.Called(ctx, stocks)
	return args.Error(0)
}

type MockFieldValidator struct {
	mock.Mock
}

func (m *MockFieldValidator) IsValidField(field string) bool {
	args := m.Called(field)
	return args.Bool(0)
}

func (m *MockFieldValidator) GetAllValidFields() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// Test for Find
func TestFind(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockValidator := new(MockFieldValidator)
	service := service.NewStockService(mockRepo, mockValidator)

	ctx := context.Background()
	pagination := domain.PaginationParams{Page: 1, PageSize: 1, SortOrder: 1, SortField: "company"}
	filters := domain.Filters{"ticker": domain.Filter{Value: "MOM", MatchMode: "contains"}}

	mockValidator.On("IsValidField", "company").Return(true)
	mockValidator.On("IsValidField", "ticker").Return(true)
	mockRepo.On("Find", ctx, pagination, filters).Return([]domain.Stock{{Ticker: "MOMO"}}, nil)
	mockRepo.On("Count", ctx, filters).Return(1, nil)

	stocks, total, err := service.Find(ctx, pagination, filters)

	assert.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "MOMO", stocks[0].Ticker)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestFind_InvalidSortField(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockValidator := new(MockFieldValidator)
	service := service.NewStockService(mockRepo, mockValidator)

	ctx := context.Background()
	pagination := domain.PaginationParams{Page: 1, PageSize: 1, SortOrder: 1, SortField: "invalid_field"}
	filters := domain.Filters{}

	mockValidator.On("IsValidField", "invalid_field").Return(false)

	stocks, total, err := service.Find(ctx, pagination, filters)

	assert.Error(t, err)
	assert.Nil(t, stocks)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "invalid sort field: invalid_field")

	mockValidator.AssertExpectations(t)
}

func TestFind_InvalidFilterField(t *testing.T) {
	mockRepo := new(MockStockRepository)
	mockValidator := new(MockFieldValidator)
	service := service.NewStockService(mockRepo, mockValidator)

	ctx := context.Background()
	pagination := domain.PaginationParams{Page: 1, PageSize: 10}
	filters := domain.Filters{"invalid_field": domain.Filter{Value: "value", MatchMode: "contains"}}

	mockValidator.On("IsValidField", "time").Return(true)
	mockValidator.On("IsValidField", "invalid_field").Return(false)

	stocks, total, err := service.Find(ctx, pagination, filters)

	assert.Error(t, err)
	assert.Nil(t, stocks)
	assert.Equal(t, 0, total)
	assert.EqualError(t, err, "invalid filter field: invalid_field")

	mockValidator.AssertExpectations(t)
}
