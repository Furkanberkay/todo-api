package handlers

import (
	"net/http"
	"todo-api/db"

	"github.com/labstack/echo/v4"
)

func AddTodo(c echo.Context) error {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	var completed int
	if todo.Completed {
		completed = 1
	} else {
		completed = 0
	}

	result, err := db.Conn().Exec("INSERT INTO todos (name,description,completed) VALUES (?,?,?)",
		todo.Name,
		todo.Description,
		completed,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if resultId, resultErr := result.LastInsertId(); resultErr == nil {
		todo.Id = int(resultId)
	}

	return c.JSON(http.StatusCreated, todo)

}
