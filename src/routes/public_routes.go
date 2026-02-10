package routes

import (
	"api-page/main/src/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register route group for /v1/versions.
	versions := route.Group("/versions")
	versions.Get("/published", controllers.GetPublishedVersionByAppName)
	versions.Get("/:id/menus", controllers.GetMenusByVersionID)

	// Register route group for v1/pages
	pages := route.Group("/pages")
	pages.Get("/:menuItemId/:locale/published", controllers.GetPublishedPageByID)
}
