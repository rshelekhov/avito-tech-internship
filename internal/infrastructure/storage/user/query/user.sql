-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash, balance, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUserByName :one
SELECT id,
       username,
       password_hash,
       balance,
       created_at,
       updated_at
FROM users
WHERE username = $1
  AND deleted_at IS NULL;

-- name: GetUserInfoByID :one
WITH user_balance AS (
    SELECT id, balance as coins
    FROM users
    WHERE users.id = $1
      AND deleted_at IS NULL
),
user_inventory AS (
    SELECT m.name as type,
         COUNT(*) as quantity
    FROM purchases p
        JOIN merch m ON p.merch_id = m.id AND m.deleted_at IS NULL
    WHERE p.user_id = $1
    GROUP BY m.name
),
received_transactions AS (
    SELECT
        sender.username as from_user,
        t.amount,
        t.created_at as date
    FROM transactions t
        JOIN users sender ON t.sender_id = sender.id AND sender.deleted_at IS NULL
    WHERE t.receiver_id = $1
    AND t.transaction_type_id = 0  -- only coin transfers
),
sent_transactions AS (
    SELECT
        receiver.username as to_user,
        t.amount,
        t.created_at as date
    FROM transactions t
        JOIN users receiver ON t.receiver_id = receiver.id AND receiver.deleted_at IS NULL
    WHERE t.sender_id = $1
    AND t.transaction_type_id = 0  -- only coin transfers
)
SELECT
    ub.id,
    ub.coins,
    COALESCE(json_agg(
        json_build_object(
            'type', ui.type,
            'quantity', ui.quantity
        )
    ) FILTER (WHERE ui.type IS NOT NULL), '[]')::jsonb as inventory,
    json_build_object(
        'received', COALESCE(
            json_agg(
                json_build_object(
                    'fromUser', rt.from_user,
                    'amount', rt.amount,
                    'date', rt.date
                )
            ) FILTER (WHERE rt.from_user IS NOT NULL),
            '[]'
        ),
        'sent', COALESCE(
            json_agg(
                json_build_object(
                    'toUser', st.to_user,
                    'amount', st.amount,
                    'date', st.date
                )
            ) FILTER (WHERE st.to_user IS NOT NULL),
            '[]'
        )
    )::jsonb as coin_history
FROM user_balance ub
LEFT JOIN user_inventory ui ON true
LEFT JOIN received_transactions rt ON true
LEFT JOIN sent_transactions st ON true
GROUP BY ub.id, ub.coins;

-- name: GetUserInfoByUsername :one
WITH user_data AS (
    SELECT id
    FROM users
    WHERE users.username = $1
      AND deleted_at IS NULL
),
user_balance AS (
    SELECT id, balance as coins
        FROM users
    WHERE id = (SELECT id FROM user_data)
        AND deleted_at IS NULL
),
user_inventory AS (
    SELECT m.name as type,
        COUNT(*) as quantity
    FROM purchases p
        JOIN merch m ON p.merch_id = m.id AND m.deleted_at IS NULL
    WHERE p.user_id = (SELECT id FROM user_data)
    GROUP BY m.name
),
received_transactions AS (
    SELECT
        sender.username as from_user,
        t.amount,
        t.created_at as date
    FROM transactions t
        JOIN users sender ON t.sender_id = sender.id AND sender.deleted_at IS NULL
    WHERE t.receiver_id = (SELECT id FROM user_data)
        AND t.transaction_type_id = 0  -- only coin transfers
),
sent_transactions AS (
    SELECT
        receiver.username as to_user,
        t.amount,
        t.created_at as date
    FROM transactions t
        JOIN users receiver ON t.receiver_id = receiver.id AND receiver.deleted_at IS NULL
    WHERE t.sender_id = (SELECT id FROM user_data)
        AND t.transaction_type_id = 0  -- only coin transfers
)
SELECT
    ub.id,
    ub.coins,
    COALESCE(json_agg(
             json_build_object(
                     'type', ui.type,
                     'quantity', ui.quantity
             )
                     ) FILTER (WHERE ui.type IS NOT NULL), '[]')::jsonb as inventory,
    json_build_object(
            'received', COALESCE(
                    json_agg(
                    json_build_object(
                            'fromUser', rt.from_user,
                            'amount', rt.amount,
                            'date', rt.date
                    )
                            ) FILTER (WHERE rt.from_user IS NOT NULL),
                    '[]'
                        ),
            'sent', COALESCE(
                            json_agg(
                            json_build_object(
                                    'toUser', st.to_user,
                                    'amount', st.amount,
                                    'date', st.date
                            )
                                    ) FILTER (WHERE st.to_user IS NOT NULL),
                            '[]'
                    )
    )::jsonb as coin_history
FROM user_balance ub
         LEFT JOIN user_inventory ui ON true
         LEFT JOIN received_transactions rt ON true
         LEFT JOIN sent_transactions st ON true
GROUP BY ub.id, ub.coins;