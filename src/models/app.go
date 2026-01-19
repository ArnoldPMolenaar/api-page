package models

type App struct {
	Name string `gorm:"primaryKey:true;autoIncrement:false"`
}
