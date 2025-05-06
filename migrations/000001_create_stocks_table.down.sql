-- Drop indexes if they exist
DROP INDEX IF EXISTS idx_stocks_ticker;

DROP INDEX IF EXISTS idx_stocks_time;

DROP INDEX IF EXISTS idx_stocks_deleted_at;

-- Drop the table stocks if it exists
DROP TABLE IF EXISTS stocks;