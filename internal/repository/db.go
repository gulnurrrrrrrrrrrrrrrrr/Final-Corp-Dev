package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"quadlingo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) error {
	var err error
	for i := 0; i < 10; i++ {
		DB, err = pgxpool.New(context.Background(), cfg.DSN())
		if err == nil {
			if err = DB.Ping(context.Background()); err == nil {
				log.Println("Successfully connected to PostgreSQL!")
				return nil
			}
		}
		log.Printf("Failed to connect to DB (attempt %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return fmt.Errorf("failed to connect to database after 10 attempts: %w", err)
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
func Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            role VARCHAR(20) NOT NULL DEFAULT 'user',
            points INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT true,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );`,

		`CREATE TABLE IF NOT EXISTS lessons (
            id SERIAL PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            content TEXT NOT NULL,
            "order" INTEGER DEFAULT 0,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            created_by INTEGER REFERENCES users(id) ON DELETE SET NULL
        );`,

		`CREATE TABLE IF NOT EXISTS tests (
            id SERIAL PRIMARY KEY,
            lesson_id INTEGER REFERENCES lessons(id) ON DELETE CASCADE,
            title VARCHAR(255) NOT NULL
        );`,

		`CREATE TABLE IF NOT EXISTS questions (
            id SERIAL PRIMARY KEY,
            test_id INTEGER REFERENCES tests(id) ON DELETE CASCADE,
            question_text TEXT NOT NULL,
            options JSONB NOT NULL,
            correct_answer INTEGER NOT NULL
        );`,

		`CREATE TABLE IF NOT EXISTS user_progress (
            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
            lesson_id INTEGER REFERENCES lessons(id) ON DELETE CASCADE,
            completed BOOLEAN DEFAULT FALSE,
            test_score INTEGER,
            completed_at TIMESTAMP WITH TIME ZONE,
            PRIMARY KEY (user_id, lesson_id)
        );`,
	}

	ctx := context.Background()
	for _, q := range queries {
		_, err := DB.Exec(ctx, q)
		if err != nil {
			log.Printf("Error executing query: %v\nQuery: %s", err, q)
			return err
		}
	}

	log.Println("All tables created successfully (migrations applied)!")
	return nil

}
