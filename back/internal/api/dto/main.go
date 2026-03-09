package dto

import (
	"time"
)

// CreateTaskRequest — тело запроса на создание задачи
type CreateTaskRequest struct {
	ColumnID    string     `json:"column_id" validate:"required,uuid"`
	Title       string     `json:"title" validate:"required,min=1,max=500"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssigneeID  *string    `json:"assignee_id,omitempty" validate:"omitempty,uuid"`
	Labels      []string   `json:"labels,omitempty"`
}

// UpdateTaskRequest — тело запроса на обновление задачи
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=500"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssigneeID  *string    `json:"assignee_id,omitempty" validate:"omitempty,uuid"`
	Labels      []string   `json:"labels,omitempty"`
}

// MoveTaskRequest — тело запроса на перемещение задачи
type MoveTaskRequest struct {
	ColumnID string `json:"column_id" validate:"required,uuid"`
	Position int    `json:"position" validate:"min=0"`
}

// TaskResponse — ответ с данными задачи
type TaskResponse struct {
	ID          string     `json:"id"`
	ColumnID    string     `json:"column_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	Assignee    *UserShort `json:"assignee,omitempty"`
	CreatorID   string     `json:"creator_id"`
	Creator     *UserShort `json:"creator,omitempty"`
	Position    int        `json:"position"`
	Labels      []string   `json:"labels"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// UserShort — краткая информация о пользователе для вложенных объектов
type UserShort struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse — стандартный ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code,omitempty"` // опционально
}
