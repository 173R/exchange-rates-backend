CREATE TABLE user_observed_currencies
(
    id          bigserial PRIMARY KEY,
    user_id     bigint REFERENCES users (id)    NOT NULL,
    currency_id text REFERENCES currencies (id) NOT NULL,
    UNIQUE (user_id, currency_id)
)