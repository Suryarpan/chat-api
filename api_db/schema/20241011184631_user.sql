-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID,
    name VARCHAR(100),
    email VARCHAR(255),
    password TEXT,
    password_salt TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
