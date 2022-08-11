CREATE TABLE currencies
(
    id        text PRIMARY KEY,
    title_key text  NOT NULL,
    images    jsonb NOT NULL,
    sign      text  NOT NULL
)