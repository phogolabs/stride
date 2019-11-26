package terminal

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
	"github.com/phogolabs/stride/contract"
)

var _ contract.Reporter = &Reporter{}

// Reporter represents the terminal reporter
type Reporter struct {
	Severity contract.Severity
	Writer   io.Writer
}

// With return a reporter with severity
func (r *Reporter) With(value contract.Severity) contract.Reporter {
	return &Reporter{
		Severity: value,
		Writer:   r.Writer,
	}
}

// Notice prints a notice message
func (r *Reporter) Notice(msg string, args ...interface{}) {
	var (
		c      *color.Color
		prefix = r.prefix("notice")
	)

	switch r.Severity {
	case contract.SeverityLow:
		c = color.New(color.Faint, color.FgWhite)
	case contract.SeverityNormal:
		c = color.New(color.FgWhite)
	case contract.SeverityHigh:
		c = color.New(color.FgHiWhite)
	case contract.SeverityVeryHigh:
		c = color.New(color.Bold, color.FgHiWhite)
	default:
		c = color.New(color.FgWhite)
	}

	c.Fprintln(r.Writer, r.text(prefix, msg, args))
}

// Info writes a info
func (r *Reporter) Info(msg string, args ...interface{}) {
	var (
		c      *color.Color
		prefix = r.prefix("info")
	)

	switch r.Severity {
	case contract.SeverityLow:
		c = color.New(color.Faint)
	case contract.SeverityNormal:
		c = color.New()
	case contract.SeverityHigh:
		c = color.New()
	case contract.SeverityVeryHigh:
		c = color.New(color.Bold)
	default:
		c = color.New()
	}

	c.Fprintln(r.Writer, r.text(prefix, msg, args))
	r.clear()
}

// Warn reports a warn level
func (r *Reporter) Warn(msg string, args ...interface{}) {
	var (
		c      *color.Color
		prefix = r.prefix("warn")
	)

	switch r.Severity {
	case contract.SeverityLow:
		c = color.New(color.Faint, color.FgYellow)
	case contract.SeverityNormal:
		c = color.New(color.FgYellow)
	case contract.SeverityHigh:
		c = color.New(color.FgHiYellow)
	case contract.SeverityVeryHigh:
		c = color.New(color.Bold, color.FgHiYellow)
	default:
		c = color.New(color.FgYellow)
	}

	c.Fprintln(r.Writer, r.text(prefix, msg, args))
	r.clear()
}

// Success writes a success
func (r *Reporter) Success(msg string, args ...interface{}) {
	var (
		c      *color.Color
		prefix = r.prefix("success")
	)

	switch r.Severity {
	case contract.SeverityLow:
		c = color.New(color.Faint, color.FgGreen)
	case contract.SeverityNormal:
		c = color.New(color.FgGreen)
	case contract.SeverityHigh:
		c = color.New(color.FgHiGreen)
	case contract.SeverityVeryHigh:
		c = color.New(color.Bold, color.FgHiGreen)
	default:
		c = color.New(color.FgGreen)
	}

	c.Fprintln(r.Writer, r.text(prefix, msg, args))
	r.clear()
}

// Error reports a error level
func (r *Reporter) Error(msg string, args ...interface{}) {
	var (
		c      *color.Color
		prefix = r.prefix("error")
	)

	switch r.Severity {
	case contract.SeverityLow:
		c = color.New(color.Faint, color.FgRed)
	case contract.SeverityNormal:
		c = color.New(color.FgRed)
	case contract.SeverityHigh:
		c = color.New(color.FgHiRed)
	case contract.SeverityVeryHigh:
		c = color.New(color.Bold, color.FgHiRed)
	default:
		c = color.New(color.FgRed)
	}

	c.Fprintln(r.Writer, r.text(prefix, msg, args))
	r.clear()
}

func (r *Reporter) prefix(v string) string {
	switch v {
	case "notice":
		return " "
	case "success":
		return " "
	case "info":
		return "  "
	case "warn":
		return " "
	case "error":
		return " "
	}
	return " "
}

func (r *Reporter) text(prefix, msg string, args []interface{}) string {
	parts := []string{}
	parts = append(parts, prefix)
	parts = append(parts, fmt.Sprintf(msg, args...))
	return strings.Join(parts, " ")
}

func (r *Reporter) clear() {
	color.New(color.Reset).Fprint(r.Writer)
}
