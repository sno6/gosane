package http

var (
	Ok              = &Error{Status: StatusOk, Message: "200 OK"}
	BadRequest      = &Error{Status: StatusBadRequest, Message: "Bad request"}
	Unauthorized    = &Error{Status: StatusUnauthorized, Message: "Unauthorized"}
	NotFound        = &Error{Status: StatusNotFound, Message: "Not found"}
	Internal        = &Error{Status: StatusInternal, Message: "Internal server error"}
	Forbidden       = &Error{Status: StatusForbidden, Message: "Forbidden"}
	PaymentRequired = &Error{Status: StatusPaymentRequired, Message: "Payment Required"}
	Teapot          = &Error{Status: StatusTeapot, Message: "Teapot"}
)

type Error struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta"`
}

func (e *Error) Error() string {
	return e.Message
}

func ErrorFromString(s string) *Error {
	for _, err := range []*Error{
		Ok,
		BadRequest,
		Teapot,
		Forbidden,
		Unauthorized,
		NotFound,
		Internal,
		PaymentRequired,
	} {
		if s == err.Error() {
			return err
		}
	}
	// If a given error can not be found return an internal server error.
	return Internal
}
