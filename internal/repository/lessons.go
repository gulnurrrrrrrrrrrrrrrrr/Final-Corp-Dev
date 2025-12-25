package repository

import (
	"context"

	"quadlingo/internal/models"
)

// Создать новый урок (для manager/admin)
func CreateLesson(lesson *models.Lesson) error {
	query := `
        INSERT INTO lessons (title, description, content, "order", created_by)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at
    `

	err := DB.QueryRow(context.Background(), query, lesson.Title, lesson.Description, lesson.Content, lesson.Order, lesson.CreatedBy).
		Scan(&lesson.ID, &lesson.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Получить все уроки (для всех пользователей — карточки на главной)
func GetAllLessons() ([]models.Lesson, error) {
	query := `
        SELECT id, title, description, content, "order", created_at, created_by
        FROM lessons
        ORDER BY "order" ASC, created_at ASC
    `

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []models.Lesson
	for rows.Next() {
		var l models.Lesson
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.Content, &l.Order, &l.CreatedAt, &l.CreatedBy)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, l)
	}

	return lessons, nil
}

// Получить урок по ID (для изучения)
func GetLessonByID(id int) (models.Lesson, error) {
	query := `
        SELECT id, title, description, content, "order", created_at, created_by
        FROM lessons
        WHERE id = $1
    `

	var lesson models.Lesson
	err := DB.QueryRow(context.Background(), query, id).Scan(
		&lesson.ID, &lesson.Title, &lesson.Description, &lesson.Content, &lesson.Order, &lesson.CreatedAt, &lesson.CreatedBy)
	if err != nil {
		return models.Lesson{}, err
	}

	return lesson, nil
}

// Обновить урок (для manager/admin)
func UpdateLesson(lesson *models.Lesson) error {
	query := `
        UPDATE lessons
        SET title = $1, description = $2, content = $3, "order" = $4
        WHERE id = $5
    `

	_, err := DB.Exec(context.Background(), query, lesson.Title, lesson.Description, lesson.Content, lesson.Order, lesson.ID)
	return err
}

// Удалить урок (для manager/admin)
func DeleteLesson(id int) error {
	query := `DELETE FROM lessons WHERE id = $1`
	_, err := DB.Exec(context.Background(), query, id)
	return err
}
