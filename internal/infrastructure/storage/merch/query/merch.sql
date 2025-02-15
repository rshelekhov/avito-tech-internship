-- name: GetMerchByName :one
SELECT
    id,
    name,
    price
FROM merch
WHERE name = $1
  AND deleted_at IS NULL;

-- name: AddToInventory :exec
INSERT INTO purchases (id, user_id, merch_id, created_at)
VALUES ($1, $2, $3, $4);