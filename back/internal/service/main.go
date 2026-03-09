package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"rri/task-back/internal/api/dto"
	"rri/task-back/internal/cache"
	"rri/task-back/internal/models"
	"rri/task-back/internal/repository"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type TaskService interface {
	CreateTask(ctx context.Context, userID string, req *dto.CreateTaskRequest) (*dto.TaskResponse, error)
	GetTask(ctx context.Context, taskID string, userID string) (*dto.TaskResponse, error)
	UpdateTask(ctx context.Context, taskID string, userID string, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error)
	DeleteTask(ctx context.Context, taskID string, userID string) error
	MoveTask(ctx context.Context, taskID string, userID string, req *dto.MoveTaskRequest) (*dto.TaskResponse, error)
}

type taskService struct {
	taskRepo    repository.TaskRepository
	columnRepo  repository.ColumnRepository        // для проверки существования колонки
	memberRepo  repository.ProjectMemberRepository // для проверки прав доступа
	cacheClient *cache.Client
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	columnRepo repository.ColumnRepository,
	memberRepo repository.ProjectMemberRepository,
	cacheClient *cache.Client,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		columnRepo:  columnRepo,
		memberRepo:  memberRepo,
		cacheClient: cacheClient,
	}
}

// checkProjectAccess проверяет, имеет ли пользователь доступ к проекту через колонку/доску
func (s *taskService) checkProjectAccess(ctx context.Context, columnID, userID string, requiredRole string) error {
	// Получаем проект через колонку
	column, err := s.columnRepo.GetByID(ctx, columnID)
	if err != nil {
		return fmt.Errorf("column not found: %w", err)
	}
	// Получаем роль пользователя в проекте
	role, err := s.memberRepo.GetUserRole(ctx, column.Board.ProjectID, userID)
	if err != nil {
		return errors.New("access denied: user is not a member of this project")
	}
	if requiredRole == "editor" && (role != "editor" && role != "owner") {
		return errors.New("access denied: editor or owner role required")
	}
	// Для просмотра достаточно быть участником (role может быть viewer/editor/owner)
	return nil
}

// mapTaskToResponse преобразует модель Task в DTO
func mapTaskToResponse(task *models.Task) *dto.TaskResponse {
	resp := &dto.TaskResponse{
		ID:          task.ID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		AssigneeID:  task.AssigneeID,
		Position:    task.Position,
		Labels:      task.Labels,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		CreatorID:   task.CreatorID,
	}
	if task.Assignee != nil {
		resp.Assignee = &dto.UserShort{
			ID:    task.Assignee.ID,
			Name:  task.Assignee.Name,
			Email: task.Assignee.Email,
		}
	}
	if &task.Creator != nil {
		resp.Creator = &dto.UserShort{
			ID:    task.Creator.ID,
			Name:  task.Creator.Name,
			Email: task.Creator.Email,
		}
	}
	return resp
}

func (s *taskService) CreateTask(ctx context.Context, userID string, req *dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	// Валидация запроса
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Проверка доступа: для создания задачи нужно быть редактором проекта
	if err := s.checkProjectAccess(ctx, req.ColumnID, userID, "editor"); err != nil {
		return nil, err
	}

	// Создаём модель задачи
	task := &models.Task{
		ColumnID:    req.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		AssigneeID:  req.AssigneeID,
		CreatorID:   userID,
		Labels:      req.Labels,
		// Позиция будет установлена в последнюю по умолчанию? Лучше определить в репозитории
	}

	// Устанавливаем позицию как максимальную + 1 в данной колонке (или 0 если колонка пуста)
	tasks, err := s.taskRepo.ListByColumn(ctx, req.ColumnID)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		task.Position = 0
	} else {
		maxPos := tasks[len(tasks)-1].Position
		task.Position = maxPos + 1
	}

	// Сохраняем в БД
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	// Инвалидируем кэш списка задач колонки, если он кэшируется
	s.cacheClient.Delete(fmt.Sprintf("column:%s:tasks", req.ColumnID))

	// Для ответа загружаем задачу с присоединёнными пользователями
	task, err = s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}

	return mapTaskToResponse(task), nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string, userID string) (*dto.TaskResponse, error) {
	// Сначала проверяем кэш
	cacheKey := fmt.Sprintf("task:%s", taskID)
	if cached, err := s.cacheClient.Get(cacheKey); err == nil {
		var resp dto.TaskResponse
		if err := json.Unmarshal(cached.Value, &resp); err == nil {
			return &resp, nil
		}
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Проверка доступа: пользователь должен быть участником проекта
	// Пока проверка прав не реализована
	// Проверка доступа: нужно быть редактором проекта
	//if err := s.checkProjectAccess(ctx, task.ColumnID, userID, "editor"); err != nil {
	//	return nil, err
	//}
	resp := mapTaskToResponse(task)

	// Сохраняем в кэш на 5 минут
	data, _ := json.Marshal(resp)
	s.cacheClient.Set(&memcache.Item{Key: cacheKey, Value: data, Expiration: 300})

	return resp, nil
}

func (s *taskService) UpdateTask(ctx context.Context, taskID string, userID string, req *dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	// Валидация (необязательные поля)
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Получаем задачу
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Проверка доступа: нужно быть редактором проекта
	if err := s.checkProjectAccess(ctx, task.ColumnID, userID, "editor"); err != nil {
		return nil, err
	}

	// Обновляем только переданные поля
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.Labels != nil {
		task.Labels = req.Labels
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Инвалидируем кэш задачи и списка колонки
	s.cacheClient.Delete(fmt.Sprintf("task:%s", taskID))
	s.cacheClient.Delete(fmt.Sprintf("column:%s:tasks", task.ColumnID))

	// Получаем обновлённую задачу с прелоадами
	task, err = s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return mapTaskToResponse(task), nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID string, userID string) error {
	// Получаем задачу (нужна для проверки прав)
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	// Проверка доступа: редактор или владелец
	if err := s.checkProjectAccess(ctx, task.ColumnID, userID, "editor"); err != nil {
		return err
	}

	// Удаляем задачу (репозиторий сам обработает каскады, если надо)
	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		return err
	}

	// Инвалидируем кэш
	s.cacheClient.Delete(fmt.Sprintf("task:%s", taskID))
	s.cacheClient.Delete(fmt.Sprintf("column:%s:tasks", task.ColumnID))

	return nil
}

func (s *taskService) MoveTask(ctx context.Context, taskID string, userID string, req *dto.MoveTaskRequest) (*dto.TaskResponse, error) {
	// Валидация
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Получаем задачу (для проверки прав и старой колонки)
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Проверка доступа: нужно быть редактором
	if err := s.checkProjectAccess(ctx, task.ColumnID, userID, "editor"); err != nil {
		return nil, err
	}
	// Также проверяем доступ к новой колонке (она может быть в другом проекте? но обычно в пределах доски)
	if err := s.checkProjectAccess(ctx, req.ColumnID, userID, "editor"); err != nil {
		return nil, err
	}

	// Выполняем перемещение
	if err := s.taskRepo.Move(ctx, taskID, req.ColumnID, req.Position); err != nil {
		return nil, err
	}

	// Инвалидируем кэш старой и новой колонки
	s.cacheClient.Delete(fmt.Sprintf("column:%s:tasks", task.ColumnID))
	s.cacheClient.Delete(fmt.Sprintf("column:%s:tasks", req.ColumnID))
	s.cacheClient.Delete(fmt.Sprintf("task:%s", taskID))

	// Возвращаем обновлённую задачу
	updatedTask, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return mapTaskToResponse(updatedTask), nil
}
