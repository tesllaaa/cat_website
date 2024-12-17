package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"server/internal/config"
)

// NewDatabase инициализация подключения к бд
func NewDatabase() (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database connection: %v", err)
	}

	fmt.Println("Successfully connected to database. Creating tables")
	CreateTable(db)

	return db, nil
}

// CreateTable cоздание всех таблиц
func CreateTable(db *sqlx.DB) {
	db.MustExec(createUserTable)
	db.MustExec(createCatTable)
	db.MustExec(createFavoritesTable)
}
