package repository

import (
	"context"

	"quadlingo/internal/models"
)

// Получить всех пользователей (для админа)
func GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, role, points, is_active FROM users ORDER BY id`
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var roleStr string
		var isActive bool
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &roleStr, &u.Points, &isActive)
		if err != nil {
			return nil, err
		}
		u.Role = models.Role(roleStr)
		u.IsActive = isActive
		users = append(users, u)
	}
	return users, nil
}

// Изменить роль пользователя (только админ)
func UpdateUserRole(userID int, newRole string) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	_, err := DB.Exec(context.Background(), query, newRole, userID)
	return err
}
