-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
  $1,
  $2,
  $3,
  $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE name = $1 LIMIT 1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetFeedFollowsForUser :many
SELECT ff.id, u.name as user_name, f.name as feed_name FROM users as U
  INNER JOIN feed_follows as ff ON ff.user_id = u.id
  INNER JOIN feeds as f ON ff.feed_id = f.id
  WHERE u.name = $1;
