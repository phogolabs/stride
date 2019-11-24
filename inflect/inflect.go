package inflect

import (
	"fmt"
	"strings"

	"github.com/go-openapi/inflect"
)

func init() {
	inflect.AddAcronym("API")
	inflect.AddAcronym("UUID")
	inflect.AddAcronym("ID")
}

// Camelize camelizes the text
func Camelize(text string, tail ...string) string {
	const star = "*"

	text = strings.TrimPrefix(text, star)

	items := []string{}
	items = append(items, text)
	items = append(items, tail...)

	text = strings.Join(items, "-")

	return inflect.Camelize(text)
}

// Dasherize dasherizes the text
func Dasherize(text string) string {
	const star = "*"

	text = strings.TrimPrefix(text, star)
	return inflect.Dasherize(text)
}

// Underscore underscores the text
func Underscore(text string) string {
	const star = "*"

	text = strings.TrimPrefix(text, star)
	return inflect.Underscore(text)
}

// Singularize makes a word singularized
func Singularize(word string) string {
	return inflect.Singularize(word)
}

// Pointer makes a type pointer
func Pointer(text string) string {
	const star = "*"

	if !strings.HasPrefix(text, star) {
		text = fmt.Sprintf("%s%s", star, text)
	}

	return text
}
