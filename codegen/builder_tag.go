package codegen

import (
	"bytes"
	"fmt"

	"github.com/fatih/structtag"
)

// TagBuilder builds the tags for given type
type TagBuilder struct{}

// Build returns the tag for given property
func (builder *TagBuilder) Build(property *PropertyDescriptor) string {
	tags := &structtag.Tags{}

	tags.Set(&structtag.Tag{
		Key:     "json",
		Name:    property.Name,
		Options: builder.omitempty(property),
	})

	// TODO: uncomment when you add xml support
	// tags.Set(&structtag.Tag{
	// 	Key:     "xml",
	// 	Name:    property.Name,
	// 	Options: builder.omitempty(property),
	// })

	if value := property.PropertyType.Default; value != nil {
		tags.Set(&structtag.Tag{
			Key:  "default",
			Name: fmt.Sprintf("%v", value),
		})
	}

	if options := builder.validate(property); len(options) > 0 {
		tags.Set(&structtag.Tag{
			Key:     "validate",
			Name:    options[0],
			Options: options[1:],
		})
	}

	return fmt.Sprintf("`%s`", tags.String())
}

func (builder *TagBuilder) omitempty(property *PropertyDescriptor) []string {
	options := []string{}
	if !property.Required {
		options = append(options, "omitempty")
	}

	return options
}

func (builder *TagBuilder) validate(property *PropertyDescriptor) []string {
	var (
		options  = []string{}
		metadata = element(property.PropertyType).Metadata
	)

	if property.Required {
		options = append(options, "required")
	}

	for k, v := range metadata {
		switch k {
		case "unique":
			if unique, ok := v.(bool); ok {
				if unique {
					options = append(options, "unique")
				}
			}
		case "pattern":
			if value, ok := v.(string); ok {
				if value != "" {
					// TODO: add support for regex
					// options = append(options, fmt.Sprintf("regex=%v", value))
				}
			}
		case "multiple_of":
			if value, ok := v.(*float64); ok {
				if value != nil {
					// TODO: add support for multileof
					// options = append(options, fmt.Sprintf("multipleof=%v", value))
				}
			}
		case "min":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := metadata["min_exclusive"].(bool); ok {
						if exclusive {
							options = append(options, fmt.Sprintf("gt=%v", *value))
						} else {
							options = append(options, fmt.Sprintf("gte=%v", *value))
						}
					}
				}
			}
		case "max":
			if value, ok := v.(*float64); ok {
				if value != nil {
					if exclusive, ok := metadata["max_exclusive"].(bool); ok {
						if exclusive {
							options = append(options, fmt.Sprintf("lt=%v", *value))
						} else {
							options = append(options, fmt.Sprintf("lte=%v", *value))
						}
					}
				}
			}
		case "values":
			if values, ok := v.([]interface{}); ok {
				if len(values) > 0 {
					options = append(options, fmt.Sprintf("oneof=%v", builder.oneof(values)))
				}
			}
		}
	}

	if len(options) == 0 {
		options = append(options, "-")
	}

	return options
}

func (builder *TagBuilder) oneof(values []interface{}) string {
	buffer := &bytes.Buffer{}

	for index, value := range values {
		if index > 0 {
			fmt.Fprint(buffer, " ")
		}

		fmt.Fprintf(buffer, "%v", value)
	}

	return buffer.String()
}
