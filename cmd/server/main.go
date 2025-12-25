package main

import (
	"log"
	"net/http"

	"quadlingo/internal/config"
	"quadlingo/internal/handlers"
	"quadlingo/internal/middleware"
	"quadlingo/internal/repository"
	"quadlingo/internal/utils"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Инициализируем секрет для JWT
	utils.InitJWT(cfg.JWTSecret)

	if err := repository.InitDB(cfg); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer repository.CloseDB()

	if err := repository.Migrate(); err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := mux.NewRouter()

	// Главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("QuadLingo API is running! 🇰🇿"))
	}).Methods("GET")

	r.HandleFunc("/register", handlers.RegisterHandler(cfg)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(cfg)).Methods("POST")

	r.HandleFunc("/lessons", handlers.GetAllLessonsHandler).Methods("GET")

	// Защищённые роуты
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Только manager и admin могут создавать уроки
	protected.HandleFunc("/lessons", handlers.CreateLessonHandler).Methods("POST")

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
