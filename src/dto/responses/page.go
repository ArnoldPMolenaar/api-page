package responses

import (
	"api-page/main/src/models"
	"time"
)

type Page struct {
	MenuItemID      uint           `json:"menuItemID"`
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
	p.Name = page.Name
	p.NewTabEnabled = page.NewTabEnabled
	p.UrlEnabled = page.UrlEnabled
	p.CreatedAt = page.CreatedAt
	p.UpdatedAt = page.UpdatedAt

	if page.Plugin.Valid {
		p.Plugin = &page.Plugin.String
	}
	if page.MetaTitle.Valid {
		p.MetaTitle = &page.MetaTitle.String
	}
	if page.MetaDescription.Valid {
		p.MetaDescription = &page.MetaDescription.String
	}
	if page.Hashtag.Valid {
		p.Hashtag = &page.Hashtag.String
	}
	if page.Url.Valid {
		p.Url = &page.Url.String
	}
	if page.EnabledAt.Valid {
		p.EnabledAt = &page.EnabledAt.Time
	}

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
