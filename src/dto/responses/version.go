package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type Version struct {
	ID          uint       `json:"id"`
	PublishID   string     `json:"publishId"`
	AppName     string     `json:"appName"`
	Name        string     `json:"name"`
	EnabledAt   *time.Time `json:"enabledAt"`
	PublishedAt *time.Time `json:"publishedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// SetVersion method to set version data from models.Version{}.
func (v *Version) SetVersion(version *models.Version) {
	v.ID = version.ID
	v.PublishID = version.PublishID.String()
	v.AppName = version.AppName
	v.Name = version.Name
	v.EnabledAt = utils.PtrFromNullTime(version.EnabledAt)
	v.PublishedAt = utils.PtrFromNullTime(version.PublishedAt)
	v.CreatedAt = version.CreatedAt
	v.UpdatedAt = version.UpdatedAt
}
