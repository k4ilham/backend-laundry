package handlers

import (
	"laundry-backend/database"
	"laundry-backend/models"
	"laundry-backend/utils"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	role := c.Query("role", "")
	status := c.Query("status", "active") // active, archived, all
	sortBy := c.Query("sort_by", "id")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	db := database.DB

	if status == "archived" {
		db = db.Unscoped().Where("deleted_at IS NOT NULL")
	} else if status == "all" {
		db = db.Unscoped()
	} else {
		db = db.Where("deleted_at IS NULL")
	}

	if search != "" {
		db = db.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if role != "" {
		db = db.Where("role = ?", role)
	}

	db.Model(&models.User{}).Count(&total)

	if err := db.Order(sortBy + " " + order).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users", err.Error())
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return utils.SendSuccessResponse(c, "Users fetched successfully", fiber.Map{
		"users": users,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"last_page": lastPage,
		},
	})
}

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	hashedPassword, _ := utils.HashPassword("password123") // Default password
	user.Password = hashedPassword

	if err := database.DB.Create(&user).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create user", err.Error())
	}

	return utils.SendSuccessResponse(c, "User created successfully", user)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.Unscoped().First(&user, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	return utils.SendSuccessResponse(c, "User found", user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	if err := c.BodyParser(&user); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	database.DB.Save(&user)

	return utils.SendSuccessResponse(c, "User updated successfully", user)
}

func ArchiveUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	database.DB.Delete(&user)

	return utils.SendSuccessResponse(c, "User archived successfully", nil)
}

func RestoreUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.Unscoped().First(&user, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	database.DB.Unscoped().Model(&user).Update("deleted_at", gorm.DeletedAt{})

	return utils.SendSuccessResponse(c, "User restored successfully", user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := database.DB.Unscoped().First(&user, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	database.DB.Unscoped().Delete(&user)

	return utils.SendSuccessResponse(c, "User deleted permanently", nil)
}
