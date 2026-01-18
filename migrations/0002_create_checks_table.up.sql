CREATE TABLE IF NOT EXISTS site_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    site_id UUID NOT NULL,
    http_status INT,
    response_time_ms INT,
    is_available BOOLEAN NOT NULL,
    error_message TEXT,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_site
        FOREIGN KEY(site_id)
        REFERENCES sites(id)
        ON DELETE CASCADE
);
