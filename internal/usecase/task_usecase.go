package usecase

import (
	"context"
	"errors"
	"log"
	"time"
	"to-do-list/internal/model"

	"github.com/sirupsen/logrus"
)

type TaskUsecase struct {
	taskRepo model.ITaskRepository
}

func NewTaskUsecase(taskRepo model.ITaskRepository) model.ITaskUsecase {
	return &TaskUsecase{
		taskRepo: taskRepo,
	}
}

func (t *TaskUsecase) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.Task, error) {
	log := logrus.WithFields(logrus.Fields{
		"filter": filter,
	})

	tasks, err := t.taskRepo.FindAll(ctx, filter)
	if err != nil {
		log.Error("Error fetching tasks: ", err)
		return nil, err
	}

	return tasks, nil
}

func (t *TaskUsecase) FindById(ctx context.Context, id int64) (*model.Task, error) {
	task, err := t.taskRepo.FindById(ctx, id)
	if err != nil {
		log.Printf("Error fetching task by ID: %v", err)
		return nil, err
	}

	if task == nil {
		return nil, errors.New("task not found")
	}

	if task.DeletedAt != nil {
		return nil, errors.New("task is deleted")
	}

	return task, nil
}

func (t *TaskUsecase) Create(ctx context.Context, in model.CreateTaskInput) error {
	log := logrus.WithFields(logrus.Fields{
		"input": in,
	})

	if err := validateCreateTaskInput(in); err != nil {
		return err
	}

	task := model.Task{
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		Priority:    in.Priority,
		DueDate:     in.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := t.taskRepo.Create(ctx, task); err != nil {
		log.Printf("Error creating task: %v", err)
		return err
	}

	return nil
}

func (t *TaskUsecase) Update(ctx context.Context, id int64, in model.UpdateTaskInput) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	existingTask, err := t.taskRepo.FindById(ctx, id)
	if err != nil {
		log.Printf("Error fetching task by ID: %v", err)
		return err
	}

	if existingTask == nil {
		return errors.New("task not found")
	}

	if existingTask.DeletedAt != nil {
		return errors.New("task is deleted and cannot be updated")
	}

	func(task *model.Task, input model.UpdateTaskInput) {
		task.Title = input.Title
		task.Description = input.Description
		task.Status = input.Status
		task.Priority = input.Priority
		task.DueDate = input.DueDate
		task.UpdatedAt = time.Now()
	}(existingTask, in)

	if err := t.taskRepo.Update(ctx, *existingTask); err != nil {
		log.Printf("Error updating task: %v", err)
		return err
	}

	return nil
}

func (t *TaskUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"id": id,
	})

	task, err := t.taskRepo.FindById(ctx, id)
	if err != nil {
		log.Error("Failed to find task for deletion: ", err)
		return err
	}

	if task == nil {
		log.Error("Task not found")
		return errors.New("task not found")
	}

	if task.DeletedAt != nil {
		log.Error("Task already deleted")
		return errors.New("task already deleted")
	}

	err = t.taskRepo.Delete(ctx, id)
	if err != nil {
		log.Error("Failed to delete task: ", err)
		return err
	}

	log.Info("Successfully deleted task with ID: ", id)
	return nil
}

func validateCreateTaskInput(in model.CreateTaskInput) error {
	if in.Title == "" || len(in.Title) < 3 || len(in.Title) > 255 {
		return errors.New("invalid title: must be between 3 and 255 characters")
	}
	if in.Description == "" {
		return errors.New("description is required")
	}
	if in.Status != "pending" && in.Status != "in_progress" && in.Status != "completed" {
		return errors.New("invalid status: must be one of pending, in_progress, or completed")
	}
	if in.Priority != "low" && in.Priority != "medium" && in.Priority != "high" {
		return errors.New("invalid priority: must be one of low, medium, or high")
	}
	return nil
}
