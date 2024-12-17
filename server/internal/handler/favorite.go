package handler

import (
	"github.com/gofiber/fiber/v2"
	"server/internal/entities"
	"server/internal/log"
	"server/internal/repository/postgres"
)

// GetFavoriteCats
// @Tags         cat
// @Summary      Получение списка любимых кошек пользователя
// @Description  Получение списка кошек, которые являются любимыми у пользователя по его идентификатору с логированием ошибок
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {array}   entities.Cat "Успешное получение списка любимых кошек"
// @Failure      400  {object}  entities.ErrorResponse "Некорректный идентификатор пользователя"
// @Failure      500  {object}  entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/favorites [get]
// @Security ApiKeyAuth
func (h *Handler) GetFavoriteCats(c *fiber.Ctx) error {
	id, ok := c.Locals("id").(int)
	if !ok {
		return c.SendStatus(fiber.StatusForbidden)
	}

	h.logger.Debug().Msg("call postgres.DBFormulaExists")
	cats, err := postgres.DBGetFavoriteCats(h.db, int(id))
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(cats)
}

// AddFavoriteCat
// @Tags         cat
// @Summary      Добавление кошки в список любимых
// @Description  Добавление кошки в список любимых пользователя по его идентификатору и идентификатору кошки с логированием ошибок
// @Accept       json
// @Produce      json
// @Param        body  body      entities.UpdateFavoriteRequest  true  "Информация о кошке и пользователе"
// @Success      200   {object}  entities.Favorite "Успешное добавление кошки в список любимых"
// @Failure      400   {object}  entities.ErrorResponse "Некорректные данные запроса"
// @Failure      500   {object}  entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/favorites [post]
// @Security ApiKeyAuth
func (h *Handler) AddFavoriteCat(c *fiber.Ctx) error {
	var req entities.UpdateFavoriteRequest
	if err := c.BodyParser(&req); err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	favorite := &entities.Favorite{
		UserID: req.UserID,
		CatID:  req.CatID,
	}
	h.logger.Debug().Msg("call postgres.AddFavoriteCat")
	res, err := postgres.DBAddFavoriteCat(h.db, favorite)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(res)
}

// DeleteFavoriteCat
// @Tags         cat
// @Summary      Удаление кошки из списка любимых
// @Description  Удаление кошки из списка любимых пользователя по его идентификатору и идентификатору кошки с логированием ошибок
// @Accept       json
// @Produce      json
// @Param        body  body      entities.UpdateFavoriteRequest  true  "Информация о кошке и пользователе"
// @Success      200   {object}  map[string]string "Успешное удаление кошки из списка любимых"
// @Failure      400   {object}  entities.ErrorResponse "Некорректные данные запроса"
// @Failure      500   {object}  entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /auth/favorites [delete]
// @Security ApiKeyAuth
func (h *Handler) DeleteFavoriteCat(c *fiber.Ctx) error {
	var req entities.UpdateFavoriteRequest
	if err := c.BodyParser(&req); err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	favorite := &entities.Favorite{
		UserID: req.UserID,
		CatID:  req.CatID,
	}
	h.logger.Debug().Msg("call postgres.DBRemoveFavoriteCat")
	err := postgres.DBRemoveFavoriteCat(h.db, favorite)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}