package main

import (
	"hex-arch-fiber/adapter/database"
	"hex-arch-fiber/adapter/database/repository"
	"hex-arch-fiber/adapter/web/routes"
	"hex-arch-fiber/application/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db := database.SetupDB()
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)

	app := fiber.New()
	routes.User(app, userService)

	app.Listen(":8080")
}
