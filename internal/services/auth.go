package services

import (
	"context"
	"fmt"
	"log"

	"quadlingo/internal/config"
	"quadlingo/internal/models"
	"quadlingo/internal/repository"
	"quadlingo/internal/utils"
)

func Register(req models.RegisterRequest, cfg *config.Config) (*models.AuthResponse, error) {
	// Проверяем, есть ли уже пользователи
	var count int
	err := repository.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return nil, err
	}

	role := models.RoleUser
	if count == 0 {
		role = models.RoleAdmin
		log.Println("First user registered — granted ADMIN role!")
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         role,
		Points:       0,
	}

	query := `INSERT INTO users (username, email, password_hash, role, points) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err = repository.DB.QueryRow(context.Background(), query, user.Username, user.Email, user.PasswordHash, user.Role, user.Points).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role), cfg)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func Login(req models.LoginRequest, cfg *config.Config) (*models.AuthResponse, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, role, points, created_at, updated_at 
              FROM users WHERE username = $1 OR email = $1`

	err := repository.DB.QueryRow(context.Background(), query, req.UsernameOrEmail).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Points, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err // пользователь не найден или ошибка
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateJWT(user.ID, string(user.Role), cfg)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}
