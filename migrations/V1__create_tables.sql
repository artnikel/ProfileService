CREATE TABLE IF NOT EXISTS users (
	id uuid,
	login VARCHAR,
	password VARCHAR,
	refreshToken VARCHAR,
	primary key (id)
);