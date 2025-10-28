CREATE TABLE settlement_types (
    id CHAR(2) PRIMARY KEY,           
    name VARCHAR(50) NOT NULL,
    slug VARCHAR(50) NOT NULL         
);

COMMENT ON COLUMN settlement_types.id IS 'SEPOMEX settlement type code';
COMMENT ON COLUMN settlement_types.name IS 'Settlement type name';
COMMENT ON COLUMN settlement_types.slug IS 'Settlement type slug';

CREATE INDEX idx_settlement_types_name ON settlement_types(name);
CREATE INDEX idx_settlement_types_slug ON settlement_types(slug);
