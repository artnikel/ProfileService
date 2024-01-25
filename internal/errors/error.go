// Package errors contains business errors
package errors

const (
	// LoginAlreadyExist is error code if login already exist
	LoginAlreadyExist = "LOGIN_ALREADY_EXIST"
	// UserDoesntExists is error code if user doesn`t exist
	UserDoesntExists = "USER_DOESNT_EXISTS"
)

// BusinessError is struct for business errors
type BusinessError struct {
	Code string
}

// New is constructor for manage business errors
func New(code string) *BusinessError {
	return &BusinessError{Code: code}
}

// Error is method for creating business errors
func (bs *BusinessError) Error() string {
	return bs.Code
}
