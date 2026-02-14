CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    current_stock INTEGER NOT NULL,
    minimum_stock INTEGER NOT NULL,
    average_daily_sales INTEGER NOT NULL,
    lead_time_days INTEGER NOT NULL,
    unit_cost NUMERIC(10,2) NOT NULL,
    criticality_level INTEGER NOT NULL
);
