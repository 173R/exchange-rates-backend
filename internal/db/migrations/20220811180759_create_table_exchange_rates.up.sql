CREATE TABLE exchange_rates
(
    id               bigserial PRIMARY KEY,
    base_currency_id text REFERENCES currencies (id) NOT NULL,
    timestamp        timestamp                       NOT NULL,
    rates            jsonb                           NOT NULL
)