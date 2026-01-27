package controllers

import (
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/errors"
	"api-page/main/src/models"
	"api-page/main/src/services"
	"time"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// GetMenu func for getting all menus paginated.
func GetMenu(c *fiber.Ctx) error {
	paginationModel, err := services.GetMenus(c)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetMenuByID func for getting a menu by ID.
func GetMenuByID(c *fiber.Ctx) error {
	menuIDParam := c.Params("id")
	menuID, err := util.StringToUint(menuIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	menu, err := services.GetMenuByID(menuID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if menu.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.MenuExists, "Menu does not exist.")
	}

	response := responses.Menu{}
	response.SetMenu(menu)

	return c.Status(fiber.StatusOK).JSON(response)
}

// IsMenuNameAvailable method to check if menu is available.
func IsMenuNameAvailable(c *fiber.Ctx) error {
	versionIDParam := c.Query("versionId")
	versionID, err := util.StringToUint(versionIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	name := c.Query("name")
	if name == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Name is required.")
	}

	var ignore *string
	ignoreParam := c.Query("ignore", "")
	if ignoreParam != "" {
		ignore = &ignoreParam
	}

	if available, err := services.IsMenuNameAvailable(versionID, name, ignore); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else {
		response := responses.Available{}
		response.SetAvailable(available)

		return c.JSON(response)
	}
}

// CreateMenu func for creating a menu.
func CreateMenu(c *fiber.Ctx) error {
	// Create a new menu struct for the request.
	menuRequest := &requests.CreateMenu{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(menuRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate menu fields.
	validate := util.NewValidator()
	if err := validate.Struct(menuRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if version exists.
	version, err := services.GetVersionByID(menuRequest.VersionID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionExists, "Version not found.")
	}

	// Check if menu name exists.
	if available, err := services.IsMenuNameAvailable(menuRequest.VersionID, menuRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.MenuAvailable, "Menu name already exist.")
	}

	// Create menu.
	menu, err := services.CreateMenu(menuRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the menu.
	response := responses.Menu{}
	response.SetMenu(menu)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateMenu func for updating a menu.
func UpdateMenu(c *fiber.Ctx) error {
	// Get the menuID parameter from the URL.
	menuIDParam := c.Params("id")
	menuID, err := util.StringToUint(menuIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new menu struct for the request.
	menuRequest := &requests.UpdateMenu{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(menuRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate menu fields.
	validate := util.NewValidator()
	if err := validate.Struct(menuRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get old menu.
	oldMenu, err := services.GetMenuByID(menuID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldMenu.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.MenuExists, "Menu does not exist.")
	}

	// Check if the menu has been modified since it was last fetched.
	if menuRequest.UpdatedAt.Unix() < oldMenu.UpdatedAt.Unix() || isMenuItemsOutOfSync(&oldMenu.MenuItemRelations, &menuRequest.Items) {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Check if menu name exists.
	if menuRequest.Name != oldMenu.Name {
		if available, err := services.IsMenuNameAvailable(oldMenu.VersionID, menuRequest.Name, &oldMenu.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.MenuAvailable, "Menu name already exist.")
		}
	}

	// Update menu.
	updatedMenu, err := services.UpdateMenu(oldMenu, menuRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the menu.
	response := responses.Menu{}
	response.SetMenu(updatedMenu)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteMenu func for deleting a menu.
func DeleteMenu(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Menu.
	menu, err := services.GetMenuByID(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if menu.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.MenuExists, "Menu does not exist.")
	}

	// Delete the Menu.
	if err := services.DeleteMenu(menu.ID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreMenu func for restoring a deleted menu.
func RestoreMenu(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check if menu is deleted.
	if isDeleted, err := services.IsMenuDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.MenuAvailable, "Menu is not deleted.")
	}

	// Restore the menu.
	if err := services.RestoreMenu(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// flattenUpdateMenuItems flattens a nested slice of UpdateMenuItem into a single-level slice.
func flattenUpdateMenuItems(items *[]requests.UpdateMenuItem) []requests.UpdateMenuItem {
	if items == nil {
		return nil
	}
	var out []requests.UpdateMenuItem
	var stack []requests.UpdateMenuItem

	// seed stack
	for i := range *items {
		stack = append(stack, (*items)[i])
	}

	// iterative DFS to avoid deep recursion
	for len(stack) > 0 {
		// pop
		n := len(stack) - 1
		cur := stack[n]
		stack = stack[:n]

		out = append(out, cur)

		// push children if present
		if cur.Items != nil {
			for i := range cur.Items {
				stack = append(stack, (cur.Items)[i])
			}
		}
	}

	return out
}

// isMenuItemsOutOfSync checks if any of the menu items are out of sync based on their UpdatedAt timestamps.
func isMenuItemsOutOfSync(menuItems *[]models.MenuItemRelation, requestItems *[]requests.UpdateMenuItem) bool {
	relations := make(map[uint]time.Time)
	for i := range *menuItems {
		relations[(*menuItems)[i].MenuItemChildID] = (*menuItems)[i].MenuItemChild.UpdatedAt
	}

	flattenItems := flattenUpdateMenuItems(requestItems)

	for i := range flattenItems {
		if flattenItems[i].ID != nil {
			if updatedAt, exists := relations[*flattenItems[i].ID]; exists {
				if flattenItems[i].UpdatedAt.Unix() < updatedAt.Unix() {
					log.Debug(*flattenItems[i].ID)
					return true
				}
			}
		}
	}

	return false
}
