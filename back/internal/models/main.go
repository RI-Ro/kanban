package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Doljnost     string    `gorm:"not null;size:255" json:"doljnost"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Name         string    `gorm:"not null;size:100" json:"name"`
	AvatarURL    string    `gorm:"size:500" json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Связи (не загружаются автоматически, если не указать Preload)
	OwnedProjects []Project       `gorm:"foreignKey:OwnerID" json:"-"`
	MemberIn      []ProjectMember `gorm:"foreignKey:UserID" json:"-"`
	TasksAssigned []Task          `gorm:"foreignKey:AssigneeID" json:"-"`
	TasksCreated  []Task          `gorm:"foreignKey:CreatorID" json:"-"`
	Comments      []Comment       `gorm:"foreignKey:AuthorID" json:"-"`
}

// BeforeCreate — хук GORM для генерации UUID, если не задан
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Project struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null;size:200" json:"name"`
	Description *string        `gorm:"type:text" json:"description"`
	OwnerID     string         `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner       User           `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"owner,omitempty"`
	Settings    map[string]any `gorm:"type:jsonb;serializer:json" json:"settings"` // JSONB для гибких настроек
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`

	// Связи
	Boards  []Board         `gorm:"foreignKey:ProjectID" json:"boards,omitempty"`
	Members []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
}

// BeforeCreate для Project
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type ProjectMember struct {
	ProjectID string    `gorm:"type:uuid;primaryKey;index" json:"project_id"`
	UserID    string    `gorm:"type:uuid;primaryKey;index" json:"user_id"`
	Role      string    `gorm:"not null;size:50;default:'viewer'" json:"role"` // owner, editor, viewer
	JoinedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"joined_at"`

	// Связи
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
}

type Board struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID string    `gorm:"type:uuid;not null;index" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	Position  int       `gorm:"not null;default:0" json:"position"` // для сортировки досок в проекте
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Columns []Column `gorm:"foreignKey:BoardID;order:position asc" json:"columns,omitempty"`
}

type Column struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BoardID   string    `gorm:"type:uuid;not null;index" json:"board_id"`
	Board     Board     `gorm:"foreignKey:BoardID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	Position  int       `gorm:"not null;default:0" json:"position"`
	Color     string    `gorm:"size:20" json:"color,omitempty"` // hex-код цвета
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tasks []Task `gorm:"foreignKey:ColumnID;order:position asc" json:"tasks,omitempty"`
}

type Task struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ColumnID    string     `gorm:"type:uuid;not null;index" json:"column_id"`
	Column      Column     `gorm:"foreignKey:ColumnID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Title       string     `gorm:"not null;size:500" json:"title"`
	Description *string    `gorm:"type:text" json:"description"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *string    `gorm:"type:uuid;index" json:"assignee_id,omitempty"`
	Assignee    *User      `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"assignee,omitempty"`
	CreatorID   string     `gorm:"type:uuid;not null;index" json:"creator_id"`
	Creator     User       `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"creator,omitempty"`
	Position    int        `gorm:"not null;default:0" json:"position"`                  // порядок в колонке
	Labels      []string   `gorm:"type:text[];serializer:json" json:"labels,omitempty"` // массив текстовых меток
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Связи
	Comments    []Comment    `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
	Attachments []Attachment `gorm:"foreignKey:TaskID" json:"attachments,omitempty"`
}

type Comment struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    string    `gorm:"type:uuid;not null;index" json:"task_id"`
	Task      Task      `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	AuthorID  string    `gorm:"type:uuid;not null;index" json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"author"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Attachment struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID     string    `gorm:"type:uuid;not null;index" json:"task_id"`
	Task       Task      `gorm:"foreignKey:TaskID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	FileName   string    `gorm:"not null;size:255" json:"file_name"`
	FileURL    string    `gorm:"not null;size:500" json:"file_url"` // ссылка на S3 или локальное хранилище
	UploadedBy string    `gorm:"type:uuid;not null" json:"uploaded_by"`
	Uploader   User      `gorm:"foreignKey:UploadedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"uploader,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ActivityLog struct {
	ID         string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EntityType string         `gorm:"not null;size:50;index" json:"entity_type"` // task, project, board, column
	EntityID   string         `gorm:"not null;index" json:"entity_id"`           // UUID сущности
	Action     string         `gorm:"not null;size:50" json:"action"`            // created, updated, moved, deleted
	UserID     *string        `gorm:"type:uuid;index" json:"user_id,omitempty"`  // кто совершил действие
	User       *User          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"user,omitempty"`
	Data       map[string]any `gorm:"type:jsonb;serializer:json" json:"data"` // дополнительные данные (например, старые/новые значения)
	CreatedAt  time.Time      `json:"created_at"`
}

type Notification struct {
	ID              string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          string    `gorm:"type:uuid;not null;index" json:"user_id"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Type            string    `gorm:"not null;size:50" json:"type"` // comment, assign, mention, etc.
	RelatedEntityID *string   `json:"related_entity_id,omitempty"`  // ID задачи, проекта и т.п.
	IsRead          bool      `gorm:"default:false;index" json:"is_read"`
	CreatedAt       time.Time `json:"created_at"`
	// Можно добавить поле с кратким текстом
	Message string `gorm:"type:text" json:"message,omitempty"`
}

func CreateProjectWithDefaultBoard(db *gorm.DB, ownerID string, projectName string) (*Project, error) {
	project := &Project{
		Name:    projectName,
		OwnerID: ownerID,
		Settings: map[string]any{
			"defaultView": "kanban",
		},
	}

	// Используем транзакцию
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			return err
		}

		// Создаём доску
		board := &Board{
			ProjectID: project.ID,
			Name:      "Основная доска",
			Position:  0,
		}
		if err := tx.Create(board).Error; err != nil {
			return err
		}

		// Создаём стандартные колонки: To Do, In Progress, Done
		columns := []Column{
			{BoardID: board.ID, Name: "Нужно сделать", Position: 0, Color: "#0079bf"},
			{BoardID: board.ID, Name: "В процессе", Position: 1, Color: "#d29034"},
			{BoardID: board.ID, Name: "Готово", Position: 2, Color: "#61bd4f"},
		}
		if err := tx.Create(&columns).Error; err != nil {
			return err
		}

		// Добавляем владельца как участника с ролью owner (можно сделать триггером)
		member := &ProjectMember{
			ProjectID: project.ID,
			UserID:    ownerID,
			Role:      "owner",
		}
		return tx.Create(member).Error
	})

	return project, err
}

func GetBoardWithDetails(db *gorm.DB, boardID string) (*Board, error) {
	var board Board
	err := db.Preload("Columns.Tasks").
		Preload("Columns", func(db *gorm.DB) *gorm.DB {
			return db.Order("position asc")
		}).
		First(&board, "id = ?", boardID).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func MoveTask(db *gorm.DB, taskID string, newColumnID string, newPosition int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var task Task
		if err := tx.First(&task, "id = ?", taskID).Error; err != nil {
			return err
		}

		// Если задача перемещается в другую колонку, нужно пересчитать позиции
		oldColumnID := task.ColumnID

		if oldColumnID != newColumnID {
			// Уменьшаем позиции задач в старой колонке, которые были после этой задачи
			if err := tx.Model(&Task{}).
				Where("column_id = ? AND position > ?", oldColumnID, task.Position).
				Update("position", gorm.Expr("position - 1")).Error; err != nil {
				return err
			}

			// Увеличиваем позиции задач в новой колонке, начиная с newPosition
			if err := tx.Model(&Task{}).
				Where("column_id = ? AND position >= ?", newColumnID, newPosition).
				Update("position", gorm.Expr("position + 1")).Error; err != nil {
				return err
			}

			// Обновляем саму задачу
			task.ColumnID = newColumnID
			task.Position = newPosition
		} else {
			// Перемещение внутри той же колонки
			if newPosition > task.Position {
				// сдвиг вниз: задачи между старым и новым положением сдвигаются вверх
				if err := tx.Model(&Task{}).
					Where("column_id = ? AND position > ? AND position <= ?", oldColumnID, task.Position, newPosition).
					Update("position", gorm.Expr("position - 1")).Error; err != nil {
					return err
				}
			} else if newPosition < task.Position {
				// сдвиг вверх: задачи между новым и старым положением сдвигаются вниз
				if err := tx.Model(&Task{}).
					Where("column_id = ? AND position >= ? AND position < ?", oldColumnID, newPosition, task.Position).
					Update("position", gorm.Expr("position + 1")).Error; err != nil {
					return err
				}
			}
			task.Position = newPosition
		}

		// Сохраняем задачу
		return tx.Save(&task).Error
	})
}
