package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var errEmptyRequestBody = errors.New("empty request body")

// Because Validator was taken.. sue me.
type Validateable interface {
	Validate() error
}

type Enumer interface {
	Values() []string
}

type Validator struct {
	validator *validator.Validate
	decoder   *json.Decoder
}

func New() *Validator {
	validator := validator.New()

	// Pull JSON tag names from structs we validate.
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Custom validation tags.
	validator.RegisterValidation("isEnum", ValidateEnum, true)
	validator.RegisterValidation("isTimeString", ValidateTimeString, true)

	return &Validator{
		validator: validator,
	}
}

// Validate a struct with validator tags. If it's a struct validator test that first.
func (v *Validator) Validate(val interface{}) error {
	if err := v.validator.Struct(val); err != nil {
		errs := err.(validator.ValidationErrors)
		return formatValidationErrors(errs)
	}

	// Run custom struct validator if val implements Validateable
	inter, isValidateable := val.(Validateable)
	if isValidateable {
		if err := inter.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (v *Validator) ValidateJSON(r io.Reader, val interface{}) error {
	if err := json.NewDecoder(r).Decode(val); err != nil {
		if err == io.EOF {
			return errEmptyRequestBody
		}

		return err
	}
	return v.Validate(val)
}

func ValidateEnum(fl validator.FieldLevel) bool {
	field := fl.Field()

	// If the value is nil it should still pass the enum check.
	// It will fail if it is nil and `required` is present.
	if isFieldNil(fl, field) {
		return true
	}

	if field.Type().Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			v := field.Index(i)

			if !isValidEnumValue(v) {
				return false
			}
		}

		return true
	}

	return isValidEnumValue(field)
}

func isFieldNil(fl validator.FieldLevel, field reflect.Value) bool {
	v, _, nillable := fl.ExtractType(field)
	return nillable && v.IsNil() || field.String() == ""
}

func ValidateTimeString(fl validator.FieldLevel) bool {
	field := fl.Field()

	if isFieldNil(fl, field) {
		return true
	}

	// For some reason we are validating a time string against a non string type.
	if field.Type().Kind() != reflect.String {
		return false
	}

	return isValidTimeString(field.String())
}

// Min: 00:00:00 -> Max 23:59:59
func isValidTimeString(s string) bool {
	match, err := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`, s)
	return err == nil && match
}

func isValidEnumValue(v reflect.Value) bool {
	enumer, isEnumer := v.Interface().(Enumer)
	if !isEnumer {
		return false
	}

	if !containsString(v.String(), enumer.Values()) {
		return false
	}

	return true
}

func containsString(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}
	return false
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	var buf bytes.Buffer
	for i, err := range errs {
		buf.WriteString(fmt.Sprintf("Key: '%s' failed validation on tags: '%s'", err.Field(), err.ActualTag()))

		if i != len(errs)-1 {
			buf.WriteString("\n")
		}
	}
	return errors.New(buf.String())
}
