CREATE TABLE IF NOT EXISTS promocodes (
			id SERIAL PRIMARY KEY,
			trigger character varying(100),
			description TEXT
		);

CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			user_name character varying(100),
			chat_id INTEGER
		);