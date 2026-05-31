package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PagePartialRow struct {
	ID              uint                   `json:"id"`
	PartialID       uint                   `json:"partialId"`
	Position        uint                   `json:"position"`
	NoGutters       bool                   `json:"noGutters"`
	Dense           bool                   `json:"dense"`
	Hashtag         *string                `json:"hashtag"`
	Align           *string                `json:"align"`
	AlignXxl        *string                `json:"alignXxl"`
	AlignXl         *string                `json:"alignXl"`
	AlignLg         *string                `json:"alignLg"`
	AlignMd         *string                `json:"alignMd"`
	AlignSm         *string                `json:"alignSm"`
	AlignContent    *string                `json:"alignContent"`
	AlignContentXxl *string                `json:"alignContentXxl"`
	AlignContentXl  *string                `json:"alignContentXl"`
	AlignContentLg  *string                `json:"alignContentLg"`
	AlignContentMd  *string                `json:"alignContentMd"`
	AlignContentSm  *string                `json:"alignContentSm"`
	Justify         *string                `json:"justify"`
	JustifyXxl      *string                `json:"justifyXxl"`
	JustifyXl       *string                `json:"justifyXl"`
	JustifyLg       *string                `json:"justifyLg"`
	JustifyMd       *string                `json:"justifyMd"`
	JustifySm       *string                `json:"justifySm"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
	Columns         []PagePartialRowColumn `json:"columns"`
}

// SetPagePartialRow sets the PagePartialRow response from the models.PagePartialRow model.
func (ppr *PagePartialRow) SetPagePartialRow(row *models.PagePartialRow) {
	ppr.ID = row.ID
	ppr.PartialID = row.PagePartialID
	ppr.Position = row.Position
	ppr.NoGutters = row.NoGutters
	ppr.Dense = row.Dense
	ppr.CreatedAt = row.CreatedAt
	ppr.UpdatedAt = row.UpdatedAt
	ppr.Hashtag = utils.PtrFromNullString(row.Hashtag)
	ppr.Align = utils.PtrFromNullString(row.Align)
	ppr.AlignXxl = utils.PtrFromNullString(row.AlignXxl)
	ppr.AlignXl = utils.PtrFromNullString(row.AlignXl)
	ppr.AlignLg = utils.PtrFromNullString(row.AlignLg)
	ppr.AlignMd = utils.PtrFromNullString(row.AlignMd)
	ppr.AlignSm = utils.PtrFromNullString(row.AlignSm)
	ppr.AlignContent = utils.PtrFromNullString(row.AlignContent)
	ppr.AlignContentXxl = utils.PtrFromNullString(row.AlignContentXxl)
	ppr.AlignContentXl = utils.PtrFromNullString(row.AlignContentXl)
	ppr.AlignContentLg = utils.PtrFromNullString(row.AlignContentLg)
	ppr.AlignContentMd = utils.PtrFromNullString(row.AlignContentMd)
	ppr.AlignContentSm = utils.PtrFromNullString(row.AlignContentSm)
	ppr.Justify = utils.PtrFromNullString(row.Justify)
	ppr.JustifyXxl = utils.PtrFromNullString(row.JustifyXxl)
	ppr.JustifyXl = utils.PtrFromNullString(row.JustifyXl)
	ppr.JustifyLg = utils.PtrFromNullString(row.JustifyLg)
	ppr.JustifyMd = utils.PtrFromNullString(row.JustifyMd)
	ppr.JustifySm = utils.PtrFromNullString(row.JustifySm)

	ppr.Columns = make([]PagePartialRowColumn, len(row.Columns))
	for i := range row.Columns {
		pprc := PagePartialRowColumn{}
		pprc.SetPagePartialRowColumn(&row.Columns[i])
		ppr.Columns[i] = pprc
	}
}
