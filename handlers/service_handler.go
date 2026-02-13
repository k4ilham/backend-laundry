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

func GetServices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")
	status := c.Query("status", "active") // active, archived, all
	sortBy := c.Query("sort_by", "id")
	order := c.Query("order", "asc")

	offset := (page - 1) * limit

	var services []models.Service
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
		db = db.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	db.Model(&models.Service{}).Count(&total)

	if err := db.Order(sortBy + " " + order).Limit(limit).Offset(offset).Find(&services).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch services", err.Error())
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))

	return utils.SendSuccessResponse(c, "Services fetched successfully", fiber.Map{
		"services": services,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"limit":     limit,
			"last_page": lastPage,
		},
	})
}

func CreateService(c *fiber.Ctx) error {
	var service models.Service
	if err := c.BodyParser(&service); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := database.DB.Create(&service).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create service", err.Error())
	}

	return utils.SendSuccessResponse(c, "Service created successfully", service)
}

func GetService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service models.Service

	if err := database.DB.Unscoped().First(&service, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found", nil)
	}

	return utils.SendSuccessResponse(c, "Service found", service)
}

func UpdateService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service models.Service

	if err := database.DB.First(&service, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found", nil)
	}

	if err := c.BodyParser(&service); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	database.DB.Save(&service)

	return utils.SendSuccessResponse(c, "Service updated successfully", service)
}

func ArchiveService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service models.Service

	if err := database.DB.First(&service, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found", nil)
	}

	database.DB.Delete(&service)

	return utils.SendSuccessResponse(c, "Service archived successfully", nil)
}

func RestoreService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service models.Service

	if err := database.DB.Unscoped().First(&service, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found", nil)
	}

	database.DB.Unscoped().Model(&service).Update("deleted_at", gorm.DeletedAt{})

	return utils.SendSuccessResponse(c, "Service restored successfully", service)
}

func DeleteService(c *fiber.Ctx) error {
	id := c.Params("id")
	var service models.Service

	if err := database.DB.Unscoped().First(&service, id).Error; err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Service not found", nil)
	}

	database.DB.Unscoped().Delete(&service)

	return utils.SendSuccessResponse(c, "Service deleted permanently", nil)
}
