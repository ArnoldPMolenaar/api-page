package requests

type UpdatePagePartialRow struct {
	ID        *uint `json:"id"`
	PartialID uint  `json:"partialId" validate:"required"`
	// Use a pointer to uint for Position to allow zero value and required validation.
	Position        *uint                        `json:"position" validate:"required"`
	NoGutters       bool                         `json:"noGutters"`
	Dense           bool                         `json:"dense"`
	Hashtag         *string                      `json:"hashtag"`
	Align           *string                      `json:"align"`
	AlignXxl        *string                      `json:"alignXxl"`
	AlignXl         *string                      `json:"alignXl"`
	AlignLg         *string                      `json:"alignLg"`
	AlignMd         *string                      `json:"alignMd"`
	AlignSm         *string                      `json:"alignSm"`
	AlignContent    *string                      `json:"alignContent"`
	AlignContentXxl *string                      `json:"alignContentXxl"`
	AlignContentXl  *string                      `json:"alignContentXl"`
	AlignContentLg  *string                      `json:"alignContentLg"`
	AlignContentMd  *string                      `json:"alignContentMd"`
	AlignContentSm  *string                      `json:"alignContentSm"`
	Justify         *string                      `json:"justify"`
	JustifyXxl      *string                      `json:"justifyXxl"`
	JustifyXl       *string                      `json:"justifyXl"`
	JustifyLg       *string                      `json:"justifyLg"`
	JustifyMd       *string                      `json:"justifyMd"`
	JustifySm       *string                      `json:"justifySm"`
	Columns         []UpdatePagePartialRowColumn `json:"columns" validate:"required,min=1,dive"`
}
