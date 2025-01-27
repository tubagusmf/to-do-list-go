package repository

import (
	"context"
	"time"
	"to-do-list/internal/model"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) model.ITaskRepository {
	return &TaskRepo{
		db: db,
	}
}

func (t *TaskRepo) FindAll(ctx context.Context, filter model.FindAllParam) ([]*model.Task, error) {
	var tasks []*model.Task
	query := t.db.WithContext(ctx).Model(&model.Task{})

	if filter.Limit > 0 {
		query = query.Limit(int(filter.Limit))
	}
	if filter.Page > 0 {
		offset := int((filter.Page - 1) * filter.Limit)
		query = query.Offset(offset)
	}

	err := query.Where("deleted_at IS NULL").Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TaskRepo) FindById(ctx context.Context, id int64) (*model.Task, error) {
	var task model.Task
	err := t.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&task).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

func (t *TaskRepo) Create(ctx context.Context, task model.Task) error {
	err := t.db.WithContext(ctx).Create(&task).Error

	if err != nil {
		return err
	}

	return nil
}

func (t *TaskRepo) Update(ctx context.Context, task model.Task) error {
	err := t.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ? AND deleted_at IS NULL", task.Id).
		Updates(map[string]interface{}{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"due_date":    task.DueDate,
			"updated_at":  task.UpdatedAt,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (t *TaskRepo) Delete(ctx context.Context, id int64) error {
	err := t.db.WithContext(ctx).
		Model(&model.Task{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return err
	}

	return nil
}
