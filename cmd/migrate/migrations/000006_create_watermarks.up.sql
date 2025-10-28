CREATE TABLE watermarks (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN watermarks.key IS 'Unique identifier for watermark';
COMMENT ON COLUMN watermarks.value IS 'Value of watermark';
COMMENT ON COLUMN watermarks.updated_at IS 'Timestamp of last update';

CREATE INDEX idx_watermarks_key ON watermarks(key);
