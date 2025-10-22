package routes

import (
	"hex-arch-fiber/adapter/web/handlers"
	"hex-arch-fiber/port"

	"github.com/gofiber/fiber/v2"
)

func User(app fiber.Router, service port.UserService) {
	app.Get("/users/{username}", handlers.FindByUsername(service))
	app.Get("/users", handlers.FindAll(service))
}
