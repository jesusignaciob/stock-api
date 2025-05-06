CREATE TABLE
    stocks (
        id SERIAL PRIMARY KEY,
        created_at TIMESTAMP
        WITH
            TIME ZONE,
            updated_at TIMESTAMP
        WITH
            TIME ZONE,
            deleted_at TIMESTAMP
        WITH
            TIME ZONE,
            ticker VARCHAR(10) NOT NULL,
            target_from VARCHAR(20),
            target_to VARCHAR(20),
            company VARCHAR(255) NOT NULL,
            action VARCHAR(100),
            brokerage VARCHAR(255) NOT NULL,
            rating_from VARCHAR(50),
            rating_to VARCHAR(50),
            classifications TEXT[] DEFAULT '{"Neutral"}',
            time TIMESTAMP
        WITH
            TIME ZONE NOT NULL
    );

-- Create indexes as defined in the GORM model
CREATE INDEX idx_stocks_ticker ON stocks (ticker);

CREATE INDEX idx_stocks_time ON stocks (time);

CREATE INDEX idx_stocks_deleted_at ON stocks (deleted_at);