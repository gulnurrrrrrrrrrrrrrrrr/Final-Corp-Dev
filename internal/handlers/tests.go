package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"quadlingo/internal/middleware"
	"quadlingo/internal/repository"
)

func CreateTestHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetCurrentUser(r)
	if currentUser.Role != "manager" {
		http.Error(w, `{"error": "Только менеджер может создавать тесты"}`, http.StatusForbidden)
		return
	}

	var req struct {
		LessonID  int `json:"lesson_id"`
		Questions []struct {
			QuestionText  string   `json:"question_text"`
			Options       []string `json:"options"`
			CorrectAnswer int      `json:"correct_answer"`
		} `json:"questions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Ошибка парсинга теста: %v", err)
		http.Error(w, `{"error": "Неверный формат данных"}`, http.StatusBadRequest)
		return
	}

	if len(req.Questions) == 0 {
		http.Error(w, `{"error": "Тест должен содержать хотя бы один вопрос"}`, http.StatusBadRequest)
		return
	}

	err := repository.CreateTestForLesson(req.LessonID, req.Questions)
	if err != nil {
		log.Printf("Ошибка создания теста: %v", err)
		http.Error(w, `{"error": "Ошибка создания теста"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("Менеджер %d создал тест для урока %d (%d вопросов)", currentUser.ID, req.LessonID, len(req.Questions))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "Тест успешно создан!",
		"lesson_id":       req.LessonID,
		"questions_count": len(req.Questions),
	})
}
