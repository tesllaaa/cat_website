package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"server/internal/config"
	"server/internal/log"
	"server/pkg"

	//"server/pkg"

	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "server/docs"
)

// Handler Инициализация структуры ручки
type Handler struct {
	db     *sqlx.DB
	logger *zerolog.Logger
}

// NewHandler Инициализация экземпляра ручки
func NewHandler(db *sqlx.DB, logger *zerolog.Logger) *Handler {
	return &Handler{db: db, logger: logger}
}

// Router Инициализация всех запросов
func (h *Handler) Router() *fiber.App {
	f := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
	})

	// CORS middleware
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		//AllowCredentials: true,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	}))
	f.Use(log.RequestLogger(h.logger)) // Logger middleware

	f.Get("/swagger/*", fiberSwagger.WrapHandler)
	f.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("healthy")
	})

	f.Post("/signup", h.SignUp)
	f.Post("/login", h.Login)
	f.Get("/user/:id", h.GetUserDataByID)

	f.Post("/cat", h.CatCreate)
	f.Put("/cat", h.CatUpdate)
	f.Delete("/cat/id/:id", h.CatDelete)
	f.Get("/cat/id/:id", h.CatGetByID)
	f.Get("/cat", h.CatGetAll)

	// Ручки доступные после авторизации пользователя
	authGroup := f.Group("/auth")
	authGroup.Use(func(c *fiber.Ctx) error {
		return pkg.WithJWTAuth(c, config.SigningKey)
	})

	return f
}
