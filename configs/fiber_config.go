package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	// Return Fiber configuration.
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}

func FiberLoggerConfig() logger.Config {
	fiberlog.SetLevel(fiberlog.LevelDebug)
	return logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}
}
