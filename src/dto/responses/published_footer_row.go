package responses

import (
	"api-page/main/src/models"
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
	if row.Hashtag.Valid {
		pfr.Hashtag = &row.Hashtag.String
	}
	if row.Align.Valid {
		pfr.Align = &row.Align.String
	}
	if row.AlignXxl.Valid {
		pfr.AlignXxl = &row.AlignXxl.String
	}
	if row.AlignXl.Valid {
		pfr.AlignXl = &row.AlignXl.String
	}
	if row.AlignLg.Valid {
		pfr.AlignLg = &row.AlignLg.String
	}
	if row.AlignMd.Valid {
		pfr.AlignMd = &row.AlignMd.String
	}
	if row.AlignSm.Valid {
		pfr.AlignSm = &row.AlignSm.String
	}
	if row.AlignContent.Valid {
		pfr.AlignContent = &row.AlignContent.String
	}
	if row.AlignContentXxl.Valid {
		pfr.AlignContentXxl = &row.AlignContentXxl.String
	}
	if row.AlignContentXl.Valid {
		pfr.AlignContentXl = &row.AlignContentXl.String
	}
	if row.AlignContentLg.Valid {
		pfr.AlignContentLg = &row.AlignContentLg.String
	}
	if row.AlignContentMd.Valid {
		pfr.AlignContentMd = &row.AlignContentMd.String
	}
	if row.AlignContentSm.Valid {
		pfr.AlignContentSm = &row.AlignContentSm.String
	}
	if row.Justify.Valid {
		pfr.Justify = &row.Justify.String
	}
	if row.JustifyXxl.Valid {
		pfr.JustifyXxl = &row.JustifyXxl.String
	}
	if row.JustifyXl.Valid {
		pfr.JustifyXl = &row.JustifyXl.String
	}
	if row.JustifyLg.Valid {
		pfr.JustifyLg = &row.JustifyLg.String
	}
	if row.JustifyMd.Valid {
		pfr.JustifyMd = &row.JustifyMd.String
	}
	if row.JustifySm.Valid {
		pfr.JustifySm = &row.JustifySm.String
	}

	pfr.Columns = make([]PublishedFooterRowColumn, len(row.Columns))
	for i := range row.Columns {
		pfrc := PublishedFooterRowColumn{}
		pfrc.SetFooterRowColumn(&row.Columns[i])
		pfr.Columns[i] = pfrc
	}
}
