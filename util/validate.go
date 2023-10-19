package util

import (
	"app/config"
	"log"

	"github.com/go-playground/validator/v10"
)

func Validate(r interface{}) error {
	// Validation
	validate := validator.New()
	validate_err := validate.Struct(r)

	if validate_err != nil {
		if config.DEBUG {
			log.Println("Validation Error: ", validate_err)
		}
		return validate_err
	}
	return nil
}
