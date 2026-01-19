// Package domainerr provides domain-driven error types for the application.
//
// Error types are classified by semantic category rather than HTTP status codes,
// keeping the domain layer independent of transport layer details.
package domainerr

// ErrorType represents the semantic category of a business error.
//
// ErrorType determines how errors map to HTTP status codes in the interface layer:
//   - ValidationError: HTTP 400 (Bad Request)
//   - NotFoundError: HTTP 404 (Not Found)
//   - PermissionError: HTTP 403 (Forbidden)
//   - ConflictError: HTTP 409 (Conflict)
//   - AuthenticationError: HTTP 401 (Unauthorized)
//   - InternalError: HTTP 500 (Internal Server Error)
type ErrorType string

const (
	// ValidationError indicates invalid input or validation failures.
	ValidationError ErrorType = "validation"

	// NotFoundError indicates a requested resource does not exist.
	NotFoundError ErrorType = "not_found"

	// PermissionError indicates authorization or permission issues.
	PermissionError ErrorType = "permission"

	// ConflictError indicates a conflict with the current state of the resource.
	ConflictError ErrorType = "conflict"

	// AuthenticationError indicates authentication failures.
	AuthenticationError ErrorType = "authentication"

	// InternalError indicates unexpected system errors.
	InternalError ErrorType = "internal_error"
)

// BusinessError represents a domain-level error with structured metadata.
//
// BusinessError should be preferred over plain error values in the domain layer
// to provide richer error context to callers.
type BusinessError struct {
	Code          string    // Unique error code for programmatic handling
	Type          ErrorType // Semantic error category
	Message       string    // Human-readable error message
	InternalError error     // Wrapped underlying error for context
}

// Error returns the error message, including any wrapped internal error.
func (e BusinessError) Error() string {
	if e.InternalError != nil {
		return e.Message + ": " + e.InternalError.Error()
	}
	return e.Message
}

// Unwrap returns the wrapped internal error, supporting Go error chaining.
func (e BusinessError) Unwrap() error {
	return e.InternalError
}
