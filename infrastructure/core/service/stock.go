package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/port"
)

type StockService struct {
	repo           port.StockRepository
	fieldValidator port.FieldValidator
}

func NewStockService(userRepo port.StockRepository, fieldValidator port.FieldValidator) *StockService {
	return &StockService{repo: userRepo, fieldValidator: fieldValidator}
}

func (s *StockService) RegisterStock(ctx context.Context, stock *domain.Stock) error {
	if stock == nil {
		return errors.New("stock cannot be nil")
	}
	if err := s.repo.Create(ctx, stock); err != nil {
		return err
	}
	return nil
}

func (s *StockService) Find(ctx context.Context, pagination domain.PaginationParams, filters domain.Filters) ([]domain.Stock, int, error) {
	log.Printf("Entro en Find: %v", filters)

	// Validar campo de ordenamiento
	if pagination.SortField != "" && !s.fieldValidator.IsValidField(pagination.SortField) {
		return nil, 0, fmt.Errorf("invalid sort field: %s", pagination.SortField)
	}

	// Validar campos de filtro
	for field := range filters {
		if !s.fieldValidator.IsValidField(field) {
			return nil, 0, fmt.Errorf("invalid filter field: %s", field)
		}
	}

	stocks, err := s.repo.Find(ctx, pagination, filters)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

func (s *StockService) FindAllStocks(ctx context.Context, order string, page, limit int) ([]domain.Stock, error) {
	stocks, err := s.repo.FindAll(ctx, order, page, limit)
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (s *StockService) FindStockByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	stock, err := s.repo.FindByTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, errors.New("stock not found")
	}
	return stock, nil
}

func (s *StockService) DeleteStock(ctx context.Context, stock *domain.Stock, id uint) error {
	if stock == nil {
		return errors.New("stock cannot be nil")
	}
	if err := s.repo.Delete(ctx, stock, id); err != nil {
		return err
	}
	return nil
}
