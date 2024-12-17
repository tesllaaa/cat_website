package entities

type Cat struct {
	ID             int    `json:"id" db:"id" example:"7"`
	Breed          string `json:"breed" db:"breed" example:"Мейн-кун"`
	Fur            string `json:"fur" db:"fur" example:"Длинношерстная"`
	Temper         string `json:"temper" db:"temper" example:"Спокойный"`
	CareComplexity int    `json:"care_complexity" db:"care_complexity" example:"4"`
	ImagePath      string `json:"image_path" db:"image_path" example:"/images/cat.png"`
}

type CreateCatRequest struct {
	Breed          string `form:"breed"`
	Fur            string `form:"fur"`
	Temper         string `form:"temper"`
	CareComplexity int    `form:"care_complexity"`
	Image          string `form:"image"`
}

type UpdateCatRequest struct {
	ID             int    `json:"id" db:"id" example:"7"`
	Breed          string `json:"breed" db:"breed" example:"Мейн-кун"`
	Fur            string `json:"fur" db:"fur" example:"Длинношерстная"`
	Temper         string `json:"temper" db:"temper" example:"Спокойный"`
	CareComplexity int    `json:"care_complexity" db:"care_complexity" example:"4"`
}

type FavoriteCat struct {
	Breed     string `json:"breed" db:"breed" example:"Мейн-кун"`
	ID        int    `json:"id" db:"id" example:"7"`
	ImagePath string `json:"image_path" db:"image_path" example:"/images/cat.png"`
}
