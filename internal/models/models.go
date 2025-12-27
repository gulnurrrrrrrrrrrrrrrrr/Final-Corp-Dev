package models

import "time"

type Role string

const (
	RoleUser    Role = "user"
	RoleManager Role = "manager"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=32"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         Role      `json:"role" db:"role" validate:"oneof=user manager admin"`
	Points       int       `json:"points" db:"points"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Lesson struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description" db:"description"`
	Content     string    `json:"content" db:"content"`
	Order       int       `json:"order" db:"order_num"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
}

type Test struct {
	ID       int    `json:"id" db:"id"`
	LessonID int    `json:"lesson_id" db:"lesson_id"`
	Title    string `json:"title" db:"title"`
}

type Question struct {
	ID            int      `json:"id" db:"id"`
	TestID        int      `json:"test_id" db:"test_id"`
	QuestionText  string   `json:"question_text" db:"question_text"`
	Options       []string `json:"options" db:"options"`
	CorrectAnswer int      `json:"correct_answer" db:"correct_answer"`
}

type TestSubmission struct {
	TestID  int         `json:"test_id"`
	Answers map[int]int `json:"answers"` // question_id -> selected_option_index
}

type TestResult struct {
	Score        int  `json:"score"`
	Total        int  `json:"total"`
	PointsEarned int  `json:"points_earned"`
	Passed       bool `json:"passed"`
}
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Points   int    `json:"points"`
	IsActive bool   `json:"is_active"`
}

type ChangeRoleRequest struct {
	NewRole string `json:"new_role" validate:"required,oneof=user manager"`
}
