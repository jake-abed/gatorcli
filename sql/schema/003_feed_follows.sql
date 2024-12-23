-- +goose Up
CREATE TABLE feed_follows (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
  feed_id uuid NOT NULL REFERENCES feeds ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (feed_id) REFERENCES feeds(id),
  UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;

