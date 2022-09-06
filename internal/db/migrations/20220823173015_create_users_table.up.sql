CREATE TABLE users
(
    id               bigserial PRIMARY KEY,
    telegram_uid     bigint                          NOT NULL,
    base_currency_id text REFERENCES currencies (id) NOT NULL,
    lang             text                            NOT NULL
)