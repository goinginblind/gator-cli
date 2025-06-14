-- +goose Up
CREATE INDEX idx_feeds_user_id ON feeds(user_id);

-- -goose Down
DROP INDEX IF EXISTS idx_feeds_user_id;