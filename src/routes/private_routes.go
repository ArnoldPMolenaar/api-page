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

	// Register route group for /v1/menus.
	menus := route.Group("/menus")
	menus.Get("/", middleware.MachineProtected(), controllers.GetMenu)
	menus.Post("/", middleware.MachineProtected(), controllers.CreateMenu)
	menus.Get("/lookup", middleware.MachineProtected(), controllers.GetMenuLookup)
	menus.Get("/name/available", middleware.MachineProtected(), controllers.IsMenuNameAvailable)
	menus.Get("/:id", middleware.MachineProtected(), controllers.GetMenuByID)
	menus.Put("/:id", middleware.MachineProtected(), controllers.UpdateMenu)
	menus.Delete("/:id", middleware.MachineProtected(), controllers.DeleteMenu)
	menus.Put("/:id/restore", middleware.MachineProtected(), controllers.RestoreMenu)

	// Register route group for /v1/menu-items.
	menuItems := route.Group("/menu-items")
	menuItems.Get("/:id/app/available", middleware.MachineProtected(), controllers.IsMenuItemWithAppNameAvailable)

	// Register route group for /v1/pages.
	pages := route.Group("/pages")
	pages.Get("/:menuItemId/:locale", middleware.MachineProtected(), controllers.GetOrCreatePageByID)
	pages.Put("/:menuItemId/:locale", middleware.MachineProtected(), controllers.UpdatePage)
	pages.Delete("/:menuItemId/:locale", middleware.MachineProtected(), controllers.DeletePage)
	pages.Put("/:menuItemId/:locale/restore", middleware.MachineProtected(), controllers.RestorePage)
	pages.Post("/:menuItemId/:locale/partials", middleware.MachineProtected(), controllers.CreatePagePartial)
	pages.Put("/:menuItemId/:locale/partials/:id", middleware.MachineProtected(), controllers.UpdatePagePartial)
	pages.Delete("/:menuItemId/:locale/partials/:id", middleware.MachineProtected(), controllers.DeletePagePartial)
	pages.Put("/:menuItemId/:locale/partials/:id/restore", middleware.MachineProtected(), controllers.RestorePagePartial)

	// Register route group for /v1/modules.
	modules := route.Group("/modules")
	modules.Get("/", middleware.MachineProtected(), controllers.GetModules)
	modules.Post("/", middleware.MachineProtected(), controllers.CreateModule)
	modules.Get("/lookup", middleware.MachineProtected(), controllers.GetModuleLookup)
	modules.Get("/name/available", middleware.MachineProtected(), controllers.IsModuleNameAvailable)
	modules.Get("/:id", middleware.MachineProtected(), controllers.GetModuleByID)
	modules.Put("/:id", middleware.MachineProtected(), controllers.UpdateModule)
	modules.Delete("/:id", middleware.MachineProtected(), controllers.DeleteModule)
	modules.Put("/:id/restore", middleware.MachineProtected(), controllers.RestoreModule)
}
