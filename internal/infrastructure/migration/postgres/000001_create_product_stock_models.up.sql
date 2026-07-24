CREATE TABLE product_stock_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    current_stock BIGINT NOT NULL,
    minimum_stock BIGINT NOT NULL,
    average_daily_sales BIGINT NOT NULL,
    lead_time_days BIGINT NOT NULL,
    unit_cost NUMERIC(10,2) NOT NULL,
    criticality_level BIGINT NOT NULL
);
