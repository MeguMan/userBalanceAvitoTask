CREATE TABLE users (
    user_id bigserial PRIMARY KEY,
    balance int
);

INSERT INTO users (user_id) VALUES (0);
INSERT INTO users (user_id) VALUES (1);