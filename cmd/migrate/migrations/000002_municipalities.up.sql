CREATE TABLE municipalities (
    id CHAR(3) PRIMARY KEY,           
    name VARCHAR(100) NOT NULL,       
    slug VARCHAR(100) NOT NULL,       
    state_id CHAR(2) NOT NULL REFERENCES states(id)  
);

COMMENT ON COLUMN municipalities.id IS 'INEGI municipality code';
COMMENT ON COLUMN municipalities.name IS 'Municipality name';
COMMENT ON COLUMN municipalities.slug IS 'Municipality slug';
COMMENT ON COLUMN municipalities.state_id IS 'Reference to state';

CREATE INDEX idx_municipalities_name ON municipalities(name);
CREATE INDEX idx_municipalities_slug ON municipalities(slug);
CREATE INDEX idx_municipalities_state_id ON municipalities(state_id);