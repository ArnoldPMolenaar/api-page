package responses

import (
	"api-page/main/src/models"
	"time"
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
	if row.Hashtag.Valid {
		ppr.Hashtag = &row.Hashtag.String
	}
	if row.Align.Valid {
		ppr.Align = &row.Align.String
	}
	if row.AlignXxl.Valid {
		ppr.AlignXxl = &row.AlignXxl.String
	}
	if row.AlignXl.Valid {
		ppr.AlignXl = &row.AlignXl.String
	}
	if row.AlignLg.Valid {
		ppr.AlignLg = &row.AlignLg.String
	}
	if row.AlignMd.Valid {
		ppr.AlignMd = &row.AlignMd.String
	}
	if row.AlignSm.Valid {
		ppr.AlignSm = &row.AlignSm.String
	}
	if row.AlignContent.Valid {
		ppr.AlignContent = &row.AlignContent.String
	}
	if row.AlignContentXxl.Valid {
		ppr.AlignContentXxl = &row.AlignContentXxl.String
	}
	if row.AlignContentXl.Valid {
		ppr.AlignContentXl = &row.AlignContentXl.String
	}
	if row.AlignContentLg.Valid {
		ppr.AlignContentLg = &row.AlignContentLg.String
	}
	if row.AlignContentMd.Valid {
		ppr.AlignContentMd = &row.AlignContentMd.String
	}
	if row.AlignContentSm.Valid {
		ppr.AlignContentSm = &row.AlignContentSm.String
	}
	if row.Justify.Valid {
		ppr.Justify = &row.Justify.String
	}
	if row.JustifyXxl.Valid {
		ppr.JustifyXxl = &row.JustifyXxl.String
	}
	if row.JustifyXl.Valid {
		ppr.JustifyXl = &row.JustifyXl.String
	}
	if row.JustifyLg.Valid {
		ppr.JustifyLg = &row.JustifyLg.String
	}
	if row.JustifyMd.Valid {
		ppr.JustifyMd = &row.JustifyMd.String
	}
	if row.JustifySm.Valid {
		ppr.JustifySm = &row.JustifySm.String
	}

	ppr.Columns = make([]PagePartialRowColumn, len(row.Columns))
	for i := range row.Columns {
		pprc := PagePartialRowColumn{}
		pprc.SetPagePartialRowColumn(&row.Columns[i])
		ppr.Columns[i] = pprc
	}
}
