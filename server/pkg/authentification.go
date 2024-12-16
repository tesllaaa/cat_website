package pkg

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

// tokenClaims Структура для полей токена
type tokenClaims struct {
	jwt.MapClaims
	UserId int `json:"user_id"`
}

// WithJWTAuth Middleware аутентификации
func WithJWTAuth(c *fiber.Ctx, signingKey string) error {
	header := c.Get("Authorization")

	if header == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing auth token"})
	}

	tokenString := strings.Split(header, " ")

	if len(tokenString) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid auth header"})
	}

	id, err := ParseToken(tokenString[1], signingKey)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	// Записываем id в контекст, чтобы в дальнейшем использовать в других функциях
	c.Locals("id", id)
	return c.Next()
}

// GenerateAccessToken Генрация аксес токена
func GenerateAccessToken(id, expirationTime int, signingKey string) (string, error) {
	claims := &tokenClaims{
		jwt.MapClaims{
			"ExpiresAt": time.Now().Add(time.Duration(expirationTime) * time.Hour).Unix(),
			"IssuedAr":  time.Now().Unix(),
		},
		id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// GenerateRefreshToken Генерация рефреш токена
func GenerateRefreshToken(id int, signingKey string) (string, error) {
	claims := &tokenClaims{
		jwt.MapClaims{
			"ExpiresAt": time.Now().Add(724 * time.Hour).Unix(),
			"IssuedAr":  time.Now().Unix(),
		},
		id,
	}
	// Создание токена с параметрами записанными в claims и uid пользователя
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// ParseToken Парсинг токена и получение id пользователя
func ParseToken(tokenString string, signingKey string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	if time.Now().Unix() > int64(claims.MapClaims["ExpiresAt"].(float64)) {
		return 0, errors.New("Token has expired")
	}

	return claims.UserId, nil
}
