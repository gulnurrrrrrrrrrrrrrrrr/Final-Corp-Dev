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
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         Role      `json:"role" db:"role"`
	Points       int       `json:"points" db:"points"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
type Lesson struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" validate:"required"`
	Description string    `json:"description" db:"description"`
	Content     string    `json:"content" db:"content" validate:"required"`
	Order       int       `json:"order" db:"order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
}
type Test struct {
	ID       int    `json:"id" db:"id"`
	LessonID int    `json:"lesson_id" db:"lesson_id"`
	Title    string `json:"title" db:"title" validate:"required"`
}

type Question struct {
	ID            int      `json:"id" db:"id"`
	TestID        int      `json:"test_id" db:"test_id"`
	QuestionText  string   `json:"question_text" db:"question_text" validate:"required"`
	Options       []string `json:"options" db:"options" validate:"min=2"`
	CorrectAnswer int      `json:"correct_answer" db:"correct_answer" validate:"min=0"`
}

type TestSubmission struct {
	TestID  int         `json:"test_id" validate:"required"`
	Answers map[int]int `json:"answers" validate:"required"` // question_id -> selected_option_index
}

type TestResult struct {
	Score        int  `json:"score"`
	Total        int  `json:"total"`
	PointsEarned int  `json:"points_earned"`
	Passed       bool `json:"passed"`
}
