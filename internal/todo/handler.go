package todo

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"todoApp3/internal/httpx"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service   *Service
	logger    *slog.Logger
	validator *validator.Validate
}

func NewTodoHandler(service *Service, logger *slog.Logger, validator *validator.Validate) *Handler {
	return &Handler{
		service:   service,
		logger:    logger,
		validator: validator,
	}
}

func (h *Handler) GetTodos(e echo.Context) error {

	pageStr := e.QueryParam("page")
	page, pageErr := strconv.Atoi(pageStr)
	if pageErr != nil || page < 1 {
		page = 1
	}

	limitStr := e.QueryParam("limit")
	limit, limitErr := strconv.Atoi(limitStr)
	if limitErr != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	userID, ok := e.Get("userID").(uint)
	if !ok {
		h.logger.Error("userID not found in context")
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	todos, totalCount, err := h.service.GetTodos(e.Request().Context(), page, limit, userID)
	if err != nil {
		return httpx.HandlerError(e, err)
	}

	todoResponse := toTodoResponseList(todos)

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	response := PaginatedResponse{
		Data: todoResponse,
		Meta: PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalItems:  int64(totalCount),
			TotalPages:  totalPages,
		},
	}

	return e.JSON(http.StatusOK, response)
}

func (h *Handler) CreateTodo(e echo.Context) error {
	userID, ok := e.Get("userID").(uint)
	if !ok {
		h.logger.Error("userID not found in context")
		return e.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	dto := CreateTodoRequest{}

	if err := e.Bind(&dto); err != nil {
		return httpx.BindErrorResponse(e, err)
	}

	if err := h.validator.Struct(&dto); err != nil {
		validateError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateError)
	}

	createTodoInput := CreateTodoInput{
		Name:        dto.Name,
		Description: dto.Description,
	}

	todo, err := h.service.CreateTodo(e.Request().Context(), &createTodoInput, userID)

	if err != nil {
		return httpx.HandlerError(e, err)
	}

	detailResponse := TodoResponse{
		Name:        todo.Name,
		Description: todo.Description,
		ID:          todo.ID,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
		Version:     todo.Version.Int64,
	}

	return e.JSON(http.StatusCreated, detailResponse)

}

func (h *Handler) DeleteTodo(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id < 1 {
		return e.JSON(http.StatusBadRequest, httpx.ResponseError{
			Message: "Invalid ID format. ID must be a positive number.",
		})
	}

	userID, ok := e.Get("userID").(uint)

	if !ok {
		return e.JSON(http.StatusUnauthorized, httpx.ResponseError{Message: "Unauthorized"})
	}

	if err := h.service.DeleteTodo(e.Request().Context(), userID, uint(id)); err != nil {
		return httpx.HandlerError(e, err)
	}

	return e.NoContent(http.StatusNoContent)

}

func (h *Handler) UpdateTodo(e echo.Context) error {
	idStr := e.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		return e.JSON(http.StatusBadRequest, httpx.ResponseError{
			Message: "id must be positive number",
		})
	}

	userID, ok := e.Get("userID").(uint)

	if !ok {
		return e.JSON(http.StatusUnauthorized, httpx.ResponseError{
			Message: "UnAuthorized",
		})
	}

	updateDTO := new(UpdateRequest)

	if err := e.Bind(updateDTO); err != nil {
		return httpx.BindErrorResponse(e, err)
	}

	if err := h.validator.Struct(updateDTO); err != nil {
		validateError := httpx.ParseValidationErrors(err)
		return e.JSON(http.StatusBadRequest, validateError)
	}

	updateInput := UpdateRequestToUpdateInput(id, updateDTO)

	todo, serviceErr := h.service.UpdateTodo(e.Request().Context(), updateInput, userID)

	if serviceErr != nil {
		return httpx.HandlerError(e, serviceErr)
	}

	todoResponse := TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		Completed:   todo.Completed,
		Version:     todo.Version.Int64,
	}

	return e.JSON(http.StatusOK, todoResponse)

}
