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
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	if err := h.validator.Struct(&userDTO); err != nil {
		validateError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateError)
	}

	input := CreateUserInput{
		Username: userDTO.Username,
		Surname:  userDTO.Surname,
		Name:     userDTO.Name,
		Password: userDTO.Password,
		Email:    userDTO.Email,
	}
	user, err := h.service.RegisterUser(e.Request().Context(), &input)

	if err != nil {
		return httpx.HandlerError(e, err)
	}

	response := RegisterResponse{
		ID:       user.ID,
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Username: user.Username,
	}

	return e.JSON(http.StatusOK, response)
}
