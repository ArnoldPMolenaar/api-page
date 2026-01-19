package routes

import (
	"api-page/main/src/controllers"

	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register route group for /v1/apps.
	apps := route.Group("/apps")
	apps.Post("/", middleware.MachineProtected(), controllers.CreateApp)
}
