package codegen

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-openapi/inflect"
)

func camelize(text string) string {
	var (
		field  = inflect.Camelize(text)
		buffer = &bytes.Buffer{}
		suffix = "Id"
	)

	switch {
	case field == suffix:
		buffer.WriteString(strings.ToUpper(field))
	case strings.HasSuffix(field, suffix):
		buffer.WriteString(strings.TrimSuffix(field, suffix))
		buffer.WriteString(strings.ToUpper(suffix))
	default:
		buffer.WriteString(field)
	}

	return buffer.String()
}

func commentf(text string, args ...interface{}) string {
	text = fmt.Sprintf(text, args...)
	return fmt.Sprintf("// %s", text)
}

func pointer(text string) string {
	return fmt.Sprintf("*%s", text)
}
