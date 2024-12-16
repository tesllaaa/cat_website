package postgres

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"server/internal/entities"
)

// DBUserGetById получение пользователя по айди
func DBUserGetById(db *sqlx.DB, id int64) (*entities.User, error) {
	user := entities.User{}
	query := `SELECT id, name, surname, third_name, email, password FROM users WHERE id = $1`
	err := db.Get(&user, query, id)
	if err != nil {
		return &entities.User{}, nil
	}

	return &user, nil
}

// DBUserDataGetById получение данных пользователя по его айди
func DBUserDataGetById(db *sqlx.DB, id int64) (*entities.UserData, error) {
	user := entities.UserData{}
	query := `SELECT id, name, surname, third_name, email FROM users WHERE id = $1`
	err := db.Get(&user, query, id)
	if err != nil {
		return &entities.UserData{}, nil
	}

	return &user, nil
}

// DBUserGetByEmail получение пользователя по email
func DBUserGetByEmail(db *sqlx.DB, email string) (*entities.User, error) {
	user := entities.User{}
	query := `SELECT id, name, surname, third_name, email, password FROM users WHERE email = $1`
	err := db.Get(&user, query, email)
	if err != nil {
		return &entities.User{}, nil
	}

	return &user, nil

}

// DBUserExists проверка существования пользователя в бд (по почте)
func DBUserExists(db *sqlx.DB, email string) (bool, error) {
	exists := 0
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`

	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

// DBUserExistsID проверка существования пользователя в бд (по айди)
func DBUserExistsID(db *sqlx.DB, id int64) (bool, error) {
	exists := 0
	query := `SELECT 1 FROM users WHERE id = $1 LIMIT 1`

	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

// DBUserCreate создание пользователя
func DBUserCreate(db *sqlx.DB, user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (email, password, name, surname, third_name)
	VALUES (:email, :password, :name, :surname, :third_name) RETURNING id`

	stmt, err := db.PrepareNamed(query)
	if stmt == nil {
		return nil, errors.New("error preparing statement")
	}
	err = stmt.Get(&user.ID, *user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
