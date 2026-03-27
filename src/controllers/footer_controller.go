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

// GetPublishedFooterByVersionID retrieves the published Footer for a given version ID and locale, and returns it as a JSON response.
func GetPublishedFooterByVersionID(c *fiber.Ctx) error {
	versionID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Query("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	rows, err := services.GetFooterByVersionID(versionID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.PublishedFooter{}
	response.SetFooter(rows)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetFooterByVersionID retrieves the Footer for a given version ID and locale, and returns it as a JSON response.
func GetFooterByVersionID(c *fiber.Ctx) error {
	versionID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Query("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	rows, err := services.GetFooterByVersionID(versionID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.Footer{}
	response.SetFooter(rows)

	return c.Status(fiber.StatusOK).JSON(response)
}

// UpdateFooter func for updating the rows and columns within the footer.
func UpdateFooter(c *fiber.Ctx) error {
	versionID, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Query("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	// Create a new footer struct for the request.
	footerRequest := &requests.UpdateFooter{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(footerRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate footer fields.
	validate := util.NewValidator()
	if err := validate.Struct(footerRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get version to check if it exists.
	version, err := services.GetVersionByID(versionID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Version not found for the requested footer.")
	}

	// Get old rows.
	oldRows, err := services.GetFooterByVersionID(version.ID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Check if the footer has been modified since it was last fetched.
	if isFooterRowsOutOfSync(oldRows, &footerRequest.Rows) {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update footer.
	updatedFooter, err := services.UpdateFooter(version.ID, locale, oldRows, footerRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the partial.
	response := responses.Footer{}
	response.SetFooter(updatedFooter)

	return c.Status(fiber.StatusOK).JSON(response)
}

// collectFooterRowAndColumnVersions func for collecting the versions of rows and columns in a footer rows tree.
func collectFooterRowAndColumnVersions(rows []models.FooterRow, rowVersions map[uint]int64, columnVersions map[uint]int64, depth int) {
	if depth > services.MaxRowTreeDepth {
		return
	}

	for i := range rows {
		rowVersions[rows[i].ID] = rows[i].UpdatedAt.Unix()
		for j := range rows[i].Columns {
			columnVersions[rows[i].Columns[j].ID] = rows[i].Columns[j].UpdatedAt.Unix()
			collectFooterRowAndColumnVersions(rows[i].Columns[j].FooterRows, rowVersions, columnVersions, depth+1)
		}
	}
}

// areRequestedFooterRowsOutOfSync func for checking if any of the requested rows or their columns have been modified since they were last fetched.
func areRequestedFooterRowsOutOfSync(rows []requests.UpdateFooterRow, rowVersions map[uint]int64, columnVersions map[uint]int64, depth int) bool {
	if depth > services.MaxRowTreeDepth {
		return true
	}

	for i := range rows {
		row := rows[i]
		if row.ID != nil && *row.ID != 0 {
			if row.UpdatedAt == nil {
				return true
			}

			if oldVersion, exists := rowVersions[*row.ID]; exists && row.UpdatedAt.Unix() < oldVersion {
				return true
			}
		}

		for j := range row.Columns {
			col := row.Columns[j]
			if col.ID != nil && *col.ID != 0 {
				if col.UpdatedAt == nil {
					return true
				}

				if oldVersion, exists := columnVersions[*col.ID]; exists && col.UpdatedAt.Unix() < oldVersion {
					return true
				}
			}

			if areRequestedFooterRowsOutOfSync(col.Rows, rowVersions, columnVersions, depth+1) {
				return true
			}
		}
	}

	return false
}

// isFooterRowsOutOfSync checks if any of the footer rows or their columns have been modified since they were last fetched.
func isFooterRowsOutOfSync(oldRows *[]models.FooterRow, newRows *[]requests.UpdateFooterRow) bool {
	rowVersions := map[uint]int64{}
	columnVersions := map[uint]int64{}
	collectFooterRowAndColumnVersions(*oldRows, rowVersions, columnVersions, 1)

	return areRequestedFooterRowsOutOfSync(*newRows, rowVersions, columnVersions, 1)
}
