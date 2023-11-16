package routes

import (
	"github.com/gofiber/fiber/v2"

	"general/fiber-swagger/handlers"
	"general/fiber-swagger/middleware"
)

// TodoRoutes func for describe group of Todo routes.
func TodoRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	route.Get("/todos", handlers.GetTodos)                                               // get list todo
	route.Post("/todos", middleware.AuthTokenMiddleware(), handlers.CreateTodo)          // create new todo
	route.Get("/todo/:id", handlers.GetTodoByID)                                         // get todo by ID
	route.Put("/todo/:id", middleware.AuthTokenMiddleware(), handlers.UpdateTodoByID)    // update one todo by ID
	route.Delete("/todo/:id", middleware.AuthTokenMiddleware(), handlers.DeleteTodoByID) // delete one todo by ID
}
