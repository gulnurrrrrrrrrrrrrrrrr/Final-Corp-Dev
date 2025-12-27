package main

import (
	"log"
	"net/http"
	"os"

	"quadlingo/internal/config"
	"quadlingo/internal/handlers"
	"quadlingo/internal/middleware"
	"quadlingo/internal/repository"
	"quadlingo/internal/utils"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/unrolled/secure"
	"go.uber.org/zap"
)

func main() {
	// Загружаем конфиг из .env
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Инициализируем JWT секрет
	utils.InitJWT(cfg.JWTSecret)

	// Подключение к PostgreSQL
	if err := repository.InitDB(cfg); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer repository.CloseDB()

	// Миграции таблиц
	if err := repository.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Подключение к Redis
	if err := repository.InitRedis(cfg); err != nil {
		log.Fatal("Cannot connect to Redis:", err)
	}

	// Инициализация zap логгера
	var zapLogger *zap.Logger
	if os.Getenv("ENV") == "production" {
		zapLogger, _ = zap.NewProduction()
	} else {
		zapLogger, _ = zap.NewDevelopment()
	}
	defer zapLogger.Sync()

	// Передаём логгер в middleware
	middleware.InitLogger(zapLogger)
	// Роутер
	r := mux.NewRouter()

	// Prometheus метрики
	r.Handle("/metrics", promhttp.Handler())

	// Security Headers
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		STSSeconds:            31536000,
		STSIncludeSubdomains:  true,
	})

	// Глобальные middleware
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.LoggingMiddleware)
	r.Use(secureMiddleware.Handler)

	// Главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1 style='text-align:center;margin-top:100px;font-family:system-ui'>🚀 QuadLingo API жұмыс істеп тұр! 🇰🇿</h1><p style='text-align:center'>Фронтенд: <a href='/static/index.html'>/static/index.html</a> | Метрики: <a href='/metrics'>/metrics</a></p>"))
	}).Methods("GET")

	// Публичные маршруты
	r.HandleFunc("/register", handlers.RegisterHandler(cfg)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(cfg)).Methods("POST")

	// Публичные уроки (для всех)
	r.HandleFunc("/lessons", handlers.GetAllLessonsHandler).Methods("GET")
	r.HandleFunc("/lessons/{id}", handlers.GetLessonHandler).Methods("GET")

	// Защищённые маршруты — требуют авторизации
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Профиль пользователя
	protected.HandleFunc("/profile", handlers.ProfileHandler).Methods("GET")

	// === МАРШРУТЫ ТОЛЬКО ДЛЯ МЕНЕДЖЕРА ===
	managerRouter := protected.PathPrefix("").Subrouter() // пустой префикс — действует на все маршруты ниже
	managerRouter.Use(middleware.RequireRole("manager"))

	// Создание урока — только менеджер
	managerRouter.HandleFunc("/lessons", handlers.CreateLessonHandler).Methods("POST")

	// Создание теста — только менеджер (теперь с middleware проверки роли)
	managerRouter.HandleFunc("/tests", handlers.CreateTestHandler).Methods("POST")

	// Если в будущем добавишь ещё менеджерские endpoints — добавляй их сюда

	// === АДМИНСКИЕ МАРШРУТЫ ===
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)
	adminRouter.Use(middleware.RequireRole("admin"))

	adminRouter.HandleFunc("/users", handlers.GetAllUsersHandler).Methods("GET")
	adminRouter.HandleFunc("/users/{id}/role", handlers.ChangeUserRoleHandler).Methods("PATCH")

	// Статический фронтенд
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Запуск сервера
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 QuadLingo сервер запущен на http://localhost:%s", port)
	log.Printf("   Фронтенд: http://localhost:%s/static/index.html", port)
	log.Printf("   Метрики:  http://localhost:%s/metrics", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
