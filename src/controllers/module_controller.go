package controllers

import (
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/errors"
	"api-page/main/src/services"
	"api-page/main/src/validation"

	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetModules func for getting all modules paginated.
func GetModules(c *fiber.Ctx) error {
	paginationModel, err := services.GetModules(c)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetModuleLookup func for getting module lookup.
func GetModuleLookup(c *fiber.Ctx) error {
	appParam := c.Query("app")
	if appParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "App parameter is required.")
	}

	nameParam := c.Query("name")
	var name *string
	if nameParam != "" {
		name = &nameParam
	}

	modules, err := services.GetModuleLookup(appParam, name)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	response := responses.ModuleLookupList{}
	response.SetModuleLookupList(modules)

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetModuleByID func for getting a module by ID.
func GetModuleByID(c *fiber.Ctx) error {
	moduleIDParam := c.Params("id")
	moduleID, err := util.StringToUint(moduleIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	module, err := services.GetModuleByID(moduleID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if module.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ModuleExists, "Module does not exist.")
	}

	response := responses.Module{}
	response.SetModule(module)

	return c.Status(fiber.StatusOK).JSON(response)
}

// IsModuleNameAvailable method to check if module is available.
func IsModuleNameAvailable(c *fiber.Ctx) error {
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

	if available, err := services.IsModuleNameAvailable(app, name, ignore); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else {
		response := responses.Available{}
		response.SetAvailable(available)

		return c.JSON(response)
	}
}

// CreateModule func for creating a module.
func CreateModule(c *fiber.Ctx) error {
	// Create a new module struct for the request.
	moduleRequest := &requests.CreateModule{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(moduleRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate module fields.
	if err := validation.Validate.Struct(moduleRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if app exists.
	appAvailable, err := services.IsAppAvailable(moduleRequest.AppName)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !appAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppNotFound, "App not found.")
	}

	// Check if module type exists.
	if notAvailable, err := services.IsModuleTypeNotAvailable(moduleRequest.Type); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if notAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleTypeNotFound, "Module type not found.")
	}

	// Check if module exists.
	if available, err := services.IsModuleNameAvailable(moduleRequest.AppName, moduleRequest.Name, nil); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleAvailable, "Module name already exist.")
	}

	// Create module.
	module, err := services.CreateModule(moduleRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the module.
	response := responses.Module{}
	response.SetModule(module)

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateModule func for updating a module.
func UpdateModule(c *fiber.Ctx) error {
	// Get the moduleID parameter from the URL.
	moduleIDParam := c.Params("id")
	moduleID, err := util.StringToUint(moduleIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Create a new module struct for the request.
	moduleRequest := &requests.UpdateModule{}

	// Check, if received JSON data is parsed.
	if err := c.BodyParser(moduleRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate module fields.
	if err := validation.Validate.Struct(moduleRequest); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, util.ValidatorErrors(err))
	}

	// Check if module type exists.
	if notAvailable, err := services.IsModuleTypeNotAvailable(moduleRequest.Type); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if notAvailable {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleTypeNotFound, "Module type not found.")
	}

	// Get old module.
	oldModule, err := services.GetModuleByID(moduleID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if oldModule.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ModuleExists, "Module does not exist.")
	}

	// Check if the module has been modified since it was last fetched.
	if moduleRequest.UpdatedAt.Unix() < oldModule.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Check if module name exists.
	if moduleRequest.Name != oldModule.Name {
		if available, err := services.IsModuleNameAvailable(oldModule.AppName, moduleRequest.Name, &oldModule.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if !available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleAvailable, "Module name already exist.")
		}
	}

	// Update module.
	updatedModule, err := services.UpdateModule(*oldModule, moduleRequest)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the module.
	response := responses.Module{}
	response.SetModule(updatedModule)

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteModule func for deleting a module.
func DeleteModule(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the Module.
	module, err := services.GetModuleByID(id)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if module.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.ModuleExists, "Module does not exist.")
	}

	// Delete the Module.
	if err := services.DeleteModule(module.ID, module.AppName); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreModule func for restoring a deleted module.
func RestoreModule(c *fiber.Ctx) error {
	// Get the ID from the URL.
	id, err := util.StringToUint(c.Params("id"))
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Check if module is deleted.
	if isDeleted, err := services.IsModuleDeleted(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !isDeleted {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.ModuleAvailable, "Module is not deleted.")
	}

	// Restore the module.
	if err := services.RestoreModule(id); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
