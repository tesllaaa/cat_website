package entities

type Favorite struct {
	ID     int `json:"id" db:"id" example:"5"`
	UserID int `json:"user_id" db:"user_id" example:"44"`
	CatID  int `json:"cat_id" db:"cat_id" example:"3"`
}

type CreateFavoriteRequest struct {
	UserID int `json:"user_id" db:"user_id" example:"5"`
	CatID  int `json:"cat_id" db:"cat_id" example:"3"`
}
