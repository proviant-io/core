package errors

type CustomError struct {
	message string
	code int
}

func (e *CustomError) Error() string{
	return e.message
}

func (e *CustomError) Code() int{
	return e.code
}

func NewErrNotFound(message string) *CustomError {
	return &CustomError{message: message, code: 404}
}