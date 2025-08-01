-- name: CreatePost :one
-- Inserts a new post into the database.
INSERT INTO
    posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES (
        @id,
        @created_at,
        @updated_at,
        @title,
        @url,
        @description,
        @published_at,
        @feed_id
    ) RETURNING *;