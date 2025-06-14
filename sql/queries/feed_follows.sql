-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (user_id, feed_id)
    VALUES ($1,$2)
    ON CONFLICT (user_id, feed_id) DO NOTHING
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users on users.id = inserted_feed_follow.user_id
INNER JOIN feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT 
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;