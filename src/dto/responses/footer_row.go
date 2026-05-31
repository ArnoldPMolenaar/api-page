package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
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
	fr.Hashtag = utils.PtrFromNullString(row.Hashtag)
	fr.Align = utils.PtrFromNullString(row.Align)
	fr.AlignXxl = utils.PtrFromNullString(row.AlignXxl)
	fr.AlignXl = utils.PtrFromNullString(row.AlignXl)
	fr.AlignLg = utils.PtrFromNullString(row.AlignLg)
	fr.AlignMd = utils.PtrFromNullString(row.AlignMd)
	fr.AlignSm = utils.PtrFromNullString(row.AlignSm)
	fr.AlignContent = utils.PtrFromNullString(row.AlignContent)
	fr.AlignContentXxl = utils.PtrFromNullString(row.AlignContentXxl)
	fr.AlignContentXl = utils.PtrFromNullString(row.AlignContentXl)
	fr.AlignContentLg = utils.PtrFromNullString(row.AlignContentLg)
	fr.AlignContentMd = utils.PtrFromNullString(row.AlignContentMd)
	fr.AlignContentSm = utils.PtrFromNullString(row.AlignContentSm)
	fr.Justify = utils.PtrFromNullString(row.Justify)
	fr.JustifyXxl = utils.PtrFromNullString(row.JustifyXxl)
	fr.JustifyXl = utils.PtrFromNullString(row.JustifyXl)
	fr.JustifyLg = utils.PtrFromNullString(row.JustifyLg)
	fr.JustifyMd = utils.PtrFromNullString(row.JustifyMd)
	fr.JustifySm = utils.PtrFromNullString(row.JustifySm)

	fr.Columns = make([]FooterRowColumn, len(row.Columns))
	for i := range row.Columns {
		frc := FooterRowColumn{}
		frc.SetFooterRowColumn(&row.Columns[i])
		fr.Columns[i] = frc
	}
}
