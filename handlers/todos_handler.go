package handlers

import (
	"context"
	"general/fiber-swagger/configs"
	"general/fiber-swagger/models"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary		Create a new Todo
// @Description	Create a new todo item
// @Tags			Todo
// @Accept			json
// @Produce		json
// @Param			todo	body	models.TodoBody	true	"Todo object that needs to be created"
// @Security		BearerAuth
// @Success		200	{object}	models.TodoId
// @Failure		400	{object}	models.ErrorResponse
// @Failure		401	{object}	models.ErrorResponse
// @Failure		403	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/v1/todos [post]
func CreateTodo(c *fiber.Ctx) error {
	// Make a context for the request with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create a new TodoBody struct.
	todo_body := new(models.TodoBody)
	// Parse body into TodoBody struct
	if err := c.BodyParser(todo_body); err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   fiber.ErrBadRequest.Message,
			Details: err.Error(),
		})
	}
	// Make a new todo
	todo := new(models.Todo)
	todo.FetchFromBody(*todo_body)
	todo.User = c.Locals("user").(string)
	todo.Created_time = time.Now()
	// Save the todo to the mongo database
	todosCollection := configs.GetCollection(configs.DB, "todos")
	result, err := todosCollection.InsertOne(ctx, todo)
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   fiber.ErrInternalServerError.Message,
			Details: err.Error(),
		})
	}
	todoID := new(models.TodoId)
	todoID.Id = result.InsertedID.(primitive.ObjectID)
	// Return the new todoID
	return c.Status(fiber.StatusOK).JSON(todoID)
}

// @Summary		Update a Todo
// @Description	Update a todo item
// @Tags			Todo
// @Accept			json
// @Produce		json
// @Param			id		path	string			true	"Todo ID"
// @Param			todo	body	models.TodoBody	true	"Todo object that needs to be updated"
// @Security		BearerAuth
// @Success		200	{object}	models.TodoId
// @Failure		400	{object}	models.ErrorResponse
// @Failure		401	{object}	models.ErrorResponse
// @Failure		403	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/v1/todo/{id} [put]
func UpdateTodoByID(c *fiber.Ctx) error {
	// Make a context for the request with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the todo ID from the request params and make a new TodoId struct.
	todoID := new(models.TodoId)
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	todoID.Id = objId
	todoID.User = c.Locals("user").(string)
	// Create a new TodoBody struct.
	todo_body := new(models.TodoBody)
	// Parse body into TodoBody struct
	if err := c.BodyParser(todo_body); err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   fiber.ErrBadRequest.Message,
			Details: err.Error(),
		})
	}
	// Make a new todo
	todo := new(models.Todo)
	todo.FetchFromBody(*todo_body)
	todo.User = c.Locals("user").(string)
	// Save the todo to the mongo database
	todosCollection := configs.GetCollection(configs.DB, "todos")
	result, err := todosCollection.UpdateOne(ctx, todoID, bson.M{"$set": todo})
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: err.Error(),
		})
	}
	//	if ModifiedCount is 1, then the update was successful, else return 404
	if result.ModifiedCount == 1 {
		return c.Status(fiber.StatusOK).JSON(todoID)
	} else {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: "Sorry, this resource is not found.",
		})
	}
}

// @Summary		Delete a Todo
// @Description	Delete a todo item
// @Tags			Todo
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Todo ID"
// @Security		BearerAuth
// @Success		200	{object}	string
// @Failure		400	{object}	models.ErrorResponse
// @Failure		401	{object}	models.ErrorResponse
// @Failure		403	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/v1/todo/{id} [delete]
func DeleteTodoByID(c *fiber.Ctx) error {
	// Make a context for the request with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the todo ID from the request params and make a new TodoId struct.
	todoID := new(models.TodoId)
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	todoID.Id = objId
	todoID.User = c.Locals("user").(string)
	// Save the todo to the mongo database
	todosCollection := configs.GetCollection(configs.DB, "todos")
	result, err := todosCollection.DeleteOne(ctx, todoID)
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: err.Error(),
		})
	}
	// if DeletedCount is 1, then the delete was successful, else return 404
	if result.DeletedCount == 1 {
		return c.Status(fiber.StatusOK).JSON("Successfully deleted")
	} else {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: "Sorry, this resource is not found.",
		})
	}
}

// @Summary		Get todo by ID
// @Description	Get todo by ID
// @Tags			Todo
// @Accept			json
// @Produce		json
// @Param			id	path	string	true	"Todo ID"
// @Security		BearerAuth
// @Success		200	{object}	models.Todo
// @Failure		400	{object}	models.ErrorResponse
// @Failure		401	{object}	models.ErrorResponse
// @Failure		403	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/v1/todo/{id} [get]
func GetTodoByID(c *fiber.Ctx) error {
	// Make a context for the request with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Get the todo ID from the request params and make a new TodoId struct.
	todoID := new(models.TodoId)
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	todoID.Id = objId
	todoID.User = c.Locals("user").(string)
	// Make a new todo
	todo := new(models.Todo)
	// Get the todo from the mongo database, decode to the todo struct
	todosCollection := configs.GetCollection(configs.DB, "todos")
	err := todosCollection.FindOne(ctx, todoID).Decode(&todo)
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(todo)
}

// @Summary		Get a list of Todos
// @Description	Get a list of Todos
// @Tags			Todo
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Success		200	{object}	[]models.Todo
// @Failure		400	{object}	models.ErrorResponse
// @Failure		401	{object}	models.ErrorResponse
// @Failure		403	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/v1/todos [get]
func GetTodos(c *fiber.Ctx) error {
	// Make a context for the request with a timeout of 10 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Find command and returns a Cursor over the matching documents in the collection.
	todosCollection := configs.GetCollection(configs.DB, "todos")
	cur, err := todosCollection.Find(ctx, bson.M{"user": c.Locals("user").(string)})
	// Check for errors
	if err != nil {
		fiberlog.Error(err)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Error:   fiber.ErrNotFound.Message,
			Details: err.Error(),
		})
	}
	defer cur.Close(ctx)
	//create a slice of todos
	todos := make([]models.Todo, 0)
	//reading from the db in an optimal way, eterating over the cursor and appending to the slice
	for cur.Next(ctx) {
		todo := new(models.Todo)
		// Decode the document into Todo struct, check for errors
		if err = cur.Decode(todo); err != nil {
			fiberlog.Error(err)
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error:   fiber.ErrNotFound.Message,
				Details: err.Error(),
			})
		}
		todos = append(todos, *todo)
	}
	// Return todos
	return c.Status(fiber.StatusOK).JSON(todos)
}
