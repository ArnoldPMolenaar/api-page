package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedPagePartialRow struct {
	ID              uint                            `json:"id"`
	Position        uint                            `json:"position"`
	NoGutters       bool                            `json:"noGutters"`
	Dense           bool                            `json:"dense"`
	Hashtag         *string                         `json:"hashtag"`
	Align           *string                         `json:"align"`
	AlignXxl        *string                         `json:"alignXxl"`
	AlignXl         *string                         `json:"alignXl"`
	AlignLg         *string                         `json:"alignLg"`
	AlignMd         *string                         `json:"alignMd"`
	AlignSm         *string                         `json:"alignSm"`
	AlignContent    *string                         `json:"alignContent"`
	AlignContentXxl *string                         `json:"alignContentXxl"`
	AlignContentXl  *string                         `json:"alignContentXl"`
	AlignContentLg  *string                         `json:"alignContentLg"`
	AlignContentMd  *string                         `json:"alignContentMd"`
	AlignContentSm  *string                         `json:"alignContentSm"`
	Justify         *string                         `json:"justify"`
	JustifyXxl      *string                         `json:"justifyXxl"`
	JustifyXl       *string                         `json:"justifyXl"`
	JustifyLg       *string                         `json:"justifyLg"`
	JustifyMd       *string                         `json:"justifyMd"`
	JustifySm       *string                         `json:"justifySm"`
	Columns         []PublishedPagePartialRowColumn `json:"columns"`
}

// SetPagePartialRow sets the PagePartialRow response from the models.PagePartialRow model.
func (pppr *PublishedPagePartialRow) SetPagePartialRow(row *models.PagePartialRow) {
	pppr.ID = row.ID
	pppr.Position = row.Position
	pppr.NoGutters = row.NoGutters
	pppr.Dense = row.Dense
	pppr.Hashtag = utils.PtrFromNullString(row.Hashtag)
	pppr.Align = utils.PtrFromNullString(row.Align)
	pppr.AlignXxl = utils.PtrFromNullString(row.AlignXxl)
	pppr.AlignXl = utils.PtrFromNullString(row.AlignXl)
	pppr.AlignLg = utils.PtrFromNullString(row.AlignLg)
	pppr.AlignMd = utils.PtrFromNullString(row.AlignMd)
	pppr.AlignSm = utils.PtrFromNullString(row.AlignSm)
	pppr.AlignContent = utils.PtrFromNullString(row.AlignContent)
	pppr.AlignContentXxl = utils.PtrFromNullString(row.AlignContentXxl)
	pppr.AlignContentXl = utils.PtrFromNullString(row.AlignContentXl)
	pppr.AlignContentLg = utils.PtrFromNullString(row.AlignContentLg)
	pppr.AlignContentMd = utils.PtrFromNullString(row.AlignContentMd)
	pppr.AlignContentSm = utils.PtrFromNullString(row.AlignContentSm)
	pppr.Justify = utils.PtrFromNullString(row.Justify)
	pppr.JustifyXxl = utils.PtrFromNullString(row.JustifyXxl)
	pppr.JustifyXl = utils.PtrFromNullString(row.JustifyXl)
	pppr.JustifyLg = utils.PtrFromNullString(row.JustifyLg)
	pppr.JustifyMd = utils.PtrFromNullString(row.JustifyMd)
	pppr.JustifySm = utils.PtrFromNullString(row.JustifySm)

	pppr.Columns = make([]PublishedPagePartialRowColumn, len(row.Columns))
	for i := range row.Columns {
		ppprc := PublishedPagePartialRowColumn{}
		ppprc.SetPagePartialRowColumn(&row.Columns[i])
		pppr.Columns[i] = ppprc
	}
}
