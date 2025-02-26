package routes

import (
	"gin-app/controllers"
	"gin-app/database"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and defines routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize database connection
	DB := database.GetDB()

	// Initialize TodoController with the database connection
	todoController := controllers.TodoController(DB)

	// Define routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the To-Do API!",
		})
	})

	// To-Do routes
	r.GET("/todos", todoController.GetTodos)
	r.POST("/todos", todoController.CreateTodo)
	r.PUT("/todos/:id", todoController.UpdateTodo)
	r.DELETE("/todos/:id", todoController.DeleteTodo)

	return r
}
