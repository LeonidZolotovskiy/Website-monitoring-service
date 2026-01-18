CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE sites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID NOT NULL,
    http_code INT NOT NULL,
    response_time_ms INT,
    is_available BOOLEAN NOT NULL,
    error_message TEXT,
    checked_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_checks_site
        FOREIGN KEY (site_id)
        REFERENCES sites(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_checks_checked_at ON checks(checked_at);
CREATE INDEX idx_checks_site_id ON checks(site_id);
