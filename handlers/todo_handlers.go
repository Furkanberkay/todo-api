package handlers

import (
	"database/sql"
	"log"
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
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows close error: %v", cerr)
		}
	}()

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

	resultId, resultErr := result.LastInsertId()
	if resultErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(resultErr))
	}
	todo.Id = int(resultId)

	return c.JSON(http.StatusCreated, todo)

}
func DeleteTodo(c echo.Context) error {
	paramIdStr := c.Param("id")
	paramId, err := strconv.Atoi(paramIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"id must be a number"})
	}
	resultDb, errorDb := db.Conn().Exec("DELETE FROM todos WHERE id = ?", paramId)
	if errorDb != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": errorDb.Error(),
		})
	}
	rowsAffected, rowsErr := resultDb.RowsAffected()
	if rowsErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}
	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, ErrorResponse{"todo not found"})
	}
	return c.NoContent(http.StatusNoContent)
}

func PatchTodo(c echo.Context) error {

	strId := c.Param("id")
	id, errId := strconv.Atoi(strId)
	if errId != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"id must be a number"})
	}
	var patchTodo PatchTodoRequest

	if err := c.Bind(&patchTodo); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"invalid body"})
	}
	todo := models.Todo{}
	var completed int

	row := db.Conn().QueryRow("SELECT name,description,completed FROM todos WHERE id = ?", id)
	if err := row.Scan(&todo.Name, &todo.Description, &completed); err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, ErrorResponse{"todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}
	if completed == 1 {
		todo.Completed = true
	} else {
		todo.Completed = false
	}

	if patchTodo.Name != nil {
		if strings.TrimSpace(*patchTodo.Name) == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"name is required"})
		}
		todo.Name = *patchTodo.Name
	}
	if patchTodo.Description != nil {
		if strings.TrimSpace(*patchTodo.Description) == "" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{"description is required"})
		}
		todo.Description = *patchTodo.Description
	}
	if patchTodo.Completed != nil {
		if *patchTodo.Completed {
			todo.Completed = true
		} else {
			todo.Completed = false
		}
	}

	completedDB := boolToInt(todo.Completed)
	result, err := db.Conn().Exec("UPDATE todos SET name = ?, description = ?, completed = ? WHERE id = ?",
		todo.Name,
		todo.Description,
		completedDB,
		id,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}
	affected, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		return c.JSON(http.StatusInternalServerError, InternalError(affectedErr))
	}
	if affected == 0 {
		log.Printf("[PATCH /todos/%d] No changes applied — existing values identical.\n", id)
	}

	return c.JSON(http.StatusOK, todo)

}

func GetTodosById(c echo.Context) error {
	strId := c.Param("id")
	id, strErr := strconv.Atoi(strId)
	if strErr != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"id must be a number"})
	}
	rows := db.Conn().QueryRow("SELECT id, name, description, completed FROM todos WHERE id = ?", id)

	var completed int
	todo := models.Todo{}

	err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, ErrorResponse{"todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, InternalError(err))
	}
	if completed == 1 {
		todo.Completed = true
	} else {
		todo.Completed = false
	}

	return c.JSON(http.StatusOK, todo)

}
