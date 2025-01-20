-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
  INSERT INTO feed_follows (
      id,
      user_id,
      feed_id,
      created_at,
      updated_at
  ) VALUES ( $1, $2, $3, $4, $5 )
  RETURNING *
)
SELECT
    ff.*,
    u.name AS user_name,
    f.name AS feed_name
FROM inserted_feed_follow ff
INNER JOIN users u
    ON u.id = ff.user_id
INNER JOIN feeds f
    ON f.id = ff.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT 
    ff.*,
    u.name AS user_name,
    f.name AS feed_name
FROM feed_follows ff
INNER JOIN users u
    ON u.id = ff.user_id
INNER JOIN feeds f
    ON f.id = ff.feed_id
HAVING u.name = $1;
