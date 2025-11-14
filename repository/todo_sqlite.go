package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"todo-api/domain"
)

type SqliteTodoRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewSqliteTodoRepository(db *sql.DB, logger *log.Logger) domain.TodoRepository {
	return &SqliteTodoRepository{db: db, log: logger}
}

func (t *SqliteTodoRepository) GetTodos(ctx context.Context) ([]domain.Todo, error) {
	const q = `
		SELECT id, name, description, completed
		FROM todos
		ORDER BY id;
	`

	todos := []domain.Todo{}

	rows, queryErr := t.db.QueryContext(ctx, q)
	if queryErr != nil {
		log.Printf("[GetTodos] query error: %v", queryErr)
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var completed int
		todo := domain.Todo{}

		if err := rows.Scan(&todo.Id, &todo.Name, &todo.Description, &completed); err != nil {
			log.Printf("[GetTodos] scan error: %v", err)
			return nil, domain.ErrInternal
		}

		todo.Completed = completed == 1
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[GetTodos] rows iteration error: %v", err)
		return nil, domain.ErrInternal
	}

	log.Printf("[GetTodos] successfully fetched %d todos", len(todos))
	return todos, nil
}
func (t *SqliteTodoRepository) UpdateTodo(ctx context.Context, todo *domain.Todo) (*domain.Todo, error) {
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
		return nil, domain.ErrInternal
	}

	affectedRow, rowsError := result.RowsAffected()
	if rowsError != nil {
		t.log.Printf("[UpdateTodo] rowsAffected error (id=%d): %v", todo.Id, rowsError)
		return nil, domain.ErrInternal
	}

	if affectedRow == 0 {
		t.log.Printf("[UpdateTodo] todo not found (id=%d)", todo.Id)
		return nil, domain.ErrTodoNotFound
	}

	t.log.Printf("[UpdateTodo] todo updated successfully (id=%d)", todo.Id)
	return todo, nil
}
func (t *SqliteTodoRepository) CreateTodo(ctx context.Context, todo *domain.Todo) error {
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
		return domain.ErrInternal
	}
	affectedRows, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		t.log.Printf("[CreateTodo] Affected Rows Err : %v", affectedErr)
		return domain.ErrInternal
	}
	if affectedRows == 0 {
		t.log.Printf("[CreateTodo] no rows affected on insert (name=%s)", todo.Name)
		return domain.ErrInternal
	}
	lastID, lastIDErr := result.LastInsertId()
	if lastIDErr != nil {
		t.log.Printf("[CreateTodo] LastInsertId error: %v", lastIDErr)
		return domain.ErrInternal
	}

	todo.Id = int(lastID)
	t.log.Printf("[CreateTodo] todo created successfully (id=%d, name=%s)", todo.Id, todo.Name)
	return nil
}
func (t *SqliteTodoRepository) DeleteTodo(ctx context.Context, id int) error {
	const q = `DELETE FROM todos WHERE id = ?`

	result, err := t.db.ExecContext(ctx, q, id)
	if err != nil {
		t.log.Printf("[DeleteTodo] db exec error (id=%d): %v", id, err)
		return domain.ErrInternal
	}

	affectedRows, affectedErr := result.RowsAffected()
	if affectedErr != nil {
		t.log.Printf("[DeleteTodo] rowsAffected error (id=%d): %v", id, affectedErr)
		return domain.ErrInternal
	}

	if affectedRows == 0 {
		t.log.Printf("[DeleteTodo] todo not found (id=%d)", id)
		return domain.ErrTodoNotFound
	}

	t.log.Printf("[DeleteTodo] todo deleted successfully (id=%d)", id)
	return nil
}
func (t *SqliteTodoRepository) GetTodoByID(ctx context.Context, id int) (*domain.Todo, error) {
	const q = `SELECT id,name,description,completed FROM todos WHERE id = ? `
	row := t.db.QueryRowContext(ctx, q, id)
	todo := domain.Todo{}
	var completed int

	if err := row.Scan(&todo.Id, &todo.Name, &todo.Description, &completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.log.Printf("[GetTodoByID] todo not found (id : %v),%v", id, err)
			return nil, domain.ErrTodoNotFound
		}
		t.log.Printf("[GetTodoByID] row scan error : %v", err)
		return nil, domain.ErrInternal
	}
	if completed == 1 {
		todo.Completed = true
	} else {
		todo.Completed = false
	}
	return &todo, nil
}
