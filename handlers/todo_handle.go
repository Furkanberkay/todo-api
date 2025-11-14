package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"todo-api/domain"
	"todo-api/service"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}
func (h *TodoHandler) ClientRouters(e *echo.Echo) {
	e.GET("/todos", h.GetTodos)
	e.GET("/todos/:id", h.GetTodoByID)

}

func (h *TodoHandler) GetTodos(e echo.Context) error {
	todos, err := h.todoService.GetTodos(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusInternalServerError, ResponseErr{Message: err.Error()})
	}
	return e.JSON(http.StatusOK, todos)
}
func (h *TodoHandler) GetTodoByID(e echo.Context) error {
	pathId := e.Param("id")
	id, err := strconv.Atoi(pathId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "",
		})
	}
	todo, serviceErr := h.todoService.GetTodoByID(e.Request().Context(), id)
	if serviceErr != nil {
		if errors.Is(serviceErr, domain.ErrTodoNotFound) {
			return e.JSON(http.StatusNotFound, ResponseErr{Message: serviceErr.Error()})
		}
		return e.JSON(http.StatusInternalServerError, ResponseErr{Message: serviceErr.Error()})
	}
	return e.JSON(http.StatusOK, todo)
}
func (h *TodoHandler) DeleteTodo(e echo.Context) error {
	deleteId := e.Param("id")
	id, err := strconv.Atoi(deleteId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, ResponseErr{
			Message: "id must be a number",
		})
	}
	if err := h.todoService.DeleteTodo(e.Request().Context(), id); err != nil {
		if errors.Is(err, domain.ErrTodoNotFound) {
			return e.JSON(http.StatusNotFound, ResponseErr{Message: err.Error()})
		}
		return e.JSON(http.StatusInternalServerError, ResponseErr{Message: err.Error()})
	}
	return e.NoContent(http.StatusNoContent)
}
