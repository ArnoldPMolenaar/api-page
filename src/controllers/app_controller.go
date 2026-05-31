package controllers

import (
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/errors"
	"api-page/main/src/services"
	"api-page/main/src/validation"
	"strings"

	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
)

// CreateApp method to create an app.
func CreateApp(c fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateApp{}
	if err := c.Bind().Body(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate app fields.
	validate := util.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Create the app.
	app, err := services.CreateApp(request.Name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err)
	}

	// Return the app.
	response := responses.App{}
	response.SetApp(app)

	return c.JSON(response)
}

// normalizeNames trims whitespace and removes empty entries from a list.
func normalizeNames(values []string) []string {
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		normalized = append(normalized, trimmed)
	}

	return normalized
}

// SetAppModuleTypes parses and validates the request, checks the app exists,
// then synchronizes app -> module_type links to match the provided list.
func SetAppModuleTypes(c fiber.Ctx) error {
	request := &requests.SetAppTypes{}
	if err := c.Bind().Body(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	if err := validation.Validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	request.App = strings.TrimSpace(request.App)
	request.Types = normalizeNames(request.Types)

	appAvailable, err := services.IsAppAvailable(request.App)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}
	if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	response, err := services.SetAppModuleTypes(request.App, request.Types)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleTypeNotFound, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// SetAppPluginTypes parses and validates the request, checks the app exists,
// then synchronizes app -> plugin_type links to match the provided list.
func SetAppPluginTypes(c fiber.Ctx) error {
	request := &requests.SetAppTypes{}
	if err := c.Bind().Body(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	if err := validation.Validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	request.App = strings.TrimSpace(request.App)
	request.Types = normalizeNames(request.Types)

	appAvailable, err := services.IsAppAvailable(request.App)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}
	if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	response, err := services.SetAppPluginTypes(request.App, request.Types)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
