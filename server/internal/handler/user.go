package handler

import (
	"github.com/gofiber/fiber/v2"
	"server/internal/config"
	"server/internal/entities"
	"server/internal/log"
	"server/internal/repository/postgres"
	"server/pkg"
	"server/util"
	"strconv"
)

// SignUp
// @Tags         user
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя и выдает ему токен доступа
// @Accept       json
// @Produce      json
// @Param        data body entities.CreateUserRequest true "Данные для регистрации"
// @Success      200 {object} entities.CreateUserResponse "Регистрация успешна"
// @Failure      400 {object} entities.ErrorResponse "Пользователь уже существует или данные неверны"
// @Failure      500 {object} entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /signup [post]
func (h *Handler) SignUp(c *fiber.Ctx) error {
	var u entities.CreateUserRequest
	if err := c.BodyParser(&u); err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Debug().Msg("call postgres.DBUserExists")
	exists, err := postgres.DBUserExists(h.db, u.Email)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if exists {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("user already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user already exists"})
	}

	hashedPassword, err := util.HashPassword(u.Password)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user := &entities.User{
		Email:    u.Email,
		Password: hashedPassword,
		Name:     u.Name,
		Surname:  u.Surname,
	}

	h.logger.Debug().Msg("call postgres.DBUserCreate")
	r, err := postgres.DBUserCreate(h.db, user)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	TokenExpiration, err := strconv.Atoi(config.TokenExpiration)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg("wrong data")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "wrong data"})
	}
	h.logger.Debug().Msg("call pkg.GenerateAccessToken")
	accessToken, err := pkg.GenerateAccessToken(user.ID, TokenExpiration,
		config.SigningKey)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Debug().Msg("call pkg.GenerateRefreshToken")

	res := &entities.CreateUserResponse{
		ID:          r.ID,
		AccessToken: accessToken,
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(res)

}

// Login
// @Tags         user
// @Summary      Вход пользователя
// @Description  Аутентификация пользователя с возвращением токена доступа
// @Accept       json
// @Produce      json
// @Param        data body entities.LoginUserRequest true "Данные для входа"
// @Success      200 {object} entities.LoginUserResponse "Успешный вход"
// @Failure      400 {object} entities.ErrorResponse "Неверный логин или пароль"
// @Failure      500 {object} entities.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var user entities.LoginUserRequest
	if err := c.BodyParser(&user); err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	h.logger.Debug().Msg("call postgres.DBUserGetByLogin")
	u, err := postgres.DBUserGetByEmail(h.db, user.Email)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("wrong data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wrong data"})
	}
	h.logger.Debug().Msg("call util.CheckPassword")
	err = util.CheckPassword(user.Password, u.Password)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("wrong data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wrong data"})
	}

	TokenExpiration, err := strconv.Atoi(config.TokenExpiration)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("wrong data")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wrong data"})
	}
	h.logger.Debug().Msg("call pkg.GenerateAccessToken")
	accessToken, err := pkg.GenerateAccessToken(u.ID, TokenExpiration,
		config.SigningKey)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Debug().Msg("call pkg.GenerateRefreshToken")

	res := entities.LoginUserResponse{
		AccessToken: accessToken,
		ID:          u.ID,
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(res)
}

// GetUserDataByID
// @Tags         user
// @Summary      Получение данных пользователя по его ID
// @Description  Извлекает пользовательские данные из БД по ID пользователя
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  entities.UserData  "Пользовательские данные получены"
// @Failure      400  {object}  entities.ErrorResponse  "Неверный формат ID"
// @Failure      500  {object}  entities.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /user/{id} [get]
func (h *Handler) GetUserDataByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	h.logger.Debug().Msg("call postgres.DBUserExistsID")
	exists, err := postgres.DBUserExistsID(h.db, int64(id))
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if !exists {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("user not exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user not exists"})
	}

	h.logger.Debug().Msg("call postgres.DBUserGetByID")
	user, err := postgres.DBUserDataGetById(h.db, int64(id))
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(user)
}
