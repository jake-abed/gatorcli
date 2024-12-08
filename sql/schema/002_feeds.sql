-- +goose Up
CREATE TABLE feeds (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name text NOT NULL,
  url text UNIQUE NOT NULL,
  user_id uuid NOT NULL REFERENCES users ON DELETE CASCADE,
  FOREIGN KEY (user_id)
  REFERENCES users(id)
);
-- +goose Down
DROP TABLE feeds;
