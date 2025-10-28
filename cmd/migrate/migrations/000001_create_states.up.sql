CREATE TABLE states (
    id CHAR(2) PRIMARY KEY,           
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL        
);

COMMENT ON COLUMN states.id IS 'INEGI state code';
COMMENT ON COLUMN states.name IS 'State name';
COMMENT ON COLUMN states.slug IS 'State slug';

CREATE INDEX idx_states_name ON states(name);
CREATE INDEX idx_states_slug ON states(slug);
