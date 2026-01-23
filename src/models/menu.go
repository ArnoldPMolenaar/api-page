package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	VersionID uint   `gorm:"not null;index:idx_name,unique"`
	Name      string `gorm:"not null;index:idx_name,unique"`

	// Relationships.
	Version Version `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:VersionID;references:ID"`
}
