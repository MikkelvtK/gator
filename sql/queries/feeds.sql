-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    name,
    url,
    user_id,
    created_at,
    updated_at
) VALUES ( $1, $2, $3, $4, $5, $6 )
RETURNING *;

-- name: GetFeeds :many
SELECT 
  f.id,
  f.name AS feed_name, 
  f.url, 
  u.name AS user_name,
  f.created_at,
  f.updated_at
FROM feeds f
INNER JOIN users u
ON f.user_id = u.id;

-- name: GetFeedByUrl :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(),
    last_fetched_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at DESC NULLS FIRST
LIMIT 1;
