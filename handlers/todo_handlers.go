package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"todo-api/db"
	"todo-api/models"

	"github.com/labstack/echo/v4"
)

func GetTodos(c echo.Context) error {
	var todos []models.Todo

	rows, err := db.Conn().Query("SELECT * FROM todos")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
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
func UpdateTodos(c echo.Context) error {
	idStr := c.Param("id")

	id, paramErr := strconv.Atoi(idStr)
	if paramErr != nil {
		return c.JSON(http.StatusBadRequest, ValidationError("id must be a number"))
	}
	var updatedTodo UpdateTodoRequest
	if err := c.Bind(&updatedTodo); err != nil {
		return c.JSON(http.StatusBadRequest, ValidationError("invalid body"))
	}

	if strings.TrimSpace(updatedTodo.Name) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"name is required"})
	}
	if strings.TrimSpace(updatedTodo.Description) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"description is required"})
	}
	if updatedTodo.Completed == nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"completed is required"})
	}
	completed := boolToInt(*updatedTodo.Completed)

	result, execErr := db.Conn().Exec(
		"UPDATE todos SET name = ?, description = ?, completed = ? WHERE id = ?",
		updatedTodo.Name,
		updatedTodo.Description,
		completed,
		id,
	)
	if execErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(execErr))
	}
	resultInt, resultErr := result.RowsAffected()
	if resultErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(resultErr))
	}
	if resultInt == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{"todo not found"})
	}

	todo := models.Todo{
		Id:          id,
		Name:        updatedTodo.Name,
		Description: updatedTodo.Description,
		Completed:   *updatedTodo.Completed,
	}
	return c.JSON(http.StatusOK, todo)

}
func CreateTodo(c echo.Context) error {
	var todo models.Todo
	if err := c.Bind(&todo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if strings.TrimSpace(todo.Name) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"name is required"})
	}
	if strings.TrimSpace(todo.Description) == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"description is required"})
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
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}

	if resultId, resultErr := result.LastInsertId(); resultErr == nil {
		todo.Id = int(resultId)
	}

	return c.JSON(http.StatusCreated, todo)

}
