-- +goose Up
-- +goose StatementBegin
CREATE TYPE message_status AS ENUM ('sent', 'delivered', 'read');

CREATE TABLE message_meta (
    mssg_id BIGSERIAL PRIMARY KEY,
    from_pvt_id INTEGER NOT NULL REFERENCES users
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    to_pvt_id INTEGER NOT NULL REFERENCES users
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    mssg_status message_status NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE message_text (
    mssg_id BIGINT PRIMARY KEY REFERENCES message_meta
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    mssg_body TEXT NOT NULL
) PARTITION BY hash(mssg_id);

CREATE TABLE message_text_0 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 0);
CREATE TABLE message_text_1 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 1);
CREATE TABLE message_text_2 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 2);
CREATE TABLE message_text_3 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 3);
CREATE TABLE message_text_4 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 4);
CREATE TABLE message_text_5 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 5);
CREATE TABLE message_text_6 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 6);
CREATE TABLE message_text_7 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 7);
CREATE TABLE message_text_8 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 8);
CREATE TABLE message_text_9 PARTITION OF message_text FOR VALUES WITH (MODULUS 10,REMAINDER 9);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE message_text;
DROP TABLE message_meta;
DROP TYPE message_status;
-- +goose StatementEnd
