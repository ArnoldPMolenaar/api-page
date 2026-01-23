package models

import "gorm.io/gorm"

type PagePartial struct {
	gorm.Model
	MenuItemID uint   `gorm:"not null;index:idx_name,unique"`
	Locale     string `gorm:"not null;size:32;index:idx_name,unique"`
	Name       string `gorm:"not null;index:idx_name,unique"`

	// Relationships.
	MenuItem MenuItem         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID;references:ID"`
	Rows     []PagePartialRow `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PagePartialID;references:ID"`
}
