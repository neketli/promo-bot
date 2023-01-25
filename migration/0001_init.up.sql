CREATE TABLE IF NOT EXISTS promocodes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			trigger TEXT,
			description TEXT
		)

CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT,
			password TEXT,
			chat_id INTEGER
		)