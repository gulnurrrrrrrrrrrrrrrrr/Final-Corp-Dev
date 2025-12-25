package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quadlingo/internal/middleware"
	"quadlingo/internal/models"
	"quadlingo/internal/services"

	"github.com/gorilla/mux"
)

func CreateLessonHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	if user.Role != models.RoleManager && user.Role != models.RoleAdmin {
		http.Error(w, `{"error": "Only manager or admin can create lessons"}`, http.StatusForbidden)
		return
	}

	var lesson models.Lesson
	if err := json.NewDecoder(r.Body).Decode(&lesson); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	created, err := services.CreateLesson(lesson, user.ID)
	if err != nil {
		http.Error(w, `{"error": "Failed to create lesson"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func GetAllLessonsHandler(w http.ResponseWriter, r *http.Request) {
	lessons, err := services.GetAllLessons()
	if err != nil {
		http.Error(w, `{"error": "Failed to get lessons"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lessons)
}

func GetLessonHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	lesson, err := services.GetLessonByID(id)
	if err != nil {
		http.Error(w, `{"error": "Lesson not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lesson)
}
