package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"todo-api/db"

	"github.com/labstack/echo/v4"
)

func UpdateTodos(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError("id must be a number"))
	}
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError("invalid body"))
	}

	if strings.TrimSpace(todo.Name) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"name is required"})
	}
	if strings.TrimSpace(todo.Description) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"description is required"})
	}

	result, err := db.Conn().Exec(
		"UPDATE todos SET name = ?, description = ?, completed = ? WHERE id = ?",
		todo.Name,
		todo.Description,
		todo.Completed,
		id,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}
	resultInt, resultErr := result.RowsAffected()
	if resultErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(resultErr))
	}
	if resultInt == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{"todo not found"})
	}
	todo.Id = id
	return c.JSON(http.StatusOK, todo)

}
