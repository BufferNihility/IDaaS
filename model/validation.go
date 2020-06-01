package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jschaefer-io/IDaaS/reponse"
	"io"
	"net/http"
	"strings"
)

// ties to binds the given Json string into the given struct
// and returns the validation errors if there are any
func BindJson(form interface{}, reader io.Reader) error {
	var err error

	// Get json from reader
	jsonString := new(strings.Builder)

	_, err = io.Copy(jsonString, reader)
	if err != nil {
		return err
	}

	// bind json to form struct
	err = json.Unmarshal([]byte(jsonString.String()), form)
	if err != nil {
		return err
	}

	// validates the form struct
	v := validator.New()

	return v.Struct(form)
}

// Generates a validation Error response
func ValidationError(err error) reponse.Response {
	code := http.StatusUnprocessableEntity
	validationError, ok := err.(validator.ValidationErrors)
	if !ok {
		return reponse.NewError(code, err.Error())
	}
	messages := map[string]string{}
	for _, fieldErr := range validationError {
		messages[strings.ToLower(fieldErr.Field())] = fieldError{fieldErr}.String()
	}
	return reponse.NewError(code, messages)
}

// Base wrapper for the
// field validation
type fieldError struct {
	err validator.FieldError
}

// Implement stringer interface
// to cleanup the fieldError.err
func (q fieldError) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("Validation failed on field '%s'", strings.ToLower(q.err.Field())))

	// Print failed condition
	sb.WriteString(fmt.Sprintf(", condition: %s", q.err.ActualTag()))

	// Handle Parameter Validation
	if q.err.Param() != "" {
		sb.WriteString(fmt.Sprintf(" { %s }", q.err.Param()))
	}

	// Print actual value
	if q.err.Value() != nil && q.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", given value: %v", q.err.Value()))
	}
	return sb.String()
}
