-- name: CreatePost :one
INSERT INTO posts (title, url, description, published_at, feed_id)
VALUES ($1,$2,$3,$4,$5) 
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsUser :many
SELECT posts.* FROM posts
INNER JOIN feeds ON posts.feed_id = feeds.id
INNER JOIN users ON users.id = feeds.user_id
WHERE users.id = $1
ORDER BY posts.published_at DESC 
LIMIT $2;