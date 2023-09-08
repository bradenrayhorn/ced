package ced

const (
	EFORBIDDEN     = "forbidden"
	EINTERNAL      = "internal"
	EINVALID       = "invalid"
	ENOTFOUND      = "not_found"
	EUNAUTHORIZED  = "unauthorized"
	EUNPROCESSABLE = "unprocessable"
)

var (
	ErrorForbidden     = &cedError{code: EFORBIDDEN, msg: "Forbidden"}
	ErrorInternal      = &cedError{code: EINTERNAL, msg: "Internal error"}
	ErrorInvalid       = &cedError{code: EINVALID, msg: "Invalid data provided"}
	ErrorNotFound      = &cedError{code: ENOTFOUND, msg: "Not found"}
	ErrorUnauthorized  = &cedError{code: EUNAUTHORIZED, msg: "Not authenticated"}
	ErrorUnprocessable = &cedError{code: EUNPROCESSABLE, msg: "Unprocessable request"}
)

var codeToError = map[string]*cedError{
	EFORBIDDEN:     ErrorForbidden,
	EINTERNAL:      ErrorInternal,
	EINVALID:       ErrorInvalid,
	ENOTFOUND:      ErrorNotFound,
	EUNAUTHORIZED:  ErrorUnauthorized,
	EUNPROCESSABLE: ErrorUnprocessable,
}

func NewError(code string, msg string) Error {
	err := codeToError[code]
	return WrapError(err, &cedError{code, msg})
}

type Error interface {
	error
	CedError() (string, string)
}

type cedError struct {
	code string
	msg  string
}

func (e cedError) Error() string {
	return e.msg
}

func (e cedError) CedError() (string, string) {
	return e.code, e.msg
}

type wrappedCedError struct {
	error
	cedError Error
}

func (e wrappedCedError) Is(err error) bool {
	return e.cedError == err
}

func (e wrappedCedError) Unwrap() error {
	return e.error
}

func (e wrappedCedError) CedError() (string, string) {
	return e.cedError.CedError()
}

func WrapError(err error, parent Error) Error {
	return wrappedCedError{error: err, cedError: parent}
}
