package auth

import (
	"context"
	"log/slog"
	"net/http"
	"todoApp3/internal/domain"
	"todoApp3/internal/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthenticationService interface {
	Register(ctx context.Context, input *CreateUserInput) (*domain.User, error)
	Login(ctx context.Context, loginInput *LoginInput) (*LoginOutput, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*LoginOutput, error)
	Delete(ctx context.Context, role string, id uint) error
}

type Handler struct {
	service   AuthenticationService
	validator *validator.Validate
	logger    *slog.Logger
}

func NewHandler(service AuthenticationService, validator *validator.Validate, logger *slog.Logger) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *Handler) Register(e echo.Context) error {
	userDTO := new(RegisterRequest)
	if err := e.Bind(userDTO); err != nil {
		return httpx.BindErrorResponse(e)
	}

	if err := h.validator.Struct(userDTO); err != nil {
		validateError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateError)
	}

	input := userDTO.ToRegisterInput()
	user, err := h.service.Register(e.Request().Context(), input)

	if err != nil {
		return httpx.HandlerError(e, err)
	}

	response := NewRegisterResponse(user)

	return e.JSON(http.StatusCreated, response)
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
		TokenType:            "Bearer",
		AccessToken:          loginOutput.AccessToken,
		RefreshToken:         loginOutput.RefreshToken,
		AccessTokenExpiresIn: loginOutput.ExpiresIn,
		RefreshTokenExp:      loginOutput.RefreshExpiresAt,
	}

	return e.JSON(http.StatusOK, response)

}

func (h *Handler) Refresh(e echo.Context) error {
	dto := RefreshTokenRequest{}
	if err := e.Bind(&dto); err != nil {
		return httpx.BindErrorResponse(e)
	}

	if err := h.validator.Struct(&dto); err != nil {
		validateErr := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateErr)
	}

	serviceOutput, err := h.service.RefreshTokens(e.Request().Context(), dto.RefreshToken)

	if err != nil {
		return httpx.HandlerError(e, err)
	}

	if serviceOutput == nil {
		h.logger.Error("Critical: Service returned nil output without error")

		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}

	refreshResponse := LoginResponse{
		TokenType:            "Bearer",
		AccessToken:          serviceOutput.AccessToken,
		RefreshToken:         serviceOutput.RefreshToken,
		AccessTokenExpiresIn: serviceOutput.ExpiresIn,
		RefreshTokenExp:      serviceOutput.RefreshExpiresAt,
	}

	return e.JSON(http.StatusOK, &refreshResponse)
}

func (h *Handler) Delete(e echo.Context) error {
	id, err := httpx.GetIDParam(e, "id")
	if err != nil {
		return httpx.BindErrorResponse(e)
	}

	claim := httpx.GetClaims(e)
	role := claim.Role

	if err := h.service.Delete(e.Request().Context(), role, id); err != nil {
		return httpx.HandlerError(e, err)
	}
	return e.NoContent(http.StatusNoContent)
}
