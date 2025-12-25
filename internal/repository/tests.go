package repository

import (
	"context"
	"encoding/json"

	"quadlingo/internal/models"
)

func CreateTest(test *models.Test) error {
	query := `INSERT INTO tests (lesson_id, title) VALUES ($1, $2) RETURNING id`
	return DB.QueryRow(context.Background(), query, test.LessonID, test.Title).Scan(&test.ID)
}

func GetTestByLessonID(lessonID int) (models.Test, []models.Question, error) {
	var test models.Test
	query := `SELECT id, lesson_id, title FROM tests WHERE lesson_id = $1`
	err := DB.QueryRow(context.Background(), query, lessonID).Scan(&test.ID, &test.LessonID, &test.Title)
	if err != nil {
		return test, nil, err
	}

	questionsQuery := `SELECT id, test_id, question_text, options, correct_answer FROM questions WHERE test_id = $1 ORDER BY id`
	rows, err := DB.Query(context.Background(), questionsQuery, test.ID)
	if err != nil {
		return test, nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		var options []byte // JSONB
		err := rows.Scan(&q.ID, &q.TestID, &q.QuestionText, &options, &q.CorrectAnswer)
		if err != nil {
			return test, nil, err
		}
		if err := json.Unmarshal(options, &q.Options); err != nil {
			return test, nil, err
		}
		questions = append(questions, q)
	}

	return test, questions, nil
}

func CreateQuestion(q *models.Question) error {
	optionsJSON, _ := json.Marshal(q.Options)
	query := `INSERT INTO questions (test_id, question_text, options, correct_answer) VALUES ($1, $2, $3, $4) RETURNING id`
	return DB.QueryRow(context.Background(), query, q.TestID, q.QuestionText, optionsJSON, q.CorrectAnswer).Scan(&q.ID)
}

func SaveUserProgress(userID, lessonID int, score *int) error {
	query := `
        INSERT INTO user_progress (user_id, lesson_id, completed, test_score)
        VALUES ($1, $2, true, $3)
        ON CONFLICT (user_id, lesson_id) DO UPDATE SET completed = true, test_score = $3, completed_at = CURRENT_TIMESTAMP
    `
	_, err := DB.Exec(context.Background(), query, userID, lessonID, score)
	if err != nil {
		return err
	}

	// Начисляем очки (10 за каждый правильный ответ)
	if score != nil {
		points := *score * 10
		_, err = DB.Exec(context.Background(), "UPDATE users SET points = points + $1 WHERE id = $2", points, userID)
	}
	return err
}
