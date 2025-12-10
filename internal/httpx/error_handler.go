package httpx

import (
	"errors"
	"net/http"
	"todoApp3/internal/domain"

	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

func HandlerError(e echo.Context, err error) error {

	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return e.JSON(http.StatusConflict, ResponseError{
			Message: "User already exists",
		})

	case errors.Is(err, domain.InvalidInput):
		return e.JSON(http.StatusBadRequest, ResponseError{
			Message: err.Error(),
		})

	case errors.Is(err, domain.ErrConflict):
		return e.JSON(http.StatusConflict, ResponseError{
			Message: "Data has been modified by another user. Please refresh.",
		})

	case errors.Is(err, domain.ErrTodoNotFound):
		return e.JSON(http.StatusNotFound, ResponseError{
			Message: "Todo not found",
		})

	case errors.Is(err, domain.ErrInvalidCredentials):
		return e.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Invalid credentials",
		})

	case errors.Is(err, domain.ErrInvalidID):
		return e.JSON(http.StatusUnauthorized, ResponseError{
			Message: "Invalid id",
		})

	default:
		return e.JSON(http.StatusInternalServerError, ResponseError{
			Message: "Internal Server Error",
		})
	}
}

func BindErrorResponse(e echo.Context, err error) error {
	return e.JSON(http.StatusBadRequest, ResponseError{
		Message: "Invalid request format",
	})
}
