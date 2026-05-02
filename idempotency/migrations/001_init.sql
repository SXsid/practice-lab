--+goose up

CREATE TYPE payment_status AS ENUM ('pending', 'succeeded', 'failed');

CREATE TABLE customers(
    id          BIGSERIAL       PRIMARY KEY ,
    name        VARCHAR(256) NOT NULL
);



CREATE TABLE payments (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id            TEXT            NOT NULL UNIQUE,
    amount              BIGINT          NOT NULL CHECK (amount > 0),
    currency            CHAR(3)         NOT NULL,
    status              payment_status  NOT NULL DEFAULT 'pending',
    customer_id         BIGSERIAL  REFERENCES customers(id) ,
    provider_charge_id  TEXT,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()

);

CREATE INDEX idx_payments_order_id ON payments (order_id);
CREATE INDEX idx_payments_status   ON payments (status);


