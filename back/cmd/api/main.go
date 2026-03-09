package main

/*author ri-rogov*/

import (
	"fmt"
	"log"
	"net/http"
	"rri/task-back/internal/api/handlers"
	"rri/task-back/internal/cache"
	"rri/task-back/internal/config"
	"rri/task-back/internal/models"
	"rri/task-back/internal/repository"
	"rri/task-back/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	conf := config.CreateConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=secret dbname=%s port=%s sslmode=disable TimeZone=UTC", conf.DB_IPADDRESS, conf.DB_USER, conf.DB_NAME, conf.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if conf.Development {
		// Автоматическое создание/обновление схемы (только для разработки)
		if err := db.AutoMigrate(
			&models.User{},
			&models.Project{},
			&models.ProjectMember{},
			&models.Board{},
			&models.Column{},
			&models.Task{},
			&models.Comment{},
			&models.Attachment{},
			&models.ActivityLog{},
			&models.Notification{},
		); err != nil {
			log.Fatal("failed to migrate:", err)
		}
	}
	// Подключение к Memcached
	mc := cache.New([]string{"localhost:11211"})

	// Инициализация репозиториев
	taskRepo := repository.NewTaskRepository(db)
	columnRepo := repository.NewColumnRepository(db) // нужно реализовать аналогично
	memberRepo := repository.NewProjectMemberRepository(db)

	// Инициализация сервисов
	taskService := service.NewTaskService(taskRepo, columnRepo, memberRepo, mc)

	// Инициализация хендлеров
	taskHandler := handlers.NewTaskHandler(taskService)

	// Маршрутизатор chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(middleware.CORS) // настройте под свой фронтенд

	// Группа API с аутентификацией
	r.Group(func(r chi.Router) {
		//	r.Use(middleware.AuthMiddleware) // ваш middleware для проверки JWT

		r.Route("/api/tasks", func(r chi.Router) {
			r.Post("/", taskHandler.CreateTask)
			r.Get("/{id}", taskHandler.GetTask)
			r.Put("/{id}", taskHandler.UpdateTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
			r.Post("/{id}/move", taskHandler.MoveTask)
		})

		// другие маршруты для проектов, колонок и т.д.
	})

	server := fmt.Sprintf("%s:%s", conf.SERVER_IPADDRESS, conf.SERVER_PORT)

	log.Printf("Server started on %s", server)
	http.ListenAndServe(server, r)
}

/*
	user := models.User{
		Email:        "ri-rogov",
		Doljnost:     "doljnost",
		PasswordHash: "wwsdvlkndflvkmn lkn fvlwkenokdn vlkmadfl;vbm dlkfm",
		Name:         "name",
		AvatarURL:    "http://lskjvpokdnfbpokne",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	// Create a single record
	ctx := context.Background()
	err = gorm.G[models.User](db).Create(ctx, &user) // pass pointer of data to Create


*/
