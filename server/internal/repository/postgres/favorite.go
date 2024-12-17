package postgres

import (
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

func DBAddFavoriteCat(db *sqlx.DB, fav *entities.Favorite) (*entities.Favorite, error) {
	query := `INSERT INTO favorites (user_id, cat_id) VALUES (:user_id, :cat_id) RETURNING id;`

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = stmt.Get(&fav.UserID, fav)
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
