CREATE TABLE refresh_sessions
(
    id            bigserial PRIMARY KEY,
    user_id       bigint REFERENCES users (id) NOT NULL,
    refresh_token text                         NOT NULL,
    access_token  text                         NOT NULL,
    fingerprint   text                         NOT NULL,
    expires_at    timestamp                    NOT NULL,
    created_at    timestamp                    NOT NULL
)