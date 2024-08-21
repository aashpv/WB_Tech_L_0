CREATE DATABASE wb_order_service;

\c wb_order_service;

CREATE TABLE IF NOT EXISTS public.orders
(
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    locale VARCHAR(255) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id BIGINT NOT NULL,
    oof_shard VARCHAR(255) NOT NULL,
    date_created TIMESTAMP NOT NULL
);

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;

ALTER TABLE IF EXISTS public.delivery
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.payment
(
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(255) NOT NULL,
    provider VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee INTEGER NOT NULL,
    CONSTRAINT payment_fkey FOREIGN KEY (order_id)
        REFERENCES public.orders (id) ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.payment
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.items
(
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(255) NOT NULL,
    total_price INTEGER NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    status INTEGER NOT NULL,
    CONSTRAINT items_fkey FOREIGN KEY (order_id)
        REFERENCES public.orders (id) ON DELETE CASCADE
);

ALTER TABLE IF EXISTS public.items
    OWNER to postgres;

CREATE TABLE IF NOT EXISTS public.delivery
(
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    zip VARCHAR(255) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT delivery_fkey FOREIGN KEY (order_id)
        REFERENCES public.orders (id) ON DELETE CASCADE
);
