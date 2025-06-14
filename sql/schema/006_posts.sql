-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    created_at TIMESTAMP NOT NULL DEFAULT now(), 
    updated_at TIMESTAMP NOT NULL DEFAULT now(), 
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;