-- +goose Up
CREATE TABLE IF NOT EXISTS "users"(
    "id" UUID PRIMARY KEY,
    "full_name" VARCHAR(64) NOT NULL,
    "role" UUID NOT NULL,
    "username" VARCHAR(64) NOT NULL,
    "email" VARCHAR(64) NOT NULL,
    "password" TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "users";