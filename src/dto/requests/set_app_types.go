package requests

// SetAppTypes represents the request payload to sync app types associations.
type SetAppTypes struct {
	App   string   `json:"app" validate:"required"`
	Types []string `json:"types" validate:"dive,required"`
}
