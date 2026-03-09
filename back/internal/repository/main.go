package repository

import (
	"context"
	"errors"
	"time"

	"rri/task-back/internal/models"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	ListByColumn(ctx context.Context, columnID string) ([]models.Task, error)
	Move(ctx context.Context, taskID string, newColumnID string, newPosition int) error
	// Дополнительно: методы для проверки существования и т.п.
	Exists(ctx context.Context, id string) (bool, error)
}

type taskRepository struct {
	db *gorm.DB
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *taskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	// Загружаем связанных пользователей (Assignee и Creator)
	err := r.db.WithContext(ctx).
		Preload("Assignee").
		Preload("Creator").
		Preload("Column").
		Where("id = ?", id).
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound // определим свою ошибку
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	// Используем Updates, чтобы обновить только изменённые поля
	return r.db.WithContext(ctx).Model(task).Updates(task).Error
}

func (r *taskRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Task{ID: id}).Error
}

func (r *taskRepository) ListByColumn(ctx context.Context, columnID string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("column_id = ?", columnID).
		Order("position asc").
		Find(&tasks).Error
	return tasks, err
}

// Move реализует логику перемещения задачи с пересчётом позиций в транзакции
func (r *taskRepository) Move(ctx context.Context, taskID string, newColumnID string, newPosition int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Получаем задачу
		var task models.Task
		if err := tx.First(&task, "id = ?", taskID).Error; err != nil {
			return err
		}

		oldColumnID := task.ColumnID
		oldPosition := task.Position

		// Если колонка не меняется
		if oldColumnID == newColumnID {
			if newPosition == oldPosition {
				return nil // ничего не делаем
			}
			// Сдвиг внутри одной колонки
			if newPosition > oldPosition {
				// сдвиг вниз: задачи между old+1 и new включительно сдвигаются вверх на 1
				if err := tx.Model(&models.Task{}).
					Where("column_id = ? AND position > ? AND position <= ?", oldColumnID, oldPosition, newPosition).
					Update("position", gorm.Expr("position - 1")).Error; err != nil {
					return err
				}
			} else {
				// сдвиг вверх: задачи между new и old-1 сдвигаются вниз на 1
				if err := tx.Model(&models.Task{}).
					Where("column_id = ? AND position >= ? AND position < ?", oldColumnID, newPosition, oldPosition).
					Update("position", gorm.Expr("position + 1")).Error; err != nil {
					return err
				}
			}
		} else {
			// Перемещение в другую колонку
			// 1. Уменьшаем позиции в старой колонке (все, кто были после удаляемой)
			if err := tx.Model(&models.Task{}).
				Where("column_id = ? AND position > ?", oldColumnID, oldPosition).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				return err
			}

			// 2. Увеличиваем позиции в новой колонке (начиная с newPosition)
			if err := tx.Model(&models.Task{}).
				Where("column_id = ? AND position >= ?", newColumnID, newPosition).
				Update("position", gorm.Expr("position + 1")).Error; err != nil {
				return err
			}
		}

		// Обновляем саму задачу
		updates := map[string]interface{}{
			"column_id":  newColumnID,
			"position":   newPosition,
			"updated_at": time.Now(),
		}
		if err := tx.Model(&models.Task{}).Where("id = ?", taskID).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *taskRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Определим кастомную ошибку NotFound
var ErrNotFound = errors.New("record not found")

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

/*
Тут реализуем логику работы с колонками
*/

type ColumnRepository interface {
	Create(ctx context.Context, column *models.Column) error
	GetByID(ctx context.Context, id string) (*models.Column, error)
	//	Update(ctx context.Context, task *models.Task) error
	//	Delete(ctx context.Context, id string) error
	//	ListByColumn(ctx context.Context, columnID string) ([]models.Task, error)
	//	Move(ctx context.Context, taskID string, newColumnID string, newPosition int) error
	// Дополнительно: методы для проверки существования и т.п.
	//	Exists(ctx context.Context, id string) (bool, error)
}

type columnRepository struct {
	db *gorm.DB
}

func NewColumnRepository(db *gorm.DB) ColumnRepository {
	return &columnRepository{db: db}
}
func (r *columnRepository) Create(ctx context.Context, column *models.Column) error {
	return r.db.WithContext(ctx).Create(column).Error
}

func (r *columnRepository) GetByID(ctx context.Context, id string) (*models.Column, error) {
	var column models.Column
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&column).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound // определим свою ошибку
		}
		return nil, err
	}
	return &column, nil
}

/*
Тут реализуем логику работы с людьми, являющимися членами проекта
*/

type ProjectMemberRepository interface {
	Create(ctx context.Context, projectmember *models.ProjectMember) error
	GetUserRole(ctx context.Context, projectID string, userID string) (string, error)
	//	GetByID(ctx context.Context, id string) (*models.Task, error)
	//	Update(ctx context.Context, task *models.Task) error
	//	Delete(ctx context.Context, id string) error
	//	ListByColumn(ctx context.Context, columnID string) ([]models.Task, error)
	//	Move(ctx context.Context, taskID string, newColumnID string, newPosition int) error
	// Дополнительно: методы для проверки существования и т.п.
	//	Exists(ctx context.Context, id string) (bool, error)
}

type projectMemberRepository struct {
	db *gorm.DB
}

func NewProjectMemberRepository(db *gorm.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}
func (r *projectMemberRepository) Create(ctx context.Context, projectmember *models.ProjectMember) error {
	return r.db.WithContext(ctx).Create(projectmember).Error
}

func (r *projectMemberRepository) GetUserRole(ctx context.Context, projectID string, userID string) (string, error) {
	var projectmember models.ProjectMember
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Where("user_id = ?", userID).
		First(&projectmember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound // определим свою ошибку
		}
		return "", err
	}
	return projectmember.Role, nil
}
