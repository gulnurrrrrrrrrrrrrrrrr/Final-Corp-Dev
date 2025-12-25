package handlers

import (
	"encoding/json"
	"net/http"

	"quadlingo/internal/middleware"
	"quadlingo/internal/models"
	"quadlingo/internal/services"
)

func SubmitTestHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)

	var submission models.TestSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, `{"error": "Invalid request"}`, http.StatusBadRequest)
		return
	}

	result, err := services.SubmitTest(submission, user.ID)
	if err != nil {
		http.Error(w, `{"error": "Failed to submit test"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
