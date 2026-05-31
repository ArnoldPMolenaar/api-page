package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
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
	pp.Plugin = utils.PtrFromNullString(page.Plugin)
	pp.MetaTitle = utils.PtrFromNullString(page.MetaTitle)
	pp.MetaDescription = utils.PtrFromNullString(page.MetaDescription)

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
