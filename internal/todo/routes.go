package todo

import "github.com/labstack/echo/v4"

func (h *Handler) Routes(e *echo.Group) {
	e.GET("/todos", h.GetTodos)
	e.POST("/todos", h.CreateTodo)
	e.DELETE("/todos/:id", h.DeleteTodo)
}
