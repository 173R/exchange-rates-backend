CREATE TABLE currencies
(
    id                   text PRIMARY KEY,
    title_translation_id text REFERENCES translations (id) NOT NULL,
    images               jsonb                             NOT NULL,
    sign                 text                              NOT NULL
)