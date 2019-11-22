package codegen

import (
	"fmt"
	"strings"

	"github.com/dave/dst"
	"github.com/go-openapi/inflect"
)

func init() {
	inflect.AddAcronym("API")
	inflect.AddAcronym("UUID")
	inflect.AddAcronym("ID")
}

func camelize(text string, tail ...string) string {
	const star = "*"

	items := []string{}
	items = append(items, text)
	items = append(items, tail...)

	text = strings.Join(items, "-")

	text = strings.TrimPrefix(text, star)
	return inflect.Camelize(text)
}

func dasherize(text string) string {
	const star = "*"

	text = strings.TrimPrefix(text, star)
	return inflect.Dasherize(text)
}

func pointer(text string) string {
	const star = "*"

	if !strings.HasPrefix(text, star) {
		text = fmt.Sprintf("%s%s", star, text)
	}

	return text
}

func commentf(decorator *dst.Decorations, text string, args ...interface{}) {
	if text == "" {
		return
	}

	var (
		comments = decorator.All()
		index    = len(comments) - 1
	)

	text = fmt.Sprintf(text, args...)
	text = fmt.Sprintf("// %s", text)

	decorator.Clear()
	decorator.Append(comments[:index]...)
	decorator.Append(text)
	decorator.Append(comments[index:]...)
}
