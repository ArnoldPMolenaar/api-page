package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Page struct {
	MenuItemID      uint   `gorm:"primaryKey:true;autoIncrement:false"`
	Locale          string `gorm:"primaryKey:true;autoIncrement:false;size:32"`
	Plugin          sql.NullString
	Name            string `gorm:"not null"`
	MetaTitle       sql.NullString
	MetaDescription sql.NullString
	Hashtag         sql.NullString
	NewTabEnabled   bool `gorm:"not null;default:false"`
	UrlEnabled      bool `gorm:"not null;default:false"`
	Url             sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships.
	MenuItem   MenuItem       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID;references:ID"`
	PluginType PluginType     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:Plugin;references:Name"`
	Indexing   []PageIndexing `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID,Locale;references:MenuItemID,Locale"`
	Partials   []PagePartial  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID,Locale;references:MenuItemID,Locale"`
}
