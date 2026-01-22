package todo

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"todoApp3/internal/domain"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&domain.Todo{})
	return db
}

func TestNewRepository(t *testing.T) {
	type args struct {
		db     *gorm.DB
		logger *slog.Logger
	}
	tests := []struct {
		name string
		args args
		want domain.TodoRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRepository(tt.args.db, tt.args.logger), "NewRepository(%v, %v)", tt.args.db, tt.args.logger)
		})
	}
}

func TestRepository_CreateTodo(t *testing.T) {
	db := setupTestDB()

	type args struct {
		ctx  context.Context
		todo *domain.Todo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Succes Create",
			args: args{
				ctx: context.Background(),
				todo: &domain.Todo{
					Name:        "Markete Git",
					Description: "Süt al",
					UserID:      1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db, slog.New(slog.NewTextHandler(nil, nil)))
			err := r.CreateTodo(tt.args.ctx, tt.args.todo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var count int64
				db.Model(&domain.Todo{}).Where("name = ?", "Markete Git").Count(&count)
				assert.Equal(t, int64(1), count, "Veritabanında 1 kayıt olmalı")
			}
		})
	}
}

func TestRepository_DeleteTodo(t *testing.T) {
	db := setupTestDB()

	targetTodo := domain.Todo{
		Model:     gorm.Model{ID: 100},
		Name:      "Sensitive User Data",
		UserID:    2,
		Completed: false,
	}
	db.Create(&targetTodo)

	type args struct {
		ctx    context.Context
		userID uint
		todoID uint
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		shouldBeDeleted bool
	}{
		{
			name: "Error_Unauthorized_Delete_Attempt",
			args: args{
				ctx:    context.Background(),
				userID: 99,
				todoID: 100,
			},
			wantErr:         true,
			shouldBeDeleted: false,
		},
		{
			name: "Success_Authorized_Delete",
			args: args{
				ctx:    context.Background(),
				userID: 2,
				todoID: 100,
			},
			wantErr:         false,
			shouldBeDeleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository(db, slog.New(slog.NewTextHandler(io.Discard, nil)))

			err := repo.DeleteTodo(tt.args.ctx, tt.args.userID, tt.args.todoID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			var checkTodo domain.Todo
			result := db.Unscoped().First(&checkTodo, 100)

			assert.NoError(t, result.Error)

			if tt.shouldBeDeleted {
				assert.NotZero(t, checkTodo.DeletedAt.Time, "record should be soft deleted")
			} else {
				assert.Zero(t, checkTodo.DeletedAt.Time, "record should remain active")
			}
		})
	}
}

func TestRepository_GetTodoByID(t *testing.T) {
	db := setupTestDB()

	todo := domain.Todo{
		Model:  gorm.Model{ID: 1},
		UserID: 1,
		Name:   "test",
	}
	db.Create(&todo)
	type args struct {
		ctx    context.Context
		userID uint
		todoID uint
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Todo
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				todoID: 1,
			},
			wantErr: false,
			want: &domain.Todo{
				UserID: 1,
				Name:   "test",
				Model:  gorm.Model{ID: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db, new(slog.Logger))
			got, err := r.GetTodoByID(tt.args.ctx, tt.args.userID, tt.args.todoID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)

				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.UserID, got.UserID)
			}
		})
	}
}

func TestRepository_GetTodos(t *testing.T) {
	db := setupTestDB()

	for i := 1; i <= 15; i++ {
		db.Create(&domain.Todo{
			Name:   fmt.Sprintf("User1 Task %d", i),
			UserID: 1,
		})
	}

	for i := 1; i <= 5; i++ {
		db.Create(&domain.Todo{
			Name:   fmt.Sprintf("User2 Task %d", i),
			UserID: 2,
		})
	}

	type args struct {
		ctx    context.Context
		page   int
		limit  int
		userID uint
	}
	tests := []struct {
		name          string
		args          args
		wantLen       int
		wantTotalItem int64
		wantErr       bool
	}{
		{
			name: "User1_Page1_Limit10",
			args: args{
				ctx:    context.Background(),
				page:   1,
				limit:  10,
				userID: 1,
			},
			wantLen:       10,
			wantTotalItem: 15,
			wantErr:       false,
		},
		{
			name: "User1_Page2_Limit10",
			args: args{
				ctx:    context.Background(),
				page:   2,
				limit:  10,
				userID: 1,
			},
			wantLen:       5,
			wantTotalItem: 15,
			wantErr:       false,
		},
		{
			name: "User2_GetAll",
			args: args{
				ctx:    context.Background(),
				page:   1,
				limit:  100,
				userID: 2,
			},
			wantLen:       5,
			wantTotalItem: 5,
			wantErr:       false,
		},
		{
			name: "User99_NoData",
			args: args{
				ctx:    context.Background(),
				page:   1,
				limit:  10,
				userID: 99,
			},
			wantLen:       0,
			wantTotalItem: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db, slog.New(slog.NewTextHandler(io.Discard, nil)))

			got, totalCount, err := r.GetTodos(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Len(t, got, tt.wantLen, "The length of the rotating list is incorrect.")

				assert.Equal(t, tt.wantTotalItem, int64(totalCount), "Total number of records is incorrect")

				for _, todo := range got {
					assert.Equal(t, tt.args.userID, todo.UserID, "Someone else's data has been mixed up!")
				}
			}
		})
	}
}

func TestRepository_UpdateTodo(t *testing.T) {
	db := setupTestDB()

	originalTodo := domain.Todo{
		Model:       gorm.Model{ID: 1},
		Name:        "Eski İsim",
		Description: "Eski Açıklama",
		Completed:   false,
		UserID:      10,
	}
	db.Create(&originalTodo)

	type args struct {
		ctx  context.Context
		todo *domain.Todo
	}
	tests := []struct {
		name         string
		args         args
		verifyUpdate func(t *testing.T)
		wantErr      bool
	}{
		{
			name: "Success_Update_Name_And_Status",
			args: args{
				ctx: context.Background(),
				todo: &domain.Todo{
					Model:       gorm.Model{ID: 1},
					Name:        "Yeni İsim",
					Description: "Eski Açıklama",
					Completed:   true,
					UserID:      10,
				},
			},
			wantErr: false,
			verifyUpdate: func(t *testing.T) {
				var checkTodo domain.Todo
				db.First(&checkTodo, 1)

				assert.Equal(t, "Yeni İsim", checkTodo.Name, "Name should have been updated")

				assert.Equal(t, true, checkTodo.Completed, "Completed should be true")

				assert.Equal(t, "Eski Açıklama", checkTodo.Description, "Description should not have changed")

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRepository(db, slog.New(slog.NewTextHandler(io.Discard, nil)))

			err := r.UpdateTodo(tt.args.ctx, tt.args.todo)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.verifyUpdate != nil {
					tt.verifyUpdate(t)
				}
			}
		})
	}
}
