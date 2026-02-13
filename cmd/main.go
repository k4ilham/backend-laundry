package main

import (
	"laundry-backend/config"
	"laundry-backend/database"
	"laundry-backend/models"
	"laundry-backend/routes"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Connect DB
	database.Connect()

	// 2. Auto Migrate
	log.Println("Running Migrations...")
	database.DB.AutoMigrate(&models.User{}, &models.Service{}, &models.Transaction{})

	// 3. Seed Data
	database.Seed()

	// 4. Init Fiber
	app := fiber.New()

	// 5. Middleware
	// Logger
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Security Headers
	app.Use(helmet.New())

	// Panic Recovery
	app.Use(recover.New())

	// CORS configuration
	allowedOrigins := config.Get("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "https://frontend-laundry-production.up.railway.app, http://localhost:5173, http://localhost:5174"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Rate Limiting
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many requests, please try again later.",
			})
		},
	}))

	// Health Check Endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "success",
			"message":   "Server is up and running",
			"timestamp": time.Now(),
		})
	})

	// 6. Routes
	routes.SetupRoutes(app)

	// 7. Listen
	port := config.Get("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
