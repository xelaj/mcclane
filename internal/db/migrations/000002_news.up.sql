BEGIN;

CREATE TABLE IF NOT EXISTS news (
id SERIAL NOT NULL PRIMARY KEY,
text VARCHAR ,
hot_location_id INT REFERENCES hot_location(id) ON DELETE CASCADE ,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;