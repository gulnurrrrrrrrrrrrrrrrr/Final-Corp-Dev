package services

import (
	"fmt"
	"quadlingo/internal/models"
	"quadlingo/internal/repository"
)

func CreateLesson(lesson models.Lesson, userID int) (models.Lesson, error) {
	lesson.CreatedBy = userID
	err := repository.CreateLesson(&lesson)
	return lesson, err
}

func GetAllLessons() ([]models.Lesson, error) {
	return repository.GetCachedLessons()
}

func GetLessonByID(id int) (models.Lesson, error) {
	return repository.GetLessonByID(id)
}

func UpdateLesson(lesson models.Lesson, userID int) error {
	// Проверяем, что пользователь — создатель урока или admin
	existing, err := repository.GetLessonByID(lesson.ID)
	if err != nil {
		return err
	}
	if existing.CreatedBy != userID && userID != 1 { // 1 — пример admin, лучше проверять роль
		return fmt.Errorf("permission denied")
	}
	return repository.UpdateLesson(&lesson)
}

func DeleteLesson(id int, userID int) error {
	lesson, err := repository.GetLessonByID(id)
	if err != nil {
		return err
	}
	if lesson.CreatedBy != userID {
		return fmt.Errorf("permission denied")
	}
	return repository.DeleteLesson(id)
}
