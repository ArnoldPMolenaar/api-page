package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedMenuItem struct {
	ID            uint                `json:"id"`
	Position      uint                `json:"position"`
	Page          string              `json:"page"`
	Name          string              `json:"name"`
	URLName       string              `json:"urlName"`
	Hashtag       *string             `json:"hashtag"`
	Icon          *string             `json:"icon"`
	IsHome        bool                `json:"isHome"`
	IsCustom      bool                `json:"isCustom"`
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
	pmi.Page = menuItem.Name

	if page != nil {
		pmi.Name = page.Name
		pmi.URLName = utils.URLEncode(page.Name)
		pmi.UrlEnabled = page.UrlEnabled
		pmi.NewTabEnabled = page.NewTabEnabled
		pmi.Hashtag = utils.PtrFromNullString(page.Hashtag)
		pmi.Url = utils.PtrFromNullString(page.Url)
	}

	pmi.Icon = utils.PtrFromNullString(menuItem.Icon)
	pmi.IsHome = menuItem.IsHome
	pmi.IsCustom = menuItem.IsCustom

	pmi.Items = make([]PublishedMenuItem, 0)
}
