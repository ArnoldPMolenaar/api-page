package responses

import (
	"api-page/main/src/models"
)

type PublishedPage struct {
	Plugin          *string                `json:"plugin"`
	MetaTitle       *string                `json:"metaTitle"`
	MetaDescription *string                `json:"metaDescription"`
	Indexing        []PageIndexing         `json:"indexing"`
	Partials        []PublishedPagePartial `json:"partials"`
}

// SetPage sets the Page response from models.Page.
func (pp *PublishedPage) SetPage(page *models.Page) {
	if page.Plugin.Valid {
		pp.Plugin = &page.Plugin.String
	}
	if page.MetaTitle.Valid {
		pp.MetaTitle = &page.MetaTitle.String
	}
	if page.MetaDescription.Valid {
		pp.MetaDescription = &page.MetaDescription.String
	}

	pp.Indexing = make([]PageIndexing, len(page.Indexing))
	for i := range page.Indexing {
		pi := PageIndexing{}
		pi.SetPageIndexing(&page.Indexing[i])
		pp.Indexing[i] = pi
	}

	pp.Partials = make([]PublishedPagePartial, len(page.Partials))
	for i := range page.Partials {
		ppp := PublishedPagePartial{}
		ppp.SetPagePartial(&page.Partials[i])
		pp.Partials[i] = ppp
	}
}
