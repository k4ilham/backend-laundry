package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

func SendErrorResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"data":    data,
	})
}
