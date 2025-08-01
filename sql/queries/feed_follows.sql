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

-- name: GetFeedFollowsForUser :many
-- Add a GetFeedFollowsForUser query. It should return all the feed follows for a given user, and include the names of the feeds and user in the result.
SELECT
    ff.id,
    ff.created_at,
    ff.updated_at,
    ff.user_id,
    ff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM
    feed_follows ff
    INNER JOIN users u ON ff.user_id = u.id
    INNER JOIN feeds f ON ff.feed_id = f.id
WHERE
    u.name = $1;

-- name: DeleteFeedFollowByUserAndFeedUrl :one
-- Delete a feed follow record by user and feed url combination.
DELETE FROM feed_follows
WHERE
    user_id = (
        SELECT user_id
        FROM users
        WHERE
            users.name = $1
    )
    AND feed_id = (
        SELECT feed_id
        FROM feeds
        WHERE
            feeds.url = $2
    ) RETURNING *;