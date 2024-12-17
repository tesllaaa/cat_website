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
	ID      int    `json:"id" db:"id" example:"1"`
	Email   string `json:"email" db:"email" example:"petrov@mail.ru"`
	Name    string `json:"name" db:"name" example:"Петр"`
	Surname string `json:"surname" db:"surname" example:"Петров"`
}

// CreateUserRequest структура запроса на создание пользователя
type CreateUserRequest struct {
	Password string `json:"password" example:"12345678"`
	Email    string `json:"email" example:"petrov@mail.ru"`
	Name     string `json:"name" db:"name" example:"Петр"`
	Surname  string `json:"surname" db:"surname" example:"Петров"`
}

// CreateUserResponse структура ответа на создание пользователя
type CreateUserResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiRXhwaXJlc0F0IjoxNzM4MDE3MjA5LCJJc3N1ZWRBciI6MTczNDQxNzIwOX0sInVzZXJfaWQiOjF9.CX2eHEjXZa209vDdtdoz40JlbxuHukMjrc-hw2E_Jy0"`
	ID          int    `json:"id" example:"1"`
}

// LoginUserRequest структура запроса на вход
type LoginUserRequest struct {
	Email    string `json:"email" example:"petrov@mail.ru"`
	Password string `json:"password" example:"12345678"`
}

// LoginUserResponse структура ответа на вход
type LoginUserResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiRXhwaXJlc0F0IjoxNzM4MDE3MjA5LCJJc3N1ZWRBciI6MTczNDQxNzIwOX0sInVzZXJfaWQiOjF9.CX2eHEjXZa209vDdtdoz40JlbxuHukMjrc-hw2E_Jy0"`
	ID          int    `json:"id" example:"1"`
}
