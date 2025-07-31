-- name: CreateFeedFollow :one
INSERT INTO
    feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5) RETURNING feed_follows.*,
    (
        SELECT users.name
        FROM users
        WHERE
            users.id = feed_follows.user_id
    ) AS user_name,
    (
        SELECT feeds.name
        FROM feeds
        WHERE
            feeds.id = feed_follows.feed_id
    ) AS feed_name;