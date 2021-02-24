package app

import "github.com/go-playground/validator"

// PayloadValidator ...
func PayloadValidator(model interface{}) error {
	validate := validator.New()
	validateError := validate.Struct(model)
	if validateError != nil {
		//	errs := GetErrors(validateError.(validator.ValidationErrors))
		//	RespondWithErrors(w, http.StatusBadRequest, InvalidRequestPayload, errs)
		return validateError
	}
	return nil
}
