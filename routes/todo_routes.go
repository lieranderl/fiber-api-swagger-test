package routes

import (
	"github.com/gofiber/fiber/v2"

	"general/fiber-swagger/handlers"
	"general/fiber-swagger/middleware"

	swagger_docs_v1 "general/fiber-swagger/docs/v1"
)

// TodoRoutes func for describe group of Todo routes.
func TodoRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group(swagger_docs_v1.SwaggerInfo.BasePath + "/v1")

	route.Get("/todos", middleware.AuthTokenMiddleware(), handlers.GetTodos)              // get list todo
	route.Post("/todos", middleware.AuthTokenMiddleware(), handlers.CreateTodo)           // create new todo
	route.Get("/todos/:id", middleware.AuthTokenMiddleware(), handlers.GetTodoByID)       // get todo by ID
	route.Put("/todos/:id", middleware.AuthTokenMiddleware(), handlers.UpdateTodoByID)    // update one todo by ID
	route.Delete("/todos/:id", middleware.AuthTokenMiddleware(), handlers.DeleteTodoByID) // delete one todo by ID
}
