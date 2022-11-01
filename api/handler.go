package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
)

// Creates a new todo and before saving it to the database, deletes the cached todos
func CreateTodo(c echo.Context) error {
	redisClient.Del(context.Background(), "todos")

	var todo Todo

	c.Bind(&todo)

	if err := c.Validate(todo); err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{Message: "Invalid request", Data: nil})
	}

	db.Create(&todo)

	return c.JSON(http.StatusCreated, JsonResponse{Message: "Todo created", Data: todo})
}

// First checks cache for todos, if not exists; retrieves all todos from the database and cache them in redis afterwards
func GetAllTodos(c echo.Context) error {
	cacheResponse := parseTodosFromRedis("todos", "")
	if cacheResponse != nil {
		return c.JSON(http.StatusOK, cacheResponse)
	}

	var todos []Todo

	db.Find(&todos)

	if len(todos) == 0 {
		return c.JSON(http.StatusNotFound,
			JsonResponse{Message: "No todos found", Data: nil})
	}

	response, err := json.Marshal(todos)
	if err != nil {
		log.Println("Error marshalling todos", err)
	}

	if redisClient.Set(context.Background(), "todos", response, 10*time.Minute).Err() != nil {
		log.Println("Error caching todos")
	}

	return c.JSON(http.StatusOK, todos)
}

// First checks cache for the todos, if not exists; retrieves all todos from the database and cache them in redis afterwards
func GetTodoById(c echo.Context) error {
	id := c.Param("id")

	cacheResponse := parseTodosFromRedis("todos", id)
	if cacheResponse != nil {
		return c.JSON(http.StatusOK, cacheResponse)
	}

	var todo Todo
	db.First(&todo, id)

	if todo.ID == 0 {
		return c.JSON(http.StatusNotFound,
			JsonResponse{Message: "Todo not found", Data: nil})
	}

	return c.JSON(http.StatusOK, todo)
}

// Updates a todo and before saving it to the database, deletes the cached todos
func UpdateTodoById(c echo.Context) error {
	redisClient.Del(context.Background(), "todos")

	id := c.Param("id")

	var todo Todo
	db.First(&todo, id)

	if todo.ID == 0 {
		return c.JSON(http.StatusNotFound,
			JsonResponse{Message: "Todo not found", Data: nil})
	}

	c.Bind(&todo)
	if err := c.Validate(todo); err != nil {
		return c.JSON(http.StatusBadRequest, JsonResponse{Message: "Invalid request", Data: nil})
	}

	db.Save(&todo)

	return c.JSON(http.StatusOK, todo)
}

// Deletes a todo and before saving it to the database, deletes the cached todos
func DeleteTodoById(c echo.Context) error {
	redisClient.Del(context.Background(), "todos")

	id := c.Param("id")

	var todo Todo
	db.First(&todo, id)

	if todo.ID == 0 {
		return c.JSON(http.StatusNotFound, JsonResponse{Message: "Todo not found", Data: nil})
	}

	db.Delete(&todo)

	return c.JSON(http.StatusOK, JsonResponse{Message: "Todo deleted", Data: todo})
}

// Parses todo from redis based on the id, if id is empty, it returns all todos else looks for the wanted todo by looping the returned cache from redis
func parseTodosFromRedis(key string, id string) *JsonResponse {
	cachedTodosJSON := redisClient.Get(context.Background(), key)
	if cachedTodosJSON.Err() == redis.Nil {
		return nil
	}

	var cachedTodos []Todo

	err := json.Unmarshal([]byte(cachedTodosJSON.Val()), &cachedTodos)
	if err != nil {
		log.Println("Error marshalling cached todos", err)
	}

	if id == "" {
		return &JsonResponse{Message: "Todos retrieved from redis cache", Data: cachedTodos}
	}

	for _, todo := range cachedTodos {
		if strconv.Itoa(todo.ID) == id {
			return &JsonResponse{Message: "Todo retrieved from redis cache", Data: todo}
		}
	}

	return &JsonResponse{Message: "Todo not found", Data: nil}
}
