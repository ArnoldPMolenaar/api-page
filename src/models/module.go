package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Module struct {
	gorm.Model
	AppName  string         `gorm:"not null;index:idx_module_name,unique"`
	Type     string         `gorm:"not null"`
	Name     string         `gorm:"not null;index:idx_module_name,unique"`
	Settings datatypes.JSON `gorm:"not null"`

	// Relationships.
	App        App        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
	ModuleType ModuleType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:Type;references:Name"`
}
