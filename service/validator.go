package service

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

// Validator validates a swagger file
type Validator struct {
	Path string
}

// Validate validates the file
func (v *Validator) Validate() error {
	loader := openapi3.NewSwaggerLoader()

	spec, err := loader.LoadSwaggerFromFile(v.Path)
	if err != nil {
		return err
	}

	return spec.Validate(context.TODO())
}
