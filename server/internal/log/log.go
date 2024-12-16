package log

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"server/internal/config"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type LogsField struct {
	Level  string
	Method string
	Url    string
	Status int
}

func InitLogger() *zerolog.Logger {
	if config.ProductionType == "prod" {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   "/var/log/app/iqj.log",
			MaxSize:    1024,
			MaxAge:     183,
			MaxBackups: 5,
			Compress:   true,
		}

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger := zerolog.New(lumberjackLogger).With().Timestamp().Logger()

		return &logger
	}
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &logger
}

func CreateLog(log *zerolog.Logger, field LogsField) *zerolog.Event {
	var event *zerolog.Event
	if field.Level == "Info" {
		event = log.Info().Caller()
	} else if field.Level == "Error" {
		event = log.Error().Caller()
	} else if field.Level == "Warn" {
		event = log.Warn().Caller()
	} else if field.Level == "Debug" {
		event = log.Debug().Caller()
	} else if field.Level == "Fatal" {
		event = log.Fatal().Caller()
	} else {
		fmt.Println("Unknown log level")
		return nil
	}

	event.Str("method", field.Method).Str("url", field.Url).Int("status", field.Status)

	return event
}

func RequestLogger(log *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Info().
			Str("method", c.Method()).
			Str("url", c.OriginalURL()).
			Msg("incoming request")

		start := time.Now()
		defer func() {
			if time.Since(start) > time.Second*2 {
				log.Warn().
					Str("method", c.Method()).
					Str("url", c.OriginalURL()).
					Dur("elapsed_ms", time.Since(start)).
					Msg("long response time")
			}
		}()

		return c.Next()
	}
}
