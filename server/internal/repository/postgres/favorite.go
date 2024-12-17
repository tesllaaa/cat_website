package postgres

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"server/internal/entities"
)

func DBGetFavoriteCats(db *sqlx.DB, userID int) (*[]entities.FavoriteCat, error) {
	var favorites []entities.FavoriteCat

	query := `
	SELECT cats.id, cats.breed, cats.image_path FROM favorites
	JOIN cats ON favorites.cat_id = cats.id
	WHERE favorites.user_id = $1;`

	err := db.Select(&favorites, query, userID)
	if err != nil {
		return nil, err
	}
	return &favorites, nil

}

func DBFavoriteExists(db *sqlx.DB, favorite *entities.Favorite) (bool, error) {
	exists := 0
	query := `SELECT 1 FROM favorites WHERE user_id = $1 AND cat_id = $2 LIMIT 1`

	err := db.QueryRow(query, favorite.UserID, favorite.CatID).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func DBAddFavoriteCat(db *sqlx.DB, fav *entities.Favorite) (*entities.Favorite, error) {
	query := `INSERT INTO favorites (user_id, cat_id) VALUES ($1, $2) RETURNING id;`

	err := db.QueryRow(query, fav.UserID, fav.CatID).Scan(&fav.ID)
	if err != nil {
		return nil, err
	}
	return fav, nil
}

func DBRemoveFavoriteCat(db *sqlx.DB, fav *entities.Favorite) error {
	query := `DELETE FROM favorites WHERE user_id = $1 and cat_id = $2;`
	_, err := db.Exec(query, fav.UserID, fav.CatID)
	if err != nil {
		return err
	}
	return nil
}
