-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id)
VALUES ($1,$2,$3) 
RETURNING *;

-- name: GetLatestFeed :one
SELECT * FROM feeds 
WHERE user_id = $1
ORDER BY created_at DESC 
LIMIT 1;

-- name: GetFeedsWithUNames :many
SELECT feeds.name, feeds.url, users.name FROM feeds
JOIN users ON user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds 
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET last_fetched_at = now(), updated_at = now() 
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;