-- name: RegisterCoinTransfer :exec
INSERT INTO transactions (id, sender_id, receiver_id, transaction_type_id, amount, created_at)
VALUES (
           @id,
           @sender_id,
           @receiver_id,
           (SELECT id FROM transaction_types WHERE title = @transaction_type),
           @amount,
           @created_at
       );

-- name: UpdateUserCoins :exec
UPDATE users
SET
    balance = $2,
    updated_at = now()
WHERE id = $1
  AND deleted_at IS NULL;