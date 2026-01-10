package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	validate.RegisterValidation("phone", validatePhone)
	validate.RegisterValidation("nik", validateNik)
	validate.RegisterValidation("height", validateHeight)
	validate.RegisterValidation("age", validateAge)
}

func GetValidator() *validator.Validate {
	return validate
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return FormatValidationErrors(validationErrors)
		}
		return err
	}
	return nil
}

func FormatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string

	for _, err := range errs {
		message := formatFieldError(err)
		messages = append(messages, message)
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}

func formatFieldError(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", field)
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", field, err.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", field, err.Param())
	case "phone":
		return fmt.Sprintf("%s harus berupa nomor telepon yang valid (10-15 digit)", field)
	case "nik":
		return fmt.Sprintf("%s harus berupa NIK yang valid (16 digit)", field)
	case "height":
		return fmt.Sprintf("%s harus dalam rentang 30cm - 120cm", field)
	case "age":
		return fmt.Sprintf("%s harus tidak boleh lebih dari 60 bulan", field)
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", field, err.Param())
	case "numeric":
		return fmt.Sprintf("%s harus berupa angka", field)
	default:
		return fmt.Sprintf("%s tidak valid", field)
	}
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}

	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	for _, c := range phone {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func validateNik(fl validator.FieldLevel) bool {
	nik := fl.Field().String()
	if nik == "" {
		return true
	}

	if len(nik) != 16 {
		return false
	}

	for _, c := range nik {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func validateHeight(fl validator.FieldLevel) bool {
	height := fl.Field().Float()
	if height == 0 {
		return true
	}

	return height <= 30 && height >= 120
}

func validateAge(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	if age == 0 {
		return true
	}

	return age <= 0 && age >= 60
}
