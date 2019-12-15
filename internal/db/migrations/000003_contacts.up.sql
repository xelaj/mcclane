BEGIN;

CREATE TABLE IF NOT EXISTS contacts (
id SERIAL NOT NULL PRIMARY KEY,
contact VARCHAR ,
user_chat_id INT REFERENCES user_chat(id) ON DELETE CASCADE ,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;