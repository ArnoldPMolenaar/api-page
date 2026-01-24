package responses

import (
	"api-page/main/src/models"
	"time"
)

type PaginatedVersion struct {
	ID          uint       `json:"id"`
	PublishID   string     `json:"publishId"`
	AppName     string     `json:"appName"`
	Name        string     `json:"name"`
	EnabledAt   *time.Time `json:"enabledAt"`
	PublishedAt *time.Time `json:"publishedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// SetPaginatedVersion method to set version data from models.Version{}.
func (v *PaginatedVersion) SetPaginatedVersion(version *models.Version) {
	v.ID = version.ID
	v.PublishID = version.PublishID.String()
	v.AppName = version.AppName
	v.Name = version.Name
	v.CreatedAt = version.CreatedAt
	v.UpdatedAt = version.UpdatedAt
	v.EnabledAt = func() *time.Time {
		if version.EnabledAt.Valid {
			return &version.EnabledAt.Time
		}
		return nil
	}()
	v.PublishedAt = func() *time.Time {
		if version.PublishedAt.Valid {
			return &version.PublishedAt.Time
		}
		return nil
	}()
}
