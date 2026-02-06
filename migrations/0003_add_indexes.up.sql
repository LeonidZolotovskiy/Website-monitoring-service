CREATE INDEX IF NOT EXISTS idx_site_checks_checked_at
    ON site_checks (checked_at DESC);

CREATE INDEX IF NOT EXISTS idx_site_checks_site_id
    ON site_checks (site_id);

CREATE INDEX idx_check_results_site_id ON site_checks(site_id);