CREATE TABLE exchange_rates
(
    id        bigserial PRIMARY KEY,
    timestamp timestamp NOT NULL,
    rates     jsonb     NOT NULL,
    UNIQUE (timestamp)
)