package controllers

import (
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/errors"
	"api-page/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetVersions func for getting all versions paginated.
func GetVersions(c *fiber.Ctx) error {
	paginationModel, err := services.GetVersions(c)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetVersionLookup func for getting version lookup.
func GetVersionLookup(c *fiber.Ctx) error {
	appParam := c.Query("app")
	if appParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "App parameter is required.")
	}

	nameParam := c.Query("name")
	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	versions, err := services.GetVersionLookup(appParam, name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.VersionLookupList{}
	response.SetVersionLookupList(versions)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetVersionByID func for getting a version by ID.
func GetVersionByID(c *fiber.Ctx) error {
	versionIDParam := c.Params("id")
	versionID, err := util.StringToUint(versionIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	version, err := services.GetVersionByID(versionID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Version does not exist.")
	}

	response := responses.Version{}
	response.SetVersion(version)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetPublishedVersionByAppName func for getting the published version by app name.
func GetPublishedVersionByAppName(c *fiber.Ctx) error {
	appName := c.Query("app")
	if appName == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "App name parameter is required.")
	}

	version, err := services.GetPublishedVersionByAppName(appName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Published version does not exist.")
	}

	response := responses.PublishedVersion{}
	response.SetVersion(version)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetMenusByVersionID func for getting menus by version ID.
func GetMenusByVersionID(c *fiber.Ctx) error {
	versionIDParam := c.Params("id")
	versionID, err := util.StringToUint(versionIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	locale := c.Query("locale")
	if locale == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Locale parameter is required.")
	}

	if isPublished, err := services.IsVersionPublished(versionID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isPublished {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionNotPublished, "Version is not published.")
	}

	menus, err := services.GetMenusByVersionID(versionID, locale)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.PublishedMenuList{}
	response.SetMenuList(menus)

	return c.Status(fiber.StatusOK).JSON(response)
}

// IsVersionNameAvailable method to check if version is available.
func IsVersionNameAvailable(c *fiber.Ctx) error {
	app := c.Query("app")
	name := c.Query("name")

	var ignore *string
	ignoreParam := c.Query("ignore", "")
	if ignoreParam != "" {
		ignore = &ignoreParam
	}

	if name == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Name is required.")
	}

	if available, err := services.IsVersionAvailable(app, name, ignore); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else {
		response := responses.Available{}
		response.SetAvailable(available)

		return c.JSON(response)
	}
}

// CreateVersion func for creating a version.
func CreateVersion(c *fiber.Ctx) error {
	// Create a new version struct for the request.
	versionRequest := &requests.CreateVersion{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(versionRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate version fields.
	validate := util.NewValidator()
	if err := validate.Struct(versionRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if app exists.
	appAvailable, err := services.IsAppAvailable(versionRequest.AppName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	// Check if version exists.
	if available, err := services.IsVersionAvailable(versionRequest.AppName, versionRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionAvailable, "Version name already exist.")
	}

	// Create version.
	version, err := services.CreateVersion(versionRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the version.
	response := responses.Version{}
	response.SetVersion(version)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateVersion func for updating a version.
func UpdateVersion(c *fiber.Ctx) error {
	// Get the versionID parameter from the URL.
	versionIDParam := c.Params("id")
	versionID, err := util.StringToUint(versionIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new version struct for the request.
	versionRequest := &requests.UpdateVersion{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(versionRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate version fields.
	validate := util.NewValidator()
	if err := validate.Struct(versionRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Get old version.
	oldVersion, err := services.GetVersionByID(versionID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldVersion.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Version does not exist.")
	}

	// Check if the version has been modified since it was last fetched.
	if versionRequest.UpdatedAt.Unix() < oldVersion.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Check if version exists.
	if versionRequest.Name != oldVersion.Name {
		if available, err := services.IsVersionAvailable(oldVersion.AppName, versionRequest.Name, &oldVersion.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionAvailable, "Version name already exist.")
		}
	}

	if oldVersion.PublishedAt.Valid && versionRequest.EnabledAt == nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionIsPublished, "Published versions must be enabled.")
	}

	// Update version.
	updatedVersion, err := services.UpdateVersion(*oldVersion, versionRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the version.
	response := responses.Version{}
	response.SetVersion(updatedVersion)

	return c.Status(fiber.StatusOK).JSON(response)
}

// PublishVersion func for publishing a version.
func PublishVersion(c *fiber.Ctx) error {
	// Get the versionID parameter from the URL.
	versionIDParam := c.Params("id")
	versionID, err := util.StringToUint(versionIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get version.
	version, err := services.GetVersionByID(versionID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Version does not exist.")
	}

	// Check if version is enabled.
	if !version.EnabledAt.Valid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionNotEnabled, "Version is not enabled.")
	}

	// Check if version is not already published.
	if version.PublishedAt.Valid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionIsPublished, "Version is already published.")
	}

	// Publish version.
	if err := services.PublishVersion(version.AppName, version.ID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteVersion func for deleting a version.
func DeleteVersion(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Version.
	version, err := services.GetVersionByID(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if version.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.VersionExists, "Version does not exist.")
	}

	// Check if it is a published version.
	if version.PublishedAt.Valid {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionIsPublished, "Published versions cannot be deleted.")
	}

	// Delete the Version.
	if err := services.DeleteVersion(version.ID, version.AppName); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreVersion func for restoring a deleted version.
func RestoreVersion(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check if version is deleted.
	if isDeleted, err := services.IsVersionDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.VersionAvailable, "Version is not deleted.")
	}

	// Restore the version.
	if err := services.RestoreVersion(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
