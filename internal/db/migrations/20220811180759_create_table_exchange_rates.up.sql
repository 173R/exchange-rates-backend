CREATE TABLE exchange_rates
(
    id           bigserial PRIMARY KEY,
    timestamp    timestamp                       NOT NULL,
    currency_id  text REFERENCES currencies (id) NOT NULL,
    convert_rate numeric                         NOT NULL,
    UNIQUE (timestamp, currency_id)
)