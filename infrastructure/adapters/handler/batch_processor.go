package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"stock-api/infrastructure/core/domain"
	"stock-api/infrastructure/core/port"
)

type BatchProcessor struct {
	apiClient             port.APIClient
	repo                  port.StockRepository
	classificationService port.ClassificationService
	// Configuration
	batchSize int
	jwtToken  string
	apiDelay  time.Duration
}

// NewBatchProcessor creates a new instance of BatchProcessor
func NewBatchProcessor(
	apiClient port.APIClient,
	repo port.StockRepository,
	classificationService port.ClassificationService,
	batchSize int,
	token string,
	apiDelay time.Duration,
) *BatchProcessor {
	return &BatchProcessor{
		apiClient:             apiClient,
		repo:                  repo,
		classificationService: classificationService,
		// Configuration
		batchSize: batchSize,
		jwtToken:  token,
		apiDelay:  apiDelay,
	}
}

// ProcessStocks processes paginated stocks by ticker
func (bp *BatchProcessor) ProcessStocks(ctx context.Context) error {
	var (
		batch      []*domain.Stock
		lastTicker string
		total      int
		startTime  = time.Now()
	)

	for {
		// Fetch data from the API
		items, nextPage, err := bp.apiClient.FetchStocks(ctx, bp.jwtToken, lastTicker)
		if err != nil {
			return fmt.Errorf("error fetching stocks: %w", err)
		}

		if len(items) == 0 {
			break // No more data
		}

		// Update the last ticker for the next page
		lastTicker = nextPage
		total += len(items)
		batch = append(batch, items...)

		// Save in batches when the defined size is reached
		if len(batch) >= bp.batchSize {
			// Classify and save the current batch
			bp.classificationService.ClassifyBatch(batch)

			if err := bp.saveStocksBatch(ctx, batch); err != nil {
				return fmt.Errorf("error saving batch: %w", err)
			}
			batch = batch[:0] // Clear the batch while retaining capacity
		}

		// Log progress
		log.Printf("Processed %d items, last ticker: %s", total, lastTicker)

		// If there are no more pages, exit
		if nextPage == "" {
			break
		}

		// Wait before the next request
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(bp.apiDelay):
			continue
		}
	}

	// Save remaining data
	if len(batch) > 0 {
		// Classify and save the last batch
		bp.classificationService.ClassifyBatch(batch)

		// Save the batch after classification
		if err := bp.saveStocksBatch(ctx, batch); err != nil {
			return fmt.Errorf("error saving final batch: %w", err)
		}
	}

	log.Printf("Process completed. Total items processed: %d in %v", total, time.Since(startTime))
	return nil
}

// saveStocksBatch saves a batch of stocks to the repository
func (bp *BatchProcessor) saveStocksBatch(ctx context.Context, batch []*domain.Stock) error {
	log.Printf("Saving batch of %d stocks", len(batch))
	return bp.repo.SaveBatch(ctx, batch)
}
