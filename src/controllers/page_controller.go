package controllers

import (
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/errors"
	"api-page/main/src/models"
	"api-page/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

func GetPublishedPageByID(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	page, err := services.GetPublishedPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Published page not found for the specified menu item and locale.")
	}

	response := responses.PublishedPage{}
	response.SetPage(page)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetOrCreatePageByID func for getting or creating a page.
func GetOrCreatePageByID(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	page, err := services.GetOrCreatePage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page could not be created for the specified menu item and locale.")
	}

	response := responses.Page{}
	response.SetPage(page)

	return c.Status(fiber.StatusOK).JSON(response)
}

// CreatePagePartial func for creating a page partial.
func CreatePagePartial(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	// Create a new partial struct for the request.
	partialRequest := &requests.CreatePagePartial{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(partialRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate partial fields.
	validate := util.NewValidator()
	if err := validate.Struct(partialRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get the page to ensure it exists.
	page, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Check if partial name exists.
	if available, err := services.IsPagePartialNameAvailable(menuItemID, locale, partialRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.PagePartialAvailable, "Partial name already exist.")
	}

	// Create partial.
	partial, err := services.CreatePagePartial(page, partialRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the partial.
	response := responses.PagePartial{}
	response.SetPagePartial(partial)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdatePage func for updating a page.
func UpdatePage(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	// Create a new page struct for the request.
	pageRequest := &requests.UpdatePage{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(pageRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate page fields.
	validate := util.NewValidator()
	if err := validate.Struct(pageRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get old page.
	oldPage, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldPage.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Check if the page has been modified since it was last fetched.
	if pageRequest.UpdatedAt.Unix() < oldPage.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update page.
	updatedPage, err := services.UpdatePage(oldPage, pageRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the page.
	response := responses.Page{}
	response.SetPage(updatedPage)

	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdatePagePartial func for updating a page partial.
func UpdatePagePartial(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	partialID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new partial struct for the request.
	partialRequest := &requests.UpdatePagePartial{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(partialRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate partial fields.
	validate := util.NewValidator()
	if err := validate.Struct(partialRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get page to check if it exists.
	page, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Get old partial.
	oldPartial, err := services.GetPartialByID(partialID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldPartial.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PagePartialAvailable, "Partial does not exist for the specified ID.")
	}

	// Check if the partial has been modified since it was last fetched.
	if partialRequest.UpdatedAt.Unix() < oldPartial.UpdatedAt.Unix() || isPartialRowsOutOfSync(&oldPartial.Rows, &partialRequest.Rows) {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	if partialRequest.Name != oldPartial.Name {
		// Check if partial name exists.
		if available, err := services.IsPagePartialNameAvailable(menuItemID, locale, partialRequest.Name, &oldPartial.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.PagePartialAvailable, "Partial name already exist.")
		}
	}

	// Update partial.
	updatedPartial, err := services.UpdatePagePartial(page.MenuItemID, page.Locale, oldPartial, partialRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the partial.
	response := responses.PagePartial{}
	response.SetPagePartial(updatedPartial)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeletePage func for deleting a page.
func DeletePage(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	// Get page to check if it exists.
	page, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Delete the Page.
	if err := services.DeletePage(menuItemID, locale); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeletePagePartial func for deleting a partial.
func DeletePagePartial(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	partialID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get page to check if it exists.
	page, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Get old partial.
	oldPartial, err := services.GetPartialByID(partialID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldPartial.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PagePartialAvailable, "Partial does not exist for the specified ID.")
	}

	// Delete the Partial.
	if err := services.DeletePagePartial(page.MenuItemID, page.Locale, partialID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestorePage func for restoring a deleted page.
func RestorePage(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	// Check if page is deleted.
	if isDeleted, err := services.IsPageDeleted(menuItemID, locale); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.PageAvailable, "Page is not deleted.")
	}

	// Restore the page.
	if err := services.RestorePage(menuItemID, locale); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestorePagePartial func for restoring a deleted partial.
func RestorePagePartial(c *fiber.Ctx) error {
	menuItemID, err := util.StringToUint(c.Params("menuItemId"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Params("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	partialID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get page to check if it exists.
	page, err := services.GetPage(menuItemID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if page.MenuItemID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.PageExists, "Page not found for the specified menu item and locale.")
	}

	// Check if partial is deleted.
	if isDeleted, err := services.IsPagePartialDeleted(partialID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.PagePartialAvailable, "Partial is not deleted.")
	}

	// Restore the partial.
	if err := services.RestorePagePartial(partialID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// isPartialRowsOutOfSync checks if any of the partial rows or their columns have been modified since they were last fetched.
func isPartialRowsOutOfSync(oldRows *[]models.PagePartialRow, newRows *[]requests.UpdatePagePartialRow) bool {
	rows := &map[uint]models.PagePartialRow{}
	cols := &map[uint]models.PagePartialRowColumn{}
	for _, row := range *oldRows {
		(*rows)[row.ID] = row
		for _, col := range row.Columns {
			(*cols)[col.ID] = col
		}
	}

	for _, newRow := range *newRows {
		if newRow.ID == nil || newRow.UpdatedAt == nil {
			continue
		}

		if oldRow, exists := (*rows)[*newRow.ID]; exists {
			if newRow.UpdatedAt.Unix() < oldRow.UpdatedAt.Unix() {
				return true
			}

			for _, newCol := range newRow.Columns {
				if newCol.ID == nil || newCol.UpdatedAt == nil {
					continue
				}

				if oldCol, exists := (*cols)[*newCol.ID]; exists {
					if newCol.UpdatedAt.Unix() < oldCol.UpdatedAt.Unix() {
						return true
					}
				}
			}
		}
	}

	return false
}
