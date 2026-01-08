package auth

import "github.com/labstack/echo/v4"

func (h *Handler) Routes(echo *echo.Group) {
	echo.POST("/register", h.Register)
	echo.POST("/login", h.Login)
	echo.POST("/refresh", h.Refresh)
}

func (h *Handler) ProtectedRoutes(echo *echo.Group) {
	echo.DELETE("/users/:id", h.Delete)
}
