-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash, balance, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserBalanceByID :one
SELECT id, balance as coins
FROM users
WHERE users.id = $1
    AND deleted_at IS NULL;

-- name: GetUserInventory :many
SELECT m.name as type,
       COUNT(*) as quantity
FROM purchases p
    JOIN merch m ON p.merch_id = m.id AND m.deleted_at IS NULL
WHERE p.user_id = $1
GROUP BY m.name;

-- name: GetReceivedTransactions :many
SELECT
    sender.username as from_user,
    t.receiver_id as to_user,
    t.amount,
    t.created_at as date
FROM transactions t
    JOIN users sender ON t.sender_id = sender.id AND sender.deleted_at IS NULL
WHERE t.receiver_id = $1
    AND t.transaction_type_id = 0;  -- only coin transfers

-- name: GetSentTransactions :many
SELECT
    t.sender_id as from_user,
    receiver.username as to_user,
    t.amount,
    t.created_at as date
FROM transactions t
    JOIN users receiver ON t.receiver_id = receiver.id AND receiver.deleted_at IS NULL
WHERE t.sender_id = $1
    AND t.transaction_type_id = 0;  -- only coin transfers

-- name: GetUserIDByUsername :one
SELECT id
FROM users
WHERE username = $1
    AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT id,
       username,
       password_hash,
       balance,
       created_at,
       updated_at
FROM users
WHERE username = $1
  AND deleted_at IS NULL;