package middleware

import (
	"general/fiber-swagger/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

const apiKey = "1234"

func AuthTokenMiddleware() func(*fiber.Ctx) error {
	// Create config for Bearer authentication middleware.
	config := keyauth.Config{
		AuthScheme: "Bearer",
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
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
