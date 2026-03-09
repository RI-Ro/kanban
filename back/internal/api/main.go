package main

import (
	"log"
	"net/http"
	"rri/task-back/internal/api/handlers"

	//"rri/task-back/internal/api/middleware"
	"rri/task-back/internal/cache"
	"rri/task-back/internal/repository"
	"rri/task-back/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Mainapi() {
	// Подключение к БД
	dsn := "host=localhost user=postgres password=secret dbname=yougile port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
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

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", r)
}
