package entities

// Message Структура для успешных ответов
type Message struct {
	Message string `json:"message"`
}

// ErrorResponse Структура для ответов с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// Id Структура для ответов с id
type Id struct {
	Id int `json:"id"`
}
