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

	// Register route group for /v1/versions.
	versions := route.Group("/versions")
	versions.Get("/", middleware.MachineProtected(), controllers.GetVersions)
	versions.Post("/", middleware.MachineProtected(), controllers.CreateVersion)
	versions.Get("/lookup", middleware.MachineProtected(), controllers.GetVersionLookup)
	versions.Get("/name/available", middleware.MachineProtected(), controllers.IsVersionNameAvailable)
	versions.Get("/:id", middleware.MachineProtected(), controllers.GetVersionByID)
	versions.Put("/:id", middleware.MachineProtected(), controllers.UpdateVersion)
	versions.Delete("/:id", middleware.MachineProtected(), controllers.DeleteVersion)
	versions.Put("/:id/publish", middleware.MachineProtected(), controllers.PublishVersion)
	versions.Put("/:id/restore", middleware.MachineProtected(), controllers.RestoreVersion)
}
