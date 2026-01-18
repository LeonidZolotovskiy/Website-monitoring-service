CREATE TABLE check_results (
    id SERIAL PRIMARY KEY,
    site_id UUID NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    response_time_ms INT NOT NULL,
    checked_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_check_results_site_id ON check_results(site_id);