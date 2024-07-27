CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    status BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP
);


CREATE TABLE IF NOT EXISTS message_stats (
    id SERIAL PRIMARY KEY,
    processed_count INT NOT NULL,
    last_processed_message_id  INT REFERENCES messages(id),
    updated_at TIMESTAMP
);


CREATE INDEX idx_message_id ON messages(id);
CREATE INDEX idx_stats_id ON message_stats(id);