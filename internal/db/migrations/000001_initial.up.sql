BEGIN;

CREATE TABLE IF NOT EXISTS user_chat (
id SERIAL NOT NULL PRIMARY KEY,
user_name VARCHAR ,
chat_id INT ,
warning BOOL,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS hot_location (
id SERIAL NOT NULL PRIMARY KEY,
name VARCHAR ,
point_x_lat FLOAT ,
point_y_lon FLOAT ,
point_y_lat FLOAT ,
point_x_lon FLOAT ,
event_date TIMESTAMP ,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coordinates (
id SERIAL NOT NULL PRIMARY KEY,
user_chat_id INT REFERENCES user_chat(id) ON DELETE CASCADE ,
hot_location_id INT REFERENCES hot_location(id) ON DELETE CASCADE ,
latitude FLOAT ,
longitude FLOAT ,
created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;