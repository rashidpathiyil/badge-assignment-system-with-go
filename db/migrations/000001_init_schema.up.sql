CREATE TABLE event_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    schema JSONB,             -- JSON schema defining the expected event payload structure
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE condition_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    evaluation_logic TEXT,    -- Reference or JSON definition for dynamic mapping
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE badges (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE badge_criteria (
    id SERIAL PRIMARY KEY,
    badge_id INTEGER REFERENCES badges(id) ON DELETE CASCADE,
    flow_definition JSONB NOT NULL,  -- Dynamic rule definition using JSON operators
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_badges (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,   -- External user identifier
    badge_id INTEGER REFERENCES badges(id) ON DELETE CASCADE,
    awarded_at TIMESTAMP DEFAULT NOW(),
    metadata JSONB                   -- Additional details (e.g., event IDs, evaluation metrics)
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    event_type_id INTEGER REFERENCES event_types(id) ON DELETE SET NULL,
    user_id VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,          -- The raw event data
    occurred_at TIMESTAMP DEFAULT NOW()
);

-- Indices to improve query performance
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_event_type_id ON events(event_type_id);
CREATE INDEX idx_events_occurred_at ON events(occurred_at);
CREATE INDEX idx_user_badges_user_id ON user_badges(user_id);
CREATE INDEX idx_user_badges_badge_id ON user_badges(badge_id);
CREATE INDEX idx_badge_criteria_badge_id ON badge_criteria(badge_id); 
