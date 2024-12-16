package entities

// User базовая структура пользователя
type User struct {
	ID        int    `json:"id" db:"id"`
	Password  string `json:"password" db:"password"`
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	ThirdName string `json:"third_name" db:"third_name"`
}

// UserData базовая структура пользовательских данных
type UserData struct {
	ID        int    `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Name      string `json:"name" db:"name"`
	Surname   string `json:"surname" db:"surname"`
	ThirdName string `json:"third_name" db:"third_name"`
}

// CreateUserRequest структура запроса на создание пользователя
type CreateUserRequest struct {
	Password string `json:"password" example:"12345678"`
	Email    string `json:"email" example:"petrov@mail.ru"`
	FullName string `json:"full_name" example:"Петров Петр Иванович"`
}

// CreateUserResponse структура ответа на создание пользователя
type CreateUserResponse struct {
	AccessToken string `json:"access_token"`
	ID          int    `json:"id"`
}

// LoginUserRequest структура запроса на вход
type LoginUserRequest struct {
	Email    string `json:"email" example:"petrov@mail.ru"`
	Password string `json:"password" example:"12345678"`
}

// LoginUserResponse структура ответа на вход
type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	ID          int    `json:"id"`
}