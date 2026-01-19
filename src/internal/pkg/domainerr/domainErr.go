package domainerr

type ErrorType string

const (
	ValidationError ErrorType = "validation"
	NotFoundError   ErrorType = "not_found"
	PermissionError ErrorType = "permission"
	InternalError   ErrorType = "internal_error"
)

type BusinessError struct {
	Code          string
	Type          ErrorType
	Message       string
	InternalError error
}

func (e BusinessError) Error() string {
	return e.Message
}
