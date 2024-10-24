package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lukinhas563/gochat/src/controller"
	"github.com/lukinhas563/gochat/src/domain"
	"github.com/lukinhas563/gochat/src/model/database/sqlite"
	"github.com/lukinhas563/gochat/src/router"
	"github.com/lukinhas563/gochat/src/shared/service"
	"github.com/lukinhas563/gochat/src/shared/service/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	DB_PATH := os.Getenv("DB_PATH")
	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")
	if DB_PATH == "" || JWT_SECRET_KEY == "" {
		panic("Environment DB_PATH not defined")
	}

	logger.Info("About to start server application")

	server := gin.Default()

	database := sqlite.NewSqliteDatabase()
	if err := database.Connect(DB_PATH); err != nil {
		panic(err)
	}
	defer database.Close()
	logger.Info("Connected on database")

	tokenService := service.NewTokenService(JWT_SECRET_KEY)
	userDomain := domain.NewUserDomain(database, tokenService)
	userController := controller.NewUserController(userDomain)

	router.InitRouter(&server.RouterGroup, userController)

	server.Run()
}
