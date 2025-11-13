package repository

import (
	"context"
	"database/sql"
	"log"
	"todo-api/domain"
)

type sqliteTodoRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewSqliteTodoRepository(db *sql.DB, logger *log.Logger) domain.TodoRepository {
	return &sqliteTodoRepository{db: db, log: logger}
}

func (t *sqliteTodoRepository) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	const q = `
		SELECT id, name, description, completed
		FROM todos
		ORDER BY id;
	`

	todos := []domain.Todo{}

	rows, queryErr := t.db.QueryContext(ctx, q)
	if queryErr != nil {
		log.Printf("[GetTodos] query error: %v", queryErr)
		return nil, domain.InternalError
	}
	defer rows.Close()

	for rows.Next() {
		var completed int
		todo := domain.Todo{}

		if err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &completed); err != nil {
			log.Printf("[GetTodos] scan error: %v", err)
			return nil, domain.InternalError
		}

		todo.Completed = completed == 1
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetTodos] rows iteration error: %v", err)
		return nil, domain.InternalError
	}

	log.Printf("[GetTodos] successfully fetched %d todos", len(todos))
	return todos, nil
}
func (t *sqliteTodoRepository) UpdateTodo(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	const q = `
		UPDATE todos
		SET name = ?, description = ?, completed = ?
		WHERE id = ?;
	`

	completed := 0
	if todo.Completed {
		completed = 1
	}

	result, err := t.db.ExecContext(ctx, q, todo.Name, todo.Description, completed, todo.Id)
	if err != nil {
		t.log.Printf("[UpdateTodo] db exec error (id=%d): %v", todo.Id, err)
		return nil, domain.InternalError
	}

	affectedRow, rowsError := result.RowsAffected()
	if rowsError != nil {
		t.log.Printf("[UpdateTodo] rowsAffected error (id=%d): %v", todo.Id, rowsError)
		return nil, domain.InternalError
	}

	if affectedRow == 0 {
		t.log.Printf("[UpdateTodo] todo not found (id=%d)", todo.Id)
		return nil, domain.ErrTodoNotFound
	}

	t.log.Printf("[UpdateTodo] todo updated successfully (id=%d)", todo.Id)
	return todo, nil
}
func (t *sqliteTodoRepository) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	const q = `
		INSERT INTO todos(name, description, completed) VALUES (?,?,?)
		;
	`
	completed := 0
	if todo.Completed == true {
		completed = 1
	}
	result, err := t.db.ExecContext(ctx, q, todo.Name, todo.Description, completed)
	if err != nil {
		t.log.Printf("[CreateTodo] db exec error %v", err)
		return domain.InternalError
	}
	affectedRows, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		t.log.Printf("[CreateTodo] Affected Rows Err : %v", affectedErr)
		return domain.InternalError
	}
	if affectedRows == 0 {
		t.log.Printf("[CreateTodo] no rows affected on insert (name=%s)", todo.Name)
		return domain.InternalError
	}
	lastID, lastIDErr := result.LastInsertId()
	if lastIDErr != nil {
		t.log.Printf("[CreateTodo] LastInsertId error: %v", lastIDErr)
		return domain.InternalError
	}

	todo.Id = int(lastID)
	t.log.Printf("[CreateTodo] todo created successfully (id=%d, name=%s)", todo.Id, todo.Name)
	return nil
}
func (t *sqliteTodoRepository) DeleteTodo(ctx context.Context, id int) error {
	const q = `DELETE FROM todos WHERE id = ?`

	result, err := t.db.ExecContext(ctx, q, id)
	if err != nil {
		t.log.Printf("[DeleteTodo] db exec error (id=%d): %v", id, err)
		return domain.InternalError
	}

	affectedRows, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		t.log.Printf("[DeleteTodo] rowsAffected error (id=%d): %v", id, affectedErr)
		return domain.InternalError
	}

	if affectedRows == 0 {
		t.log.Printf("[DeleteTodo] todo not found (id=%d)", id)
		return domain.ErrTodoNotFound
	}

	t.log.Printf("[DeleteTodo] todo deleted successfully (id=%d)", id)
	return nil
}
func (t *sqliteTodoRepository) PatchTodo(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
	const q = `UPDATE todos SET (name = ?,description = ? ,completed = ? WHERE id = ?) `
	return todo, nil
}
