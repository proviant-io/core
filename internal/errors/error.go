package errors

import (
	"fmt"
	"github.com/brushknight/proviant/internal/i18n"
)

type CustomError struct {
	message i18n.Message
	code int
}

func (e *CustomError) Error() string{
	return fmt.Sprintf(e.message.Template, e.message.Params)
}

func (e *CustomError) Message() i18n.Message{
	return e.message
}

func (e *CustomError) Code() int{
	return e.code
}

func NewErrNotFound(message i18n.Message) *CustomError {
	return &CustomError{message: message, code: 404}
}

func NewErrBadRequest(message i18n.Message) *CustomError {
	return &CustomError{message: message, code: 400}
}

func NewInternalServer(message i18n.Message) *CustomError {
	return &CustomError{message: message, code: 500}
}