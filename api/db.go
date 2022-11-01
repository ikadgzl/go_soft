package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connects mysql database with gorm and retrying if connection fails
func ConnectToDb() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	timeoutCounter := 0

	for {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err != nil {
			log.Println("Error connecting to database", err)
			timeoutCounter++
		} else {
			log.Println("Connected to database")
			migrateInitialTodos(db)
			return db
		}

		if timeoutCounter == 10 {
			log.Fatal("Could not connect to database")
			return nil
		}

		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
}

// Connects redis database with go-redis and retrying if connection fails
func ConnectToRedis() *redis.Client {
	dsn := os.Getenv("REDIS_URL")
	timeoutCounter := 0

	for {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     dsn,
			Password: "",
			DB:       0,
		})

		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			log.Println("Error connecting to redis", err)
			timeoutCounter++
		} else {
			log.Println("Connected to redis")
			return redisClient
		}

		if timeoutCounter == 10 {
			log.Fatal("Could not connect to database")
			return nil
		}

		log.Println("Retrying in 2 seconds")
		time.Sleep(2 * time.Second)
	}
}

// Mysql migration for initial todos
func migrateInitialTodos(db *gorm.DB) {
	db.Migrator().DropTable(&Todo{})
	db.AutoMigrate(&Todo{})

	initialTodos := []Todo{
		{Title: "Create rest api with golang", IsCompleted: true},
		{Title: "Use gorm", IsCompleted: true},
		{Title: "Dockerize it", IsCompleted: true},
		{Title: "Use redis", IsCompleted: true},
		{Title: "Make it microservice", IsCompleted: false},
		{Title: "write tests", IsCompleted: false}}

	db.CreateInBatches(
		initialTodos,
		100)
}
