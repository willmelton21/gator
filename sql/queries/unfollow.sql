-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1 AND feed_id = (
  SELECT id FROM feeds WHERE url = $2
);
