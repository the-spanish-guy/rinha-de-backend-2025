CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE UNLOGGED TABLE payments (
    id UUID NOT NULL DEFAULT uuid_generate_v4()  PRIMARY KEY,
    correlation_id UUID NOT NULL,
    status VARCHAR(20),
    amount DECIMAL(10, 2) NOT NULL,
    requested_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX payments_status ON payments (status);
CREATE INDEX payments_correlation_id ON payments (correlation_id);
