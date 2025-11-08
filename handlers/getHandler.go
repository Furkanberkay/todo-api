package handlers

import (
	"net/http"
	"todo-api/db"

	"github.com/labstack/echo/v4"
)

func GetTodos(c echo.Context) error {
	var todos []Todo

	rows, err := db.Conn().Query("SELECT * FROM todos")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		var completed int

		err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &completed)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		}

		if completed == 1 {
			todo.Completed = true
		} else {
			todo.Completed = false
		}

		todos = append(todos, todo)
	}

	return c.JSON(http.StatusOK, todos)
}
