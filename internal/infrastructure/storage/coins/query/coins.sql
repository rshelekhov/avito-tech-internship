-- name: UpdateUserCoins :exec
UPDATE users
SET
    balance = $2,
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL;