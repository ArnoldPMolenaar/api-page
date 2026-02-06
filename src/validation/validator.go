package validation

import (
	"encoding/json"

	util "github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/go-playground/validator/v10"
)

var Validate = util.NewValidator()

func init() {
	_ = Validate.RegisterValidation("validjson", func(fl validator.FieldLevel) bool {
		raw, ok := fl.Field().Interface().(json.RawMessage)
		if !ok || len(raw) == 0 {
			return false
		}
		var v any
		return json.Unmarshal(raw, &v) == nil
	})
}
