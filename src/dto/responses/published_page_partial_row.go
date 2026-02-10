package responses

import (
	"api-page/main/src/models"
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
	if row.Hashtag.Valid {
		pppr.Hashtag = &row.Hashtag.String
	}
	if row.Align.Valid {
		pppr.Align = &row.Align.String
	}
	if row.AlignXxl.Valid {
		pppr.AlignXxl = &row.AlignXxl.String
	}
	if row.AlignXl.Valid {
		pppr.AlignXl = &row.AlignXl.String
	}
	if row.AlignLg.Valid {
		pppr.AlignLg = &row.AlignLg.String
	}
	if row.AlignMd.Valid {
		pppr.AlignMd = &row.AlignMd.String
	}
	if row.AlignSm.Valid {
		pppr.AlignSm = &row.AlignSm.String
	}
	if row.AlignContent.Valid {
		pppr.AlignContent = &row.AlignContent.String
	}
	if row.AlignContentXxl.Valid {
		pppr.AlignContentXxl = &row.AlignContentXxl.String
	}
	if row.AlignContentXl.Valid {
		pppr.AlignContentXl = &row.AlignContentXl.String
	}
	if row.AlignContentLg.Valid {
		pppr.AlignContentLg = &row.AlignContentLg.String
	}
	if row.AlignContentMd.Valid {
		pppr.AlignContentMd = &row.AlignContentMd.String
	}
	if row.AlignContentSm.Valid {
		pppr.AlignContentSm = &row.AlignContentSm.String
	}
	if row.Justify.Valid {
		pppr.Justify = &row.Justify.String
	}
	if row.JustifyXxl.Valid {
		pppr.JustifyXxl = &row.JustifyXxl.String
	}
	if row.JustifyXl.Valid {
		pppr.JustifyXl = &row.JustifyXl.String
	}
	if row.JustifyLg.Valid {
		pppr.JustifyLg = &row.JustifyLg.String
	}
	if row.JustifyMd.Valid {
		pppr.JustifyMd = &row.JustifyMd.String
	}
	if row.JustifySm.Valid {
		pppr.JustifySm = &row.JustifySm.String
	}

	pppr.Columns = make([]PublishedPagePartialRowColumn, len(row.Columns))
	for i := range row.Columns {
		ppprc := PublishedPagePartialRowColumn{}
		ppprc.SetPagePartialRowColumn(&row.Columns[i])
		pppr.Columns[i] = ppprc
	}
}
