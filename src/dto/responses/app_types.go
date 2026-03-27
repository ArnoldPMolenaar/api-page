package responses

import "api-page/main/src/models"

// AppTypes represents synchronized app type associations.
type AppTypes struct {
	App   string   `json:"app"`
	Types []string `json:"types"`
}

// SetAppModuleTypes maps app and module type models to response fields.
func (at *AppTypes) SetAppModuleTypes(app string, moduleTypes []models.ModuleType) {
	types := make([]string, 0, len(moduleTypes))
	for _, moduleType := range moduleTypes {
		types = append(types, moduleType.Name)
	}

	at.App = app
	at.Types = types
}

// SetAppPluginTypes maps app and plugin type models to response fields.
func (at *AppTypes) SetAppPluginTypes(app string, pluginTypes []models.PluginType) {
	types := make([]string, 0, len(pluginTypes))
	for _, pluginType := range pluginTypes {
		types = append(types, pluginType.Name)
	}

	at.App = app
	at.Types = types
}
