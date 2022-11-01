package main

func SetupRoutes() {
	e.GET("/todos", GetAllTodos)
	e.GET("/todos/:id", GetTodoById)
	e.POST("/todos", CreateTodo)
	e.PUT("/todos/:id", UpdateTodoById)
	e.DELETE("/todos/:id", DeleteTodoById)
}
