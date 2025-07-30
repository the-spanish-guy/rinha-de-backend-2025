CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE UNLOGGED TABLE payments (
    id UUID NOT NULL DEFAULT uuid_generate_v4()  PRIMARY KEY,
    correlation_id VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20),
    amount DECIMAL(10, 2) NOT NULL,
    processor varchar NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX payments_status ON payments (status);
CREATE INDEX payments_processor ON payments (processor);
CREATE INDEX payments_requested_at ON payments (requested_at);
CREATE INDEX payments_correlation_id ON payments (correlation_id);

CREATE INDEX payments_summary ON payments (processor, requested_at, amount);
CREATE INDEX payments_date_range ON payments (requested_at) WHERE requested_at IS NOT NULL;
