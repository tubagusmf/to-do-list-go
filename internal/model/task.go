package model

import (
	"context"
	"time"
)

type ITaskRepository interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*Task, error)
	FindById(ctx context.Context, id int64) (*Task, error)
	Create(ctx context.Context, task Task) error
	Update(ctx context.Context, task Task) error
	Delete(ctx context.Context, id int64) error
}

type ITaskUsecase interface {
	FindAll(ctx context.Context, filter FindAllParam) ([]*Task, error)
	FindById(ctx context.Context, id int64) (*Task, error)
	Create(ctx context.Context, in CreateTaskInput) error
	Update(ctx context.Context, id int64, in UpdateTaskInput) error
	Delete(ctx context.Context, id int64) error
}

type Task struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-"`
}

type FindAllParam struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type CreateTaskInput struct {
	Title       string     `json:"title" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"required"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date" validate:"required"`
}

type UpdateTaskInput struct {
	Title       string     `json:"title" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"required"`
	Status      string     `json:"status" validate:"required,oneof=pending in_progress completed"`
	Priority    string     `json:"priority" validate:"required,oneof=low medium high"`
	DueDate     *time.Time `json:"due_date" validate:"required"`
}
