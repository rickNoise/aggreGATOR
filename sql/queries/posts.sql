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

-- name: GetPostsForUser :many
-- Returns all saved posts for all feeds followed by the current user.
SELECT p.*, f.name
FROM
    posts p
    JOIN feeds f ON p.feed_id = f.id
    JOIN feed_follows ff ON f.id = ff.feed_id
    JOIN users u ON ff.user_id = u.id
WHERE
    u.id = sqlc.arg (user_id)
ORDER BY p.updated_at DESC NULLS LAST
LIMIT sqlc.arg (num_posts_limit);