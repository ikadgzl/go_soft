package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

var (
	e           = echo.New()
	redisClient *redis.Client
	db          *gorm.DB
)

func main() {
	redisClient = ConnectToRedis()
	db = ConnectToDb()

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())

	SetupRoutes()
	e.Logger.Fatal(e.Start(":8080"))
}
