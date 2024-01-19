package errors

const (
	LoginAlreadyExist = "LOGIN_ALREADY_EXIST"
	UserDoesntExists  = "USER_DOESNT_EXISTS"
)

type BusinessError struct {
	Code string
}

func New(code string) *BusinessError {
	return &BusinessError{Code: code}
}

func (bs *BusinessError) Error() string {
	return bs.Code
}
