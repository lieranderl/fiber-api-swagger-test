package middleware

import (
	"general/fiber-swagger/models"

	fiberlog "github.com/gofiber/fiber/v2/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

const apiKey = "123456789"

func AuthTokenMiddleware() func(*fiber.Ctx) error {
	// Create config for Bearer authentication middleware.
	config := keyauth.Config{
		AuthScheme: "Bearer",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			fiberlog.Debug("key: ", key)

			if key == apiKey {
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Return status 403 Forbidden.
			return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse{
				Error:   fiber.ErrForbidden.Message,
				Details: "Sorry you don't have permission to access this resource.",
			})
		},
	}

	return keyauth.New(config)
}
