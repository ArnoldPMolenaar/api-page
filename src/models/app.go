package models

type App struct {
	Name string `gorm:"primaryKey:true;autoIncrement:false"`

	// Relationships.
	ModuleTypes []ModuleType `gorm:"many2many:app_module_types;foreignKey:Name;joinForeignKey:AppName;references:Name;joinReferences:ModuleType;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PluginTypes []PluginType `gorm:"many2many:app_plugin_types;foreignKey:Name;joinForeignKey:AppName;references:Name;joinReferences:PluginType;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
