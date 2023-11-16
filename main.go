package main

import (
	"general/fiber-swagger/configs"
	"general/fiber-swagger/middleware"
	"general/fiber-swagger/routes"
	"general/fiber-swagger/utils"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"os"
)

//	@contact.name
//	@contact.url				http://www.swagger.io/support
//	@contact.email				support@swagger.io
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

func main() {
	// Define a new Fiber app with config.
	app := fiber.New(configs.FiberConfig())
	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.
	//run database
	configs.ConnectDB()
	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.TodoRoutes(app)    // Register a todo routes.
	routes.AuthCallbacks(app) // Register a auth callback routes.

	// Register a route for 404 Error.
	routes.NotFoundRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
