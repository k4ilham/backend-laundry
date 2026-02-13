package handlers

import (
	"laundry-backend/database"
	"laundry-backend/models"
	"laundry-backend/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Could not login", "data": nil})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": token})
}

func ChangePassword(c *fiber.Ctx) error {
	type ChangePasswordInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	var input ChangePasswordInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid input", "data": nil})
	}

	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	if !utils.CheckPasswordHash(input.OldPassword, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Old password incorrect", "data": nil})
	}

	hashedPassword, _ := utils.HashPassword(input.NewPassword)
	user.Password = hashedPassword
	database.DB.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "Password updated successfully", "data": nil})
}
