package handlers

import (
	"encoding/json"
	"net/http"

	"quadlingo/internal/middleware"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	response := map[string]interface{}{
		"message": "Сәлеметсіз бе! Добро пожаловать в профиль QuadLingo! 🇰🇿",
		"user": map[string]interface{}{
			"id":     user.ID,
			"role":   user.Role,
			"points": user.Points,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
