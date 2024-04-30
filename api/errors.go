package api

import "net/http"

type Error struct {
	Code 	int		`json:"code"`
	Err 	string		`json:"error"`
}

// impliments Error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err: err,
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err: "Unauthorized",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id given",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid json request",
	}
}

func ErrResourceFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err: res + " resource not found",
	}
}