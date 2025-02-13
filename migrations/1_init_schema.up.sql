CREATE TABLE IF NOT EXISTS users
(
    id            CHARACTER VARYING PRIMARY KEY,
    username      CHARACTER VARYING UNIQUE NOT NULL,
    password_hash CHARACTER VARYING NOT NULL,
    balance       INT NOT NULL DEFAULT 1000 CHECK (balance >= 0),
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_users ON users (username) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS transaction_types
(
    id    INT PRIMARY KEY,
    title CHARACTER VARYING NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS transactions
(
    id                  CHARACTER VARYING PRIMARY KEY,
    sender_id           CHARACTER VARYING NOT NULL,
    receiver_id         CHARACTER VARYING DEFAULT NULL,
    transaction_type_id INT NOT NULL,
    amount              INT NOT NULL CHECK (amount > 0),
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_transactions ON transactions (sender_id, receiver_id, transaction_type_id, amount);

CREATE TABLE IF NOT EXISTS merch
(
    id         CHARACTER VARYING PRIMARY KEY,
    name       CHARACTER VARYING UNIQUE NOT NULL,
    price      INT NOT NULL CHECK (price > 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_merch ON merch (name) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS purchases
(
    id         CHARACTER VARYING PRIMARY KEY,
    user_id    CHARACTER VARYING NOT NULL,
    merch_id   CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_active_purchases ON purchases (user_id, merch_id);

ALTER TABLE transactions ADD FOREIGN KEY (sender_id) REFERENCES users(id);
ALTER TABLE transactions ADD FOREIGN KEY (receiver_id) REFERENCES users(id);
ALTER TABLE transactions ADD FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id);
ALTER TABLE purchases ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE purchases ADD FOREIGN KEY (merch_id) REFERENCES merch(id);