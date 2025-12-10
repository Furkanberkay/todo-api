package auth

import (
	"net/http"
	"todoApp3/internal/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
}

func NewHandler(service *Service, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

func (h *Handler) RegisterUser(e echo.Context) error {
	userDTO := new(RegisterRequest)
	if err := e.Bind(userDTO); err != nil {
		return httpx.BindErrorResponse(e, err)
	}

	if err := h.validator.Struct(userDTO); err != nil {
		validateError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateError)
	}

	input := userDTO.ToRegisterInput()
	user, err := h.service.RegisterUser(e.Request().Context(), input)

	if err != nil {
		return httpx.HandlerError(e, err)
	}

	response := NewRegisterResponse(user)

	return e.JSON(http.StatusOK, response)
}

func (h *Handler) Login(e echo.Context) error {

	dto := LoginRequest{}
	if err := e.Bind(&dto); err != nil {
		return httpx.BindErrorResponse(e, err)
	}

	if err := h.validator.Struct(&dto); err != nil {
		validationError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validationError)
	}

	loginInput := dto.ToLoginInput()

	loginOutput, errorToken := h.service.Login(e.Request().Context(), loginInput)

	if errorToken != nil {
		return httpx.HandlerError(e, errorToken)
	}

	response := LoginResponse{
		AccessToken: loginOutput.AccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   loginOutput.ExpiresIn,
	}

	return e.JSON(http.StatusOK, response)

}
