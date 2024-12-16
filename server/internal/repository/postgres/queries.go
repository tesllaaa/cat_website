package postgres

// запросы создания таблиц
const (
	createUserTable = `
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR NOT NULL,
		surname VARCHAR NOT NULL,
		third_name VARCHAR,
		email VARCHAR NOT NULL UNIQUE,
		password VARCHAR NOT NULL
);
`

	createArticleTable = `
		CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		title VARCHAR NOT NULL,
		science VARCHAR NOT NULL,
		section VARCHAR NOT NULL,
		path VARCHAR NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE 
);
`
	createTableFormula = `
		CREATE TABLE IF NOT EXISTS formula (
		id SERIAL PRIMARY KEY,
		title VARCHAR NOT NULL,
		value VARCHAR NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
	`

	createFormulaVCSTable = `
		CREATE TABLE IF NOT EXISTS formula_vcs (
		id SERIAL PRIMARY KEY,
		formula_id INT NOT NULL,
		difference VARCHAR NOT NULL,
		hash VARCHAR NOT NULL,
		code_name VARCHAR NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		FOREIGN KEY (formula_id) REFERENCES formula (id) ON DELETE CASCADE
);
	`
)
