package requests

import "time"

type UpdatePage struct {
	Plugin          *string        `json:"plugin"`
	Name            string         `json:"name" validate:"required"`
	MetaTitle       *string        `json:"metaTitle"`
	MetaDescription *string        `json:"metaDescription"`
	Hashtag         *string        `json:"hashtag"`
	NewTabEnabled   bool           `json:"newTabEnabled"`
	UrlEnabled      bool           `json:"urlEnabled"`
	Url             *string        `json:"url"`
	EnabledAt       *time.Time     `json:"enabledAt"`
	UpdatedAt       time.Time      `json:"updatedAt" validate:"required"`
	Indexing        []PageIndexing `json:"indexing" validate:"required,dive"`
}
