package models

type PluginType struct {
	Name string `gorm:"primaryKey:true;autoIncrement:false"`
}
