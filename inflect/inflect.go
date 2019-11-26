package inflect

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-openapi/inflect"
)

func init() {
	inflect.AddAcronym("API")
	inflect.AddAcronym("UUID")
	inflect.AddAcronym("ID")
	inflect.AddAcronym("OK")
}

// Camelize camelizes the text
func Camelize(text string, tail ...string) string {
	const star = "*"

	text = strings.TrimPrefix(text, star)

	items := []string{}
	items = append(items, text)
	items = append(items, tail...)

	text = strings.Join(items, "-")

	items = splitAtCaseChangeWithTitlecase(text)
	text = strings.Join(items, "")

	return text
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

// UpperCase makes the text upper case
func UpperCase(text string) string {
	return strings.ToUpper(text)
}

// LowerCase makes the text lower case
func LowerCase(text string) string {
	return strings.ToLower(text)
}

// Titleize makes the text title
func Titleize(text string) string {
	text = strings.ToLower(text)
	return inflect.Titleize(text)
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

// Unpointer unpoints the text
func Unpointer(text string) string {
	const star = "*"
	return strings.TrimPrefix(text, star)
}

func splitAtCaseChangeWithTitlecase(s string) []string {
	text := func(rn []rune) string {
		word := string(rn)

		switch {
		case strings.EqualFold(word, "id"):
			word = "ID"
		case strings.EqualFold(word, "ok"):
			word = "OK"
		}

		return word
	}

	words := make([]string, 0)
	word := make([]rune, 0)

	for _, c := range s {
		spacer := isSpacerChar(c)
		if len(word) > 0 {
			if unicode.IsUpper(c) || spacer {
				words = append(words, text(word))
				word = make([]rune, 0)
			}
		}

		if !spacer {
			if len(word) > 0 {
				word = append(word, unicode.ToLower(c))
			} else {
				word = append(word, unicode.ToUpper(c))
			}
		}
	}

	words = append(words, text(word))
	return words
}

func isSpacerChar(c rune) bool {
	switch {
	case c == rune("_"[0]):
		return true
	case c == rune(" "[0]):
		return true
	case c == rune(":"[0]):
		return true
	case c == rune("-"[0]):
		return true
	}
	return false
}
