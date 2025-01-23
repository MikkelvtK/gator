-- name: CreatePost :exec
INSERT INTO posts (
    id,
    title,
    url,
    feed_id,
    description,
    published_at,
    created_at,
    updated_at
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 );

-- name: GetPostsForUser :many
SELECT p.*
FROM posts p
INNER JOIN feed_follows ff
    ON p.feed_id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2;
