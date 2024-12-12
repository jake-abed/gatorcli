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

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
  USING feeds, users
  WHERE users.name = $1 AND feeds.url = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
  SET updated_at = $1, last_fetched_at = $1
  WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
  ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;
