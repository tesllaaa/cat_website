package handler

import (
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"server/internal/entities"
	"server/internal/log"
	"server/internal/repository/postgres"
	"strconv"
)

// CatCreate
// @Tags         cats
// @Summary      Создание записи о кошке
// @Description  Добавление новой записи о кошке в базу данных с логированием ошибок
// @Accept       multipart/form-data
// @Produce      json
// @Param        fur            formData string true "Шерсть кошки"
// @Param        breed          formData string true "Порода кошки"
// @Param        care_complexity formData integer true "Сложность ухода за кошкой"
// @Param        temper         formData string true "Темперамент кошки"
// @Param        image          formData file true "Изображение кошки"
// @Success      200 {object} entities.Cat "Успешное создание записи"
// @Failure      400 {object} entities.ErrorResponse "Некорректные данные"
// @Failure      500 {object} entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /cat [post]
func (h *Handler) CatCreate(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("failed to retrieve file")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to retrieve file"})
	}

	if file.Header.Get("Content-Type") != "image/jpeg" {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("only JPEG images are allowed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "only JPEG images are allowed"})
	}

	path := "/.tmp"
	savePath := filepath.Join(path, file.Filename)

	if err := c.SaveFile(file, savePath); err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Err(err).Msg("failed to save file")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to save file"})
	}

	var cat entities.Cat

	careComp, err := strconv.Atoi(c.FormValue("care_complexity"))
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cat.Fur = c.FormValue("fur")
	cat.Breed = c.FormValue("breed")
	cat.CareComplexity = careComp
	cat.Temper = c.FormValue("temper")
	cat.ImagePath = savePath

	h.logger.Debug().Msg("call postgres.DBCatExistsBreed")
	exists, err := postgres.DBCatExistsBreed(h.db, cat.Breed)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if exists {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg("cat already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "article already exists"})
	}

	h.logger.Debug().Msg("call postgres.DBCatCreate")
	res, err := postgres.DBCatCreate(h.db, &cat)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(res)
}
