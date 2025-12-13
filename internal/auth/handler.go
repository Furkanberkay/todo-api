package auth

import (
	"log/slog"
	"net/http"
	"todoApp3/internal/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   *Service
	validator *validator.Validate
	logger    *slog.Logger
}

func NewHandler(service *Service, validator *validator.Validate, logger *slog.Logger) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *Handler) RegisterUser(e echo.Context) error {
	userDTO := new(RegisterRequest)
	if err := e.Bind(userDTO); err != nil {
		return httpx.BindErrorResponse(e)
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
		return httpx.BindErrorResponse(e)
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
		TokenType:       "Bearer",
		AccessToken:     loginOutput.AccessToken,
		RefreshToken:    loginOutput.RefreshToken,
		ExpiresIn:       loginOutput.ExpiresIn,
		RefreshTokenExp: loginOutput.RefreshExpiresAt,
	}

	return e.JSON(http.StatusOK, response)

}

func (h *Handler) Refresh(c echo.Context) error {
	dto := RefreshTokenRequest{}
	if err := c.Bind(&dto); err != nil {
		return httpx.BindErrorResponse(c)
	}

	if err := h.validator.Struct(&dto); err != nil {
		validateErr := httpx.ParseValidationErrors(err)
		return c.JSON(http.StatusBadRequest, validateErr)
	}

	serviceOutput, err := h.service.RefreshTokens(c.Request().Context(), dto.RefreshToken)

	if err != nil {
		return httpx.HandlerError(c, err)
	}

	if serviceOutput == nil {
		h.logger.Error("Critical: Service returned nil output without error")

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	refreshResponse := LoginResponse{
		TokenType:       "Bearer",
		AccessToken:     serviceOutput.AccessToken,
		RefreshToken:    serviceOutput.RefreshToken,
		ExpiresIn:       serviceOutput.ExpiresIn,
		RefreshTokenExp: serviceOutput.RefreshExpiresAt,
	}

	return c.JSON(http.StatusOK, &refreshResponse)
}
