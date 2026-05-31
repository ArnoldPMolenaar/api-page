package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type Page struct {
	MenuItemID      uint           `json:"menuItemId"`
	Locale          string         `json:"locale"`
	Plugin          *string        `json:"plugin"`
	Name            string         `json:"name"`
	MetaTitle       *string        `json:"metaTitle"`
	MetaDescription *string        `json:"metaDescription"`
	Hashtag         *string        `json:"hashtag"`
	NewTabEnabled   bool           `json:"newTabEnabled"`
	UrlEnabled      bool           `json:"urlEnabled"`
	Url             *string        `json:"url"`
	EnabledAt       *time.Time     `json:"enabledAt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	Indexing        []PageIndexing `json:"indexing"`
	Partials        []PagePartial  `json:"partials"`
}

// SetPage sets the Page response from models.Page.
func (p *Page) SetPage(page *models.Page) {
	p.MenuItemID = page.MenuItemID
	p.Locale = page.Locale
	p.Plugin = utils.PtrFromNullString(page.Plugin)
	p.Name = page.Name
	p.MetaTitle = utils.PtrFromNullString(page.MetaTitle)
	p.MetaDescription = utils.PtrFromNullString(page.MetaDescription)
	p.Hashtag = utils.PtrFromNullString(page.Hashtag)
	p.NewTabEnabled = page.NewTabEnabled
	p.UrlEnabled = page.UrlEnabled
	p.Url = utils.PtrFromNullString(page.Url)
	p.EnabledAt = utils.PtrFromNullTime(page.EnabledAt)
	p.CreatedAt = page.CreatedAt
	p.UpdatedAt = page.UpdatedAt

	p.Indexing = make([]PageIndexing, len(page.Indexing))
	for i := range page.Indexing {
		pi := PageIndexing{}
		pi.SetPageIndexing(&page.Indexing[i])
		p.Indexing[i] = pi
	}

	p.Partials = make([]PagePartial, len(page.Partials))
	for i := range page.Partials {
		pp := PagePartial{}
		pp.SetPagePartial(&page.Partials[i])
		p.Partials[i] = pp
	}
}
