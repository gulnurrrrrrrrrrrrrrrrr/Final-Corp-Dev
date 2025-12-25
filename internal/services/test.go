package services

import (
	"quadlingo/internal/models"
	"quadlingo/internal/repository"
)

func SubmitTest(submission models.TestSubmission, userID int) (models.TestResult, error) {
	_, questions, err := repository.GetTestByLessonID(submission.TestID)
	if err != nil {
		return models.TestResult{}, err
	}

	score := 0
	for _, q := range questions {
		if selected, ok := submission.Answers[q.ID]; ok && selected == q.CorrectAnswer {
			score++
		}
	}

	result := models.TestResult{
		Score:        score,
		Total:        len(questions),
		PointsEarned: score * 10,
		Passed:       score*10 >= len(questions)*7, // 70% = 7 из 10
	}

	// Сохраняем прогресс и начисляем очки
	err = repository.SaveUserProgress(userID, submission.TestID, &score) // lesson_id = test_id для простоты
	return result, err
}
