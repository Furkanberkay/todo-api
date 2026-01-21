package todo

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"todoApp3/internal/domain"
	"todoApp3/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestService_UpdateTodo(t *testing.T) {
	type args struct {
		ctx    context.Context
		input  *UpdateTodoInput
		userID uint
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks.TodoRepository)
		want      *domain.Todo
		wantErr   bool
	}{
		{
			name: "Basarili Guncelleme",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: &UpdateTodoInput{
					ID:          1,
					Name:        "Spor Yap",
					Description: "30 dakika koşu",
					Completed:   true,
				},
			},
			setupMock: func(m *mocks.TodoRepository) {
				existingTodo := &domain.Todo{
					UserID:      1,
					Name:        "Eski İsim",
					Description: "Eski Açıklama",
					Completed:   false,
				}
				m.On("GetTodoByID", mock.Anything, uint(1), uint(1)).Return(existingTodo, nil)
				m.On("UpdateTodo", mock.Anything, mock.MatchedBy(func(d *domain.Todo) bool {
					return d.Name == "Spor Yap" && d.Completed == true
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Todo Bulunamadi (404)",
			args: args{
				ctx:    context.Background(),
				userID: 99,
				input: &UpdateTodoInput{
					ID:   99,
					Name: "Hayalet Görev",
				},
			},
			setupMock: func(m *mocks.TodoRepository) {
				m.On("GetTodoByID", mock.Anything, uint(99), uint(99)).Return(nil, errors.New("todo not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			s := &Service{
				repository: mockRepo,
				logger:     logger,
			}

			got, err := s.UpdateTodo(tt.args.ctx, tt.args.input, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if got != nil {
					assert.Equal(t, tt.args.input.Name, got.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewTodoService(t *testing.T) {

	mockRepo := new(mocks.TodoRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	type args struct {
		repository domain.TodoRepository
		logger     *slog.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success created service",
			args: args{
				repository: mockRepo,
				logger:     logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTodoService(tt.args.repository, logger)
			assert.NotNil(t, got, "services should not nil")
			assert.Equalf(t, tt.args.repository, got.repository, "repository dogru atanmalı")
			assert.Equalf(t, tt.args.logger, got.logger, "logger dogru atanmalı")
		})
	}
}

func TestService_CreateTodo(t *testing.T) {

	type args struct {
		ctx    context.Context
		input  *CreateTodoInput
		userID uint
	}
	tests := []struct {
		name      string
		args      args
		setUpMock func(m *mocks.TodoRepository)
		want      *domain.Todo
		wantErr   bool
	}{
		{
			name: "success created",
			args: args{
				ctx:    context.Background(),
				userID: 55,
				input: &CreateTodoInput{
					Name:        "test_deneme_123",
					Description: "deneme",
				},
			},
			setUpMock: func(m *mocks.TodoRepository) {
				m.On("CreateTodo", mock.Anything, mock.MatchedBy(func(d *domain.Todo) bool {
					return d.Name == "test_deneme_123" &&
						d.Description == "deneme" &&
						d.UserID == 55
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Veritabani Hatasi",
			args: args{
				ctx:    context.Background(),
				userID: 55,
				input: &CreateTodoInput{
					Name: "Hatalı İşlem",
				},
			},
			setUpMock: func(m *mocks.TodoRepository) {
				m.On("CreateTodo", mock.Anything, mock.MatchedBy(func(_ *domain.Todo) bool {
					return true
				})).Return(errors.New("db connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)

			if tt.setUpMock != nil {
				tt.setUpMock(mockRepo)
			}

			s := &Service{
				repository: mockRepo,
				logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}

			got, err := s.CreateTodo(tt.args.ctx, tt.args.input, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if got != nil {
					assert.Equal(t, tt.args.input.Name, got.Name)
					assert.Equal(t, tt.args.userID, got.UserID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_DeleteTodo(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID uint
		todoID uint
	}
	tests := []struct {
		name      string
		setUpMock func(m *mocks.TodoRepository)
		args      args
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userID: 44,
				todoID: 12,
			},
			setUpMock: func(m *mocks.TodoRepository) {
				m.On("DeleteTodo", mock.Anything, uint(44), uint(12)).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "database error",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				todoID: 99,
			},
			setUpMock: func(m *mocks.TodoRepository) {
				m.On("DeleteTodo", mock.Anything, uint(1), uint(99)).Return(errors.New("db delete error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)

			if tt.setUpMock != nil {
				tt.setUpMock(mockRepo)
			}
			s := &Service{
				repository: mockRepo,
				logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}
			err := s.DeleteTodo(tt.args.ctx, tt.args.userID, tt.args.todoID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTodoByID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID uint
		todoID uint
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks.TodoRepository)
		want      *domain.Todo
		wantErr   bool
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				todoID: 1,
			},
			setupMock: func(m *mocks.TodoRepository) {
				expectedTodo := &domain.Todo{
					UserID:      1,
					Name:        "berk",
					Description: "ssds",
					Completed:   false,
				}
				m.On("GetTodoByID", mock.Anything, uint(1), uint(1)).Return(expectedTodo, nil)
			},
			want: &domain.Todo{
				UserID:      1,
				Name:        "berk",
				Description: "ssds",
				Completed:   false,
			},
			wantErr: false,
		},
		{
			name: "Database Error",
			args: args{
				ctx:    context.Background(),
				userID: 2,
				todoID: 2,
			},
			setupMock: func(m *mocks.TodoRepository) {
				m.On("GetTodoByID", mock.Anything, uint(2), uint(2)).Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := &Service{
				repository: mockRepo,
				logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}

			got, err := s.GetTodoByID(tt.args.ctx, tt.args.userID, tt.args.todoID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetTodos(t *testing.T) {
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		userID uint
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks.TodoRepository)
		want      []domain.Todo
		wantCount int
		wantErr   bool
	}{
		{
			name: "Success_Standard",
			args: args{
				ctx:    context.Background(),
				page:   1,
				limit:  10,
				userID: 55,
			},
			setupMock: func(m *mocks.TodoRepository) {
				fakeTodos := []domain.Todo{
					{Model: gorm.Model{ID: 1}, Name: "Task 1", UserID: 55},
					{Model: gorm.Model{ID: 2}, Name: "Task 2", UserID: 55},
				}
				m.On("GetTodos", mock.Anything, 1, 10, uint(55)).Return(fakeTodos, 2, nil)
			},
			want: []domain.Todo{
				{Model: gorm.Model{ID: 1}, Name: "Task 1", UserID: 55},
				{Model: gorm.Model{ID: 2}, Name: "Task 2", UserID: 55},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "Success_Pagination_Logic",
			args: args{
				ctx:    context.Background(),
				page:   0,
				limit:  100,
				userID: 55,
			},
			setupMock: func(m *mocks.TodoRepository) {
				m.On("GetTodos", mock.Anything, 1, 30, uint(55)).Return([]domain.Todo{}, 0, nil)
			},
			want:      []domain.Todo{},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "Database_Error",
			args: args{
				ctx:    context.Background(),
				page:   1,
				limit:  10,
				userID: 55,
			},
			setupMock: func(m *mocks.TodoRepository) {
				m.On("GetTodos", mock.Anything, 1, 10, uint(55)).Return(nil, 0, errors.New("db error"))
			},
			want:      nil,
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := &Service{
				repository: mockRepo,
				logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}

			got, count, err := s.GetTodos(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.wantCount, count)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_PatchTodo(t *testing.T) {
	updatedName := "Patched Name"

	type args struct {
		ctx    context.Context
		userID uint
		input  *PatchTodoInput
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(m *mocks.TodoRepository)
		want      *domain.Todo
		wantErr   bool
	}{
		{
			name: "Success_Patch_Name_Only",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: &PatchTodoInput{
					ID:   10,
					Name: &updatedName,
				},
			},
			setupMock: func(m *mocks.TodoRepository) {
				oldTodo := &domain.Todo{Model: gorm.Model{ID: 10}, UserID: 1, Name: "Old Name", Description: "Old Description"}

				m.On("GetTodoByID", mock.Anything, uint(1), uint(10)).Return(oldTodo, nil)

				m.On("UpdateTodo", mock.Anything, mock.MatchedBy(func(d *domain.Todo) bool {
					return d.Name == "Patched Name" && d.Description == "Old Description"
				})).Return(nil)
			},
			want: &domain.Todo{
				Model:       gorm.Model{ID: 10},
				UserID:      1,
				Name:        "Patched Name",
				Description: "Old Description",
			},
			wantErr: false,
		},
		{
			name: "Error_All_Fields_Nil",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: &PatchTodoInput{
					ID: 10,
				},
			},
			setupMock: func(m *mocks.TodoRepository) {
				
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error_Todo_Not_Found",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				input: &PatchTodoInput{
					ID:   99,
					Name: &updatedName,
				},
			},
			setupMock: func(m *mocks.TodoRepository) {
				m.On("GetTodoByID", mock.Anything, uint(1), uint(99)).Return(nil, errors.New("not found"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.TodoRepository)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := &Service{
				repository: mockRepo,
				logger:     slog.New(slog.NewTextHandler(os.Stdout, nil)),
			}

			got, err := s.PatchTodo(tt.args.ctx, tt.args.userID, tt.args.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
