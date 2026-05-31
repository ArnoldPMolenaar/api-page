package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedVersion struct {
	ID          uint       `json:"id"`
	PublishID   string     `json:"publishId"`
	Name        string     `json:"name"`
	PublishedAt *time.Time `json:"publishedAt"`
}

// SetVersion method to set version data from models.Version{}.
func (v *PublishedVersion) SetVersion(version *models.Version) {
	v.ID = version.ID
	v.PublishID = version.PublishID.String()
	v.Name = version.Name
	v.PublishedAt = utils.PtrFromNullTime(version.PublishedAt)
}
