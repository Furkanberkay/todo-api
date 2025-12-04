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
	if errors.Is(err, domain.ErrUserAlreadyExists) {
		return e.JSON(http.StatusConflict, ResponseError{
			Message: domain.ErrUserAlreadyExists.Error(),
		})
	}
	return e.JSON(http.StatusInternalServerError, ResponseError{
		Message: domain.InternalError.Error(),
	})
}

func BindErrorResponse(e echo.Context, err error) error {
	return e.JSON(http.StatusBadRequest, ResponseError{
		Message: "invalid error",
	})
}
