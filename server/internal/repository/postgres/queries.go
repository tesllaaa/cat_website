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
	createCatTable = `
		CREATE TABLE IF NOT EXISTS cats (
		id SERIAL PRIMARY KEY,
		breed VARCHAR NOT NULL,
		furname VARCHAR NOT NULL,
		temper VARCHAR NOT NULL,
		care_complexity INTEGER NOT NULL,
		image_path VARCHAR NOT NULL
);
`

	createFavoritesTable = `
		CREATE TABLE IF NOT EXISTS favorites (
		    id SERIAL PRIMARY KEY,
		    user_id INTEGER references users(id) ON DELETE CASCADE,
		    cat_id INTEGER references cats(id) ON DELETE CASCADE
);
`
)
