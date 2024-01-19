package errors

import "fmt"

const (
	LoginAlreadyExist = "LOGIN_ALREADY_EXIST"
	UserDoesntExists  = "USER_DOESNT_EXISTS"
)

type BusinessError struct {
	Code    string
	Message string
}

func New(code, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message}
}

func (bs *BusinessError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", bs.Code, bs.Message)
}
