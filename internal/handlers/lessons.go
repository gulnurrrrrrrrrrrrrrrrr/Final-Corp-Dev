package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quadlingo/internal/middleware"
	"quadlingo/internal/models"
	"quadlingo/internal/repository"

	"github.com/gorilla/mux"
)

func GetAllLessonsHandler(w http.ResponseWriter, r *http.Request) {
	lessons, err := repository.GetAllLessons()
	if err != nil {
		http.Error(w, "Ошибка загрузки уроков", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons) // Возвращаем []models.Lesson
}

func GetLessonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID урока", http.StatusBadRequest)
		return
	}

	lesson, err := repository.GetLessonByID(id)
	if err != nil {
		http.Error(w, "Урок не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}

func CreateLessonHandler(w http.ResponseWriter, r *http.Request) {
	var lesson models.Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Получаем ID менеджера из токена
	currentUser := middleware.GetCurrentUser(r)
	lesson.CreatedBy = currentUser.ID

	if err := repository.CreateLesson(&lesson); err != nil {
		http.Error(w, "Ошибка создания урока", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lesson)
}
