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
	"strings"
)

// SignUp
// @Tags         user
// @Summary      Регистрация пользователя
// @Description  Создает нового пользователя и выдает ему токен доступа
// @Accept       json
// @Produce      json
// @Param        data body entities.CreateUserRequest true "User  Данные для регистрации"
// @Success      200 {object} entities.CreateUserResponse "User  Регистрация успешна"
// @Failure      400 {object} entities.ErrorResponse "User  Пользователь уже существует или данные неверны"
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

	fullname := strings.Split(u.FullName, " ")
	if len(fullname) < 3 {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	user := &entities.User{
		Email:     u.Email,
		Password:  hashedPassword,
		Name:      fullname[1],
		Surname:   fullname[0],
		ThirdName: fullname[2],
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
// @Summary      User login
// @Description  Authenticates a user and returns access and refresh tokens.
// @Accept       json
// @Produce      json
// @Param        data body entities.LoginUserRequest true "User  login credentials"
// @Success      200 {object} entities.LoginUserResponse "User  successfully logged in"
// @Failure      400 {object} entities.ErrorResponse "Invalid email or password"
// @Failure      500 {object} entities.ErrorResponse "Internal server error"
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
// @Summary      Retrieve user data by ID
// @Description  Fetches user details from the database using the provided user ID.
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  entities.UserData  "User data retrieved successfully"
// @Failure      400  {object}  entities.ErrorResponse  "Invalid user ID format"
// @Failure      500  {object}  entities.ErrorResponse  "Internal server error"
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

// CheckAuth
// @Tags         user
// @Summary      Authorization check
// @Description  Validates the JWT token from the Authorization header, extracts user ID, and generates a new access token.
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer JWT token"
// @Success      200 {object} map[string]interface{} "New access token and user ID"
// @Failure      400 {object} map[string]interface{} "Missing auth token"
// @Failure      401 {object} map[string]interface{} "Invalid auth header or token"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /login [get]
func (h *Handler) CheckAuth(c *fiber.Ctx) error {
	header := c.Get("Authorization")

	if header == "" {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusBadRequest})
		logEvent.Msg("Missing auth token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing auth token"})
	}

	tokenString := strings.Split(header, " ")
	if len(tokenString) != 2 {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusUnauthorized})
		logEvent.Msg("Invalid auth header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid auth header"})
	}

	token := tokenString[1]

	id, err := pkg.ParseToken(token, config.SigningKey)
	if err != nil {
		logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Error", Method: c.Method(),
			Url: c.OriginalURL(), Status: fiber.StatusUnauthorized})
		logEvent.Msg(err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
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
			Url: c.OriginalURL(), Status: fiber.StatusInternalServerError})
		logEvent.Msg("user not exists")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "user not exists"})
	}

	logEvent := log.CreateLog(h.logger, log.LogsField{Level: "Info", Method: c.Method(),
		Url: c.OriginalURL(), Status: fiber.StatusOK})
	logEvent.Msg("success")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": token,
		"id":           id,
	})
}
