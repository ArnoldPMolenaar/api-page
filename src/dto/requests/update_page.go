package requests

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type UpdatePage struct {
	Plugin          *string        `json:"plugin"`
	Name            string         `json:"name" validate:"required"`
	MetaTitle       *string        `json:"metaTitle"`
	MetaDescription *string        `json:"metaDescription"`
	Hashtag         *string        `json:"hashtag"`
	NewTabEnabled   bool           `json:"newTabEnabled"`
	UrlEnabled      bool           `json:"urlEnabled"`
	Url             *string        `json:"url"`
	EnabledAt       *time.Time     `json:"enabledAt"`
	UpdatedAt       time.Time      `json:"updatedAt" validate:"required"`
	Indexing        []PageIndexing `json:"indexing" validate:"required,dive"`
}

func (u *UpdatePage) SetPage(page *models.Page) {
	indexing := make([]PageIndexing, 0, len(page.Indexing))
	for i := range page.Indexing {
		indexing = append(indexing, PageIndexing{
			Option: page.Indexing[i].Option.String(),
			Value:  utils.PtrFromNullString(page.Indexing[i].Value),
		})
	}

	u.Plugin = utils.PtrFromNullString(page.Plugin)
	u.Name = page.Name
	u.MetaTitle = utils.PtrFromNullString(page.MetaTitle)
	u.MetaDescription = utils.PtrFromNullString(page.MetaDescription)
	u.Hashtag = utils.PtrFromNullString(page.Hashtag)
	u.NewTabEnabled = page.NewTabEnabled
	u.UrlEnabled = page.UrlEnabled
	u.Url = utils.PtrFromNullString(page.Url)
	u.EnabledAt = utils.PtrFromNullTime(page.EnabledAt)
	u.UpdatedAt = page.UpdatedAt
	u.Indexing = indexing
}
