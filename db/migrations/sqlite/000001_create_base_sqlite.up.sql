CREATE TABLE IF NOT EXISTS airports (
	id INTEGER PRIMARY KEY NOT NULL,
	city TEXT NOT NULL,
	iata_code TEXT NOT NULL,
	country TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS tickets (
	id TEXT PRIMARY KEY NOT NULL,
	airline TEXT NOT NULL,
	departure_from INTEGER NOT NULL,
	arrival_at INTEGER NOT NULL,
	departure_time TEXT NOT NULL,
	arrival_time TEXT NOT NULL,
	quantity INTEGER NOT NULL,
	luggage TEXT,
	hand_baggage TEXT,
	price REAL NOT NULL,
	FOREIGN KEY (departure_from) REFERENCES airports(id)
	ON UPDATE RESTRICT ON DELETE RESTRICT,
	FOREIGN KEY (arrival_at) REFERENCES airports(id)
	ON UPDATE RESTRICT ON DELETE RESTRICT);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS buy_history (
	id INTEGER PRIMARY KEY NOT NULL,
	ticket_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	buy_time TEXT NOT NULL,
	FOREIGN KEY (ticket_id) REFERENCES tickets(id)
	ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id)
	ON UPDATE CASCADE ON DELETE CASCADE);



CREATE TABLE IF NOT EXISTS passport_data(
	id INTEGER PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	surname TEXT NOT NULL,
	patronymic TEXT,
	passport_series_and_number TEXT NOT NULL,
	gender TEXT NOT NULL CHECK(gender IN (w, m)),
	validity_period TEXT,
	date_of_birth TEXT NOT NULL,
	passport_type TEXT NOT NULL CHECK(passport_type IN (international passport, passport)),
	citizenship TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
	ON UPDATE CASCADE ON DELETE CASCADE);
