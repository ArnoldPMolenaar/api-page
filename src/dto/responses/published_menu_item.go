package responses

import (
	"api-page/main/src/models"
	"api-page/main/src/utils"
)

type PublishedMenuItem struct {
	ID            uint                `json:"id"`
	Position      uint                `json:"position"`
	Name          string              `json:"name"`
	URLName       string              `json:"urlName"`
	Hashtag       *string             `json:"hashtag"`
	Icon          *string             `json:"icon"`
	Url           *string             `json:"url"`
	UrlEnabled    bool                `json:"urlEnabled"`
	NewTabEnabled bool                `json:"newTabEnabled"`
	Items         []PublishedMenuItem `json:"items"`
}

// SetMenuItem sets the MenuItem response from the models.MenuItem model.
func (pmi *PublishedMenuItem) SetMenuItem(menuItem *models.MenuItem, position uint) {
	var page *models.Page
	if len(menuItem.Pages) > 0 {
		page = &menuItem.Pages[0]
	}

	pmi.ID = menuItem.ID
	pmi.Position = position

	if page != nil {
		pmi.Name = page.Name
		pmi.URLName = utils.URLEncode(page.Name)
		pmi.UrlEnabled = page.UrlEnabled
		pmi.NewTabEnabled = page.NewTabEnabled

		if page.Hashtag.Valid {
			pmi.Hashtag = &page.Hashtag.String
		}

		if page.Url.Valid {
			pmi.Url = &page.Url.String
		}
	}

	if menuItem.Icon.Valid {
		pmi.Icon = &menuItem.Icon.String
	}

	pmi.Items = make([]PublishedMenuItem, 0)
}
