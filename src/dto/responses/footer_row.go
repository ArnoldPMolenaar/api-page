package responses

import (
	"api-page/main/src/models"
	"time"
)

type FooterRow struct {
	ID              uint              `json:"id"`
	VersionID       uint              `json:"versionId"`
	Locale          string            `json:"locale"`
	Position        uint              `json:"position"`
	NoGutters       bool              `json:"noGutters"`
	Dense           bool              `json:"dense"`
	Hashtag         *string           `json:"hashtag"`
	Align           *string           `json:"align"`
	AlignXxl        *string           `json:"alignXxl"`
	AlignXl         *string           `json:"alignXl"`
	AlignLg         *string           `json:"alignLg"`
	AlignMd         *string           `json:"alignMd"`
	AlignSm         *string           `json:"alignSm"`
	AlignContent    *string           `json:"alignContent"`
	AlignContentXxl *string           `json:"alignContentXxl"`
	AlignContentXl  *string           `json:"alignContentXl"`
	AlignContentLg  *string           `json:"alignContentLg"`
	AlignContentMd  *string           `json:"alignContentMd"`
	AlignContentSm  *string           `json:"alignContentSm"`
	Justify         *string           `json:"justify"`
	JustifyXxl      *string           `json:"justifyXxl"`
	JustifyXl       *string           `json:"justifyXl"`
	JustifyLg       *string           `json:"justifyLg"`
	JustifyMd       *string           `json:"justifyMd"`
	JustifySm       *string           `json:"justifySm"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
	Columns         []FooterRowColumn `json:"columns"`
}

// SetFooterRow sets the FooterRow response from the models.FooterRow model.
func (fr *FooterRow) SetFooterRow(row *models.FooterRow) {
	fr.ID = row.ID
	fr.VersionID = row.VersionID
	fr.Locale = row.Locale
	fr.Position = row.Position
	fr.NoGutters = row.NoGutters
	fr.Dense = row.Dense
	fr.CreatedAt = row.CreatedAt
	fr.UpdatedAt = row.UpdatedAt
	if row.Hashtag.Valid {
		fr.Hashtag = &row.Hashtag.String
	}
	if row.Align.Valid {
		fr.Align = &row.Align.String
	}
	if row.AlignXxl.Valid {
		fr.AlignXxl = &row.AlignXxl.String
	}
	if row.AlignXl.Valid {
		fr.AlignXl = &row.AlignXl.String
	}
	if row.AlignLg.Valid {
		fr.AlignLg = &row.AlignLg.String
	}
	if row.AlignMd.Valid {
		fr.AlignMd = &row.AlignMd.String
	}
	if row.AlignSm.Valid {
		fr.AlignSm = &row.AlignSm.String
	}
	if row.AlignContent.Valid {
		fr.AlignContent = &row.AlignContent.String
	}
	if row.AlignContentXxl.Valid {
		fr.AlignContentXxl = &row.AlignContentXxl.String
	}
	if row.AlignContentXl.Valid {
		fr.AlignContentXl = &row.AlignContentXl.String
	}
	if row.AlignContentLg.Valid {
		fr.AlignContentLg = &row.AlignContentLg.String
	}
	if row.AlignContentMd.Valid {
		fr.AlignContentMd = &row.AlignContentMd.String
	}
	if row.AlignContentSm.Valid {
		fr.AlignContentSm = &row.AlignContentSm.String
	}
	if row.Justify.Valid {
		fr.Justify = &row.Justify.String
	}
	if row.JustifyXxl.Valid {
		fr.JustifyXxl = &row.JustifyXxl.String
	}
	if row.JustifyXl.Valid {
		fr.JustifyXl = &row.JustifyXl.String
	}
	if row.JustifyLg.Valid {
		fr.JustifyLg = &row.JustifyLg.String
	}
	if row.JustifyMd.Valid {
		fr.JustifyMd = &row.JustifyMd.String
	}
	if row.JustifySm.Valid {
		fr.JustifySm = &row.JustifySm.String
	}

	fr.Columns = make([]FooterRowColumn, len(row.Columns))
	for i := range row.Columns {
		frc := FooterRowColumn{}
		frc.SetFooterRowColumn(&row.Columns[i])
		fr.Columns[i] = frc
	}
}
