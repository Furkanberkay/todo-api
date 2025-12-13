package auth

import "github.com/labstack/echo/v4"

func (h *Handler) Routes(echo *echo.Group) {
	echo.POST("/register", h.RegisterUser)
	echo.POST("/login", h.Login)
	echo.POST("/refresh", h.Refresh)
}
