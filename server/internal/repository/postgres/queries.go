package postgres

// запросы создания таблиц
const (
	createUserTable = `
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR NOT NULL,
		surname VARCHAR NOT NULL,
		email VARCHAR NOT NULL UNIQUE,
		password VARCHAR NOT NULL
);
`
)
