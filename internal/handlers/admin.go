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

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetCurrentUser(r)
	if currentUser.Role != models.RoleAdmin {
		http.Error(w, `{"error": "Только админ может видеть пользователей"}`, http.StatusForbidden)
		return
	}

	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, `{"error": "Ошибка получения пользователей"}`, http.StatusInternalServerError)
		return
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = models.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
			Points:   u.Points,
			IsActive: true,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ChangeUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := middleware.GetCurrentUser(r)
	if currentUser.Role != models.RoleAdmin {
		http.Error(w, `{"error": "Только админ может менять роли"}`, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"])

	var req models.ChangeRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Неверный запрос"}`, http.StatusBadRequest)
		return
	}

	if req.NewRole != "user" && req.NewRole != "manager" {
		http.Error(w, `{"error": "Роль может быть только user или manager"}`, http.StatusBadRequest)
		return
	}

	err := repository.UpdateUserRole(userID, req.NewRole)
	if err != nil {
		http.Error(w, `{"error": "Ошибка изменения роли"}`, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "Роль успешно изменена"}`))
}
