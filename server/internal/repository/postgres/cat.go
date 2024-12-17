package postgres

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"server/internal/entities"
)

func DBCatCreate(db *sqlx.DB, cat *entities.Cat) (*entities.Cat, error) {
	query := `
		INSERT INTO cats (breed, fur, temper, care_complexity, image_path)
		VALUES (:breed, :fur, :temper, :care_complexity, :image_path)
	`

	stmt, err := db.PrepareNamed(query)
	if stmt == nil {
		return nil, err
	}
	err = stmt.Get(&cat.ID, cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func DBCatExistsID(db *sqlx.DB, catID int) (bool, error) {
	exists := 0
	query := `SELECT 1 FROM cats WHERE id = $1 LIMIT 1`

	err := db.QueryRow(query, catID).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func DBCatExistsBreed(db *sqlx.DB, breed string) (bool, error) {
	exists := 0
	query := `SELECT 1 FROM cats WHERE breed = $1 LIMIT 1`

	err := db.QueryRow(query, breed).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}
