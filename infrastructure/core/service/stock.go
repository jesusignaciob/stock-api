package service

import (
	"context"
	"errors"

	"stock-api/infrastructure/adapters/repository"
	"stock-api/infrastructure/core/domain"
)

type StockService struct {
	repo *repository.StockBDRepository
}

func NewStockService(userRepo *repository.StockBDRepository) *StockService {
	return &StockService{repo: userRepo}
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
