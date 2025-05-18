package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"stock-api/infrastructure/core/domain"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("error closing db: %v", err)
		}
	}
	return gormDB, mock, cleanup
}

func TestFindByClassification(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewStockBDRepository(db)
	classification := "Potential Growth"

	rows := sqlmock.NewRows([]string{"id", "ticker", "company", "classifications"}).
		AddRow(1, "AAPL", "Apple Inc.", pq.StringArray{classification})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE classifications @> $1`)).
		WithArgs(pq.StringArray{classification}).
		WillReturnRows(rows)

	stocks, err := repo.FindByClassification(context.Background(), classification)
	assert.NoError(t, err)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
	assert.Contains(t, stocks[0].Classifications, classification)
}

func TestApplyFilter_ClassificationsContains(t *testing.T) {
	db, mock, cleanup := setupTestDB(t)
	defer cleanup()

	classifications := []string{"Tech", "Potential Growth"}

	rows := sqlmock.NewRows([]string{"id", "ticker", "company", "classifications"}).
		AddRow(1, "AAPL", "Apple Inc.", pq.StringArray{"Potential Growth"})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stocks" WHERE classifications && $1 AND "stocks"."deleted_at" IS NULL`)).
		WithArgs(pq.Array(classifications)).
		WillReturnRows(rows)

	query := db.WithContext(t.Context()).Model(&domain.Stock{})
	filter := domain.Filter{
		Value:     classifications,
		MatchMode: "overlap",
	}
	q, err := applyFilter(query, "classifications", filter)
	assert.NoError(t, err)
	var stocks []domain.Stock
	err = q.Find(&stocks).Error
	assert.NoError(t, err)
	assert.Len(t, stocks, 1)
	assert.Equal(t, "AAPL", stocks[0].Ticker)
}
