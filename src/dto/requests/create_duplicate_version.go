package requests

import "time"

// CreateDuplicateVersion represents the request payload for duplicating a version.
type CreateDuplicateVersion struct {
	AppName   string     `json:"appName" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	EnabledAt *time.Time `json:"enabledAt"`
	Footer    bool       `json:"footer"`
	Locales   []string   `json:"locales" validate:"min=1,unique,dive,required"`
	Menus     *[]uint    `json:"menus"`
	MenuItems *[]uint    `json:"menuItems"`
}
