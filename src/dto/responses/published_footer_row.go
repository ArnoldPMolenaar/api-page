package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedFooterRow struct {
	ID              uint                       `json:"id"`
	Position        uint                       `json:"position"`
	NoGutters       bool                       `json:"noGutters"`
	Dense           bool                       `json:"dense"`
	Hashtag         *string                    `json:"hashtag"`
	Align           *string                    `json:"align"`
	AlignXxl        *string                    `json:"alignXxl"`
	AlignXl         *string                    `json:"alignXl"`
	AlignLg         *string                    `json:"alignLg"`
	AlignMd         *string                    `json:"alignMd"`
	AlignSm         *string                    `json:"alignSm"`
	AlignContent    *string                    `json:"alignContent"`
	AlignContentXxl *string                    `json:"alignContentXxl"`
	AlignContentXl  *string                    `json:"alignContentXl"`
	AlignContentLg  *string                    `json:"alignContentLg"`
	AlignContentMd  *string                    `json:"alignContentMd"`
	AlignContentSm  *string                    `json:"alignContentSm"`
	Justify         *string                    `json:"justify"`
	JustifyXxl      *string                    `json:"justifyXxl"`
	JustifyXl       *string                    `json:"justifyXl"`
	JustifyLg       *string                    `json:"justifyLg"`
	JustifyMd       *string                    `json:"justifyMd"`
	JustifySm       *string                    `json:"justifySm"`
	Columns         []PublishedFooterRowColumn `json:"columns"`
}

// SetFooterRow sets the FooterRow response from the models.FooterRow model.
func (pfr *PublishedFooterRow) SetFooterRow(row *models.FooterRow) {
	pfr.ID = row.ID
	pfr.Position = row.Position
	pfr.NoGutters = row.NoGutters
	pfr.Dense = row.Dense
	pfr.Hashtag = utils.PtrFromNullString(row.Hashtag)
	pfr.Align = utils.PtrFromNullString(row.Align)
	pfr.AlignXxl = utils.PtrFromNullString(row.AlignXxl)
	pfr.AlignXl = utils.PtrFromNullString(row.AlignXl)
	pfr.AlignLg = utils.PtrFromNullString(row.AlignLg)
	pfr.AlignMd = utils.PtrFromNullString(row.AlignMd)
	pfr.AlignSm = utils.PtrFromNullString(row.AlignSm)
	pfr.AlignContent = utils.PtrFromNullString(row.AlignContent)
	pfr.AlignContentXxl = utils.PtrFromNullString(row.AlignContentXxl)
	pfr.AlignContentXl = utils.PtrFromNullString(row.AlignContentXl)
	pfr.AlignContentLg = utils.PtrFromNullString(row.AlignContentLg)
	pfr.AlignContentMd = utils.PtrFromNullString(row.AlignContentMd)
	pfr.AlignContentSm = utils.PtrFromNullString(row.AlignContentSm)
	pfr.Justify = utils.PtrFromNullString(row.Justify)
	pfr.JustifyXxl = utils.PtrFromNullString(row.JustifyXxl)
	pfr.JustifyXl = utils.PtrFromNullString(row.JustifyXl)
	pfr.JustifyLg = utils.PtrFromNullString(row.JustifyLg)
	pfr.JustifyMd = utils.PtrFromNullString(row.JustifyMd)
	pfr.JustifySm = utils.PtrFromNullString(row.JustifySm)

	pfr.Columns = make([]PublishedFooterRowColumn, len(row.Columns))
	for i := range row.Columns {
		pfrc := PublishedFooterRowColumn{}
		pfrc.SetFooterRowColumn(&row.Columns[i])
		pfr.Columns[i] = pfrc
	}
}
