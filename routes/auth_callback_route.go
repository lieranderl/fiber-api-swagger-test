package routes

import (
	swagger_docs_v1 "general/fiber-swagger/docs/v1"
	"general/fiber-swagger/handlers"

	"github.com/gofiber/fiber/v2"
)

// AuthCallbacks func for describe group of API Docs routes.
func AuthCallbacks(a *fiber.App) {
	// Create routes group.
	route := a.Group(swagger_docs_v1.SwaggerInfo.BasePath + "/v1/auth")
	// Github callback
	route.Get("/github/callback", handlers.GithubCallback)
}
