-- name: CreateFeed :one
INSERT INTO
    feeds (
        id,
        created_at,
        updated_at,
        name,
        url,
        user_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetFeeds :many
SELECT
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
FROM feeds;

-- name: GetFeedByUrl :one
SELECT
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
FROM feeds
WHERE
    feeds.url = $1;

-- name: MarkFeedFetched :exec
-- Sets the last_fetched_at and updated_at columns to the current time for a given feed by feed id.
UPDATE feeds
SET
    updated_at = NOW(),
    last_fetched_at = NOW()
WHERE
    feeds.id = sqlc.arg (feed_id);

-- name: GetNextFeedToFetch :one
-- Returns the next feed we should fetch posts from. Always fetch the oldest one first.
SELECT
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id,
    last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;