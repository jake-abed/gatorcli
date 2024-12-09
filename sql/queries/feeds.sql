-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: GetFeeds :many

SELECT feeds.name, feeds.url, users.name FROM feeds
  INNER JOIN users ON users.id = feeds.user_id;

-- name: CreateFeedFollow :one

WITH insert_feed_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
  )
  RETURNING *
) SELECT
  insert_feed_follow.*,
  feeds.name AS feed_name,
  users.name AS user_name
FROM insert_feed_follow
INNER JOIN feeds ON insert_feed_follow.feed_id = feeds.id
INNER JOIN users ON insert_feed_follow.user_id = users.id;
