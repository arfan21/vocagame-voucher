package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrEmailAlreadyRegistered  = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "email already registered"}
	ErrEmailOrPasswordInvalid  = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "email or password invalid"}
	ErrInvalidUUID             = errors.New("invalid uuid length or format")
	ErrUnauthorizedAccess      = &ErrWithCode{HTTPStatusCode: http.StatusUnauthorized, Message: "unauthorized access"}
	ErrOutOfStock              = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "out of stock"}
	ErrProductNotFound         = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "product not found"}
	ErrPaymentMethodNotFound   = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "payment method not found"}
	ErrPaymentMethoUnavailable = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "payment method not available"}
	ErrInvalidSignatureKey     = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "invalid signature key"}
)

type ErrWithCode struct {
	HTTPStatusCode int
	Message        string
}

func (e *ErrWithCode) Error() string {
	return e.Message
}
