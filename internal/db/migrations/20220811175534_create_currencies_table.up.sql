CREATE TABLE currencies
(
    id     text PRIMARY KEY,
    title  jsonb NOT NULL,
    images jsonb,
    sign   text  NOT NULL
)