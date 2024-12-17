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
		VALUES (:breed, :fur, :temper, :care_complexity, :image_path) RETURNING id
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

func DBCatUpdate(db *sqlx.DB, cat *entities.UpdateCatRequest) error {
	query := `UPDATE cats SET breed = $1, fur = $2, temper = $3, care_complexity = $4 WHERE id = $5`
	_, err := db.Exec(query, cat.Breed, cat.Fur, cat.Temper, cat.CareComplexity, cat.ID)
	if err != nil {
		return err
	}
	return nil
}

func DBCatDelete(db *sqlx.DB, cat *entities.Cat) error {
	query := `DELETE FROM cats WHERE id = $1`
	_, err := db.Exec(query, cat.ID)
	if err != nil {
		return err
	}
	return nil
}

func DBCatGetByID(db *sqlx.DB, catID int) (*entities.Cat, error) {
	cat := entities.Cat{}
	query := `select * from cats where id = $1
	`

	err := db.Get(&cat, query, catID)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func DBCatGetAll(db *sqlx.DB) (*[]entities.Cat, error) {
	var cats []entities.Cat
	query := `SELECT * FROM cats ORDER BY id`
	err := db.Select(&cats, query)
	if err != nil {
		return nil, err
	}
	return &cats, nil
}
