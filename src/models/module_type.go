package models

import "gorm.io/gorm"

type ModuleType struct {
	Name      string         `gorm:"primaryKey:true;autoIncrement:false"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
