package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PageIndexing struct {
	Option string  `json:"option"`
	Value  *string `json:"value"`
}

// SetPageIndexing sets the PageIndexing response from the models.PageIndexing model.
func (pi *PageIndexing) SetPageIndexing(indexing *models.PageIndexing) {
	pi.Option = indexing.Option.String()
	pi.Value = utils.PtrFromNullString(indexing.Value)
}
