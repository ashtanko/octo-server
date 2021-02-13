CREATE TABLE IF NOT EXISTS account
(
    id              bigint  NOT NULL PRIMARY KEY DEFAULT 0,
    account_balance NUMERIC NOT NULL             DEFAULT 0.0 CHECK (account_balance >= 0.0)
);

INSERT INTO account
VALUES (0, 0);

CREATE TABLE IF NOT EXISTS transaction_type
(
    t_type TEXT PRIMARY KEY
);

INSERT INTO transaction_type
VALUES ('game'),
       ('server'),
       ('payment')
;

CREATE TABLE IF NOT EXISTS transaction_status
(
    t_status TEXT PRIMARY KEY
);

INSERT INTO transaction_status
VALUES ('DONE'),
       ('CANCELLED'),
       ('DELETED');


CREATE TABLE IF NOT EXISTS transactions
(
    id            SERIAL,
    type_key      TEXT    NOT NULL REFERENCES transaction_type (t_type),
    status_key    TEXT    NOT NULL REFERENCES transaction_status (t_status),
    amount        NUMERIC(15, 2) DEFAULT 0.0,
    timestamptz   TIMESTAMP      DEFAULT CURRENT_TIMESTAMP,
    transactionId VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS jobs
(
    id               SERIAL PRIMARY KEY,
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO jobs
VALUES (0, CURRENT_TIMESTAMP)
ON CONFLICT DO NOTHING;