package handlers

import (
	"hex-arch-fiber/adapter/database/entities"
	"hex-arch-fiber/adapter/web/response"
	"hex-arch-fiber/port"

	"github.com/gofiber/fiber/v2"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
)

func FindByUsername(service port.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		param := c.Params("username")

		user, err := service.FindByUsername(c.UserContext(), param)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error:": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(response.UserResponse{User: *user})
	}
}

func FindAll(service port.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		number := c.Query("pageNumber", "1")
		size := c.Query("pageSize", "10")

		pageRequest, err := pagination.PageRequestFrom(
			number,
			size,
			pagination.WithFilter(entities.UserFilter{
				Username:    c.Get("username"),
				PersonEmail: c.Get("person.email"),
			}),
			pagination.WithSortOrder(
				c.Query("order", "id"),
				sort.DirectionFromString(c.Query("direction", "asc")),
			),
		)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error:": err.Error()})
		}

		page, err := service.FindAll(c.UserContext(), pageRequest)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error:": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(response.UsersResponse{
			Users: page.Elements,
			PageInfo: response.PageInfo{
				Number: number,
				Size:   size,
				Total:  page.Total,
			},
		})
	}
}
