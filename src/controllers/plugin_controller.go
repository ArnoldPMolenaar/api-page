package controllers

import (
	"api-page/main/src/dto/responses"
	"api-page/main/src/services"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/gofiber/fiber/v2"
)

// GetPluginTypeLookup func for getting plugin type lookup.
func GetPluginTypeLookup(c *fiber.Ctx) error {
	appParam := c.Query("app")
	var appName *string
	if appParam != "" {
		appName = &appParam
	}

	pluginTypes, err := services.GetPluginTypeLookup(appName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	types := make([]string, len(*pluginTypes))
	for i, pluginType := range *pluginTypes {
		types[i] = pluginType.Name
	}

	response := responses.TypeLookupList{}
	response.SetTypeLookupList(types)

	return c.Status(fiber.StatusOK).JSON(response)
}
