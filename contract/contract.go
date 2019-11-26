package contract

//go:generate counterfeiter -fake-name Reporter -o ../fake/reporter.go . Reporter

// Reporter reports the execution
type Reporter interface {
	// With returns a reporter with given severity
	With(value Severity) Reporter
	// Info prints a info message
	Info(text string, args ...interface{})
	// Notice prints a notice message
	Notice(text string, args ...interface{})
	// Warn prints a warn message
	Warn(text string, args ...interface{})
	// Success prints a success message
	Success(text string, args ...interface{})
	// Error prints a error message
	Error(text string, args ...interface{})
}

// Severity represent a reporter severity
type Severity int

const (
	// SeverityLow represents a low severity
	SeverityLow Severity = -1
	// SeverityNormal represents a normal severity
	SeverityNormal Severity = 0
	// SeverityHigh represents a high severity
	SeverityHigh Severity = 1
	// SeverityVeryHigh represents a very high severity
	SeverityVeryHigh Severity = 2
)
