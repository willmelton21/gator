-- name: MarkFeedFetched :exec
update feeds
SET last_fetched_At = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
where id = $1;


-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_At NULLS FIRST, last_fetched_at ASC
LIMIT 1;

