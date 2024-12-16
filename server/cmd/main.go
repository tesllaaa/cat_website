package main

import (
	"fmt"
	_ "server/docs"
	"server/internal/handler"
	logger "server/internal/log"
	"server/internal/repository/postgres"
	//"server/util"
)

// @title Kotiki API
// @version 1.0
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Инициализация логера
	log := logger.InitLogger()
	// Инициализация бд
	db, err := postgres.NewDatabase()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("could not initialize database connection: %s", err))
	}
	// Создание директорий для временных файлов
	//util.CreateDirectory()
	// Инициализация ручек
	handlers := handler.NewHandler(db, log)

	// Запуск сервера
	app := handlers.Router()
	app.Listen(":8080")
}
