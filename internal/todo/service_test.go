package todo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"todoApp3/internal/domain"
	"todoApp3/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	type fields struct {
		repository domain.TodoRepository
		logger     *slog.Logger
	}
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Todo
		want1   int
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
				logger:     tt.fields.logger,
			}
			got, got1, err := s.GetTodos(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetTodos(%v, %v, %v, %v)", tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetTodos(%v, %v, %v, %v)", tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)
			assert.Equalf(t, tt.want1, got1, "GetTodos(%v, %v, %v, %v)", tt.args.ctx, tt.args.page, tt.args.limit, tt.args.userID)
		})
	}
}

func TestService_PatchTodo(t *testing.T) {
	type fields struct {
		repository domain.TodoRepository
		logger     *slog.Logger
	}
	type args struct {
		ctx    context.Context
		userID uint
		input  *PatchTodoInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Todo
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
				logger:     tt.fields.logger,
			}
			got, err := s.PatchTodo(tt.args.ctx, tt.args.userID, tt.args.input)
			if !tt.wantErr(t, err, fmt.Sprintf("PatchTodo(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.input)) {
				return
			}
			assert.Equalf(t, tt.want, got, "PatchTodo(%v, %v, %v)", tt.args.ctx, tt.args.userID, tt.args.input)
		})
	}
}

func TestService_UpdateTodo1(t *testing.T) {
	type fields struct {
		repository domain.TodoRepository
		logger     *slog.Logger
	}
	type args struct {
		ctx    context.Context
		input  *UpdateTodoInput
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Todo
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				repository: tt.fields.repository,
				logger:     tt.fields.logger,
			}
			got, err := s.UpdateTodo(tt.args.ctx, tt.args.input, tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateTodo(%v, %v, %v)", tt.args.ctx, tt.args.input, tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateTodo(%v, %v, %v)", tt.args.ctx, tt.args.input, tt.args.userID)
		})
	}
}
