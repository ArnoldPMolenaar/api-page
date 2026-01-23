package models

import (
	"database/sql"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Version struct {
	gorm.Model
	EnabledAt   sql.NullTime
	PublishedAt sql.NullTime
	PublishID   datatypes.UUID `gorm:"not null;uniqueIndex"`
	AppName     string         `gorm:"not null;index:idx_name,unique"`
	Name        string         `gorm:"not null;index:idx_name,unique"`

	// Relationships.
	App App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppName;references:Name"`
}

// BeforeCreate sets a time-ordered UUIDv7 for PublishID if it's not already set.
func (v *Version) BeforeCreate(_ *gorm.DB) error {
	if v.PublishID.IsNil() {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		v.PublishID = datatypes.UUID(id)
	}
	return nil
}
