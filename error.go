package whois

import "fmt"

// ErrorCode defines the type of error occurred during a WHOIS lookup.
type ErrorCode string

const (
	ErrInvalidDomain  ErrorCode = "INVALID_DOMAIN"
	ErrServerNotFound ErrorCode = "SERVER_NOT_FOUND"
	ErrQuery          ErrorCode = "QUERY_ERROR"
)

// Error represents a structured error returned by the whois client.
type Error struct {
	Domain string    // Domain name that was queried
	Server string    // WHOIS server involved in the error (if any)
	Code   ErrorCode // Machine-readable error code
	Err    error     // Original error, for wrapping
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("whois error [%s] (domain: %s, server: %s): %v", e.Code, e.Domain, e.Server, e.Err)
}

// Unwrap allows errors.Is and errors.As to work.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new structured whois error.
func NewError(domain, server string, code ErrorCode, err error) *Error {
	return &Error{
		Domain: domain,
		Server: server,
		Code:   code,
		Err:    err,
	}
}
