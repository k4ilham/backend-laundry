package routes

import (
	"laundry-backend/handlers"
	"laundry-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/change-password", middleware.Protected(), handlers.ChangePassword)

	// User Management (Protected)
	user := api.Group("/users", middleware.Protected())
	user.Get("/", handlers.GetUsers)
	user.Post("/", handlers.CreateUser)
	user.Get("/:id", handlers.GetUser)
	user.Put("/:id", handlers.UpdateUser)
	user.Delete("/:id", handlers.ArchiveUser)
	user.Post("/:id/restore", handlers.RestoreUser)
	user.Delete("/:id/permanent", handlers.DeleteUser)

	// Service Management (Protected)
	svc := api.Group("/services", middleware.Protected())
	svc.Get("/", handlers.GetServices)
	svc.Post("/", handlers.CreateService)
	svc.Get("/:id", handlers.GetService)
	svc.Put("/:id", handlers.UpdateService)
	svc.Delete("/:id", handlers.ArchiveService)
	svc.Post("/:id/restore", handlers.RestoreService)
	svc.Delete("/:id/permanent", handlers.DeleteService)
}
