package db

import (
	"database/sql"
	"fmt"
)

// InitSchema creates all tables if they don't exist.
func InitSchema(db *sql.DB) error {
	schema := `
	-- Cached from Nutrition Reference.md
	CREATE TABLE IF NOT EXISTS foods (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		category TEXT,
		calories_per_serving REAL,
		protein_per_serving REAL,
		fat_per_serving REAL,
		carbs_per_serving REAL,
		serving_description TEXT,
		is_favorite BOOLEAN DEFAULT FALSE,
		source_line TEXT
	);

	-- Custom foods added via app (no macros)
	CREATE TABLE IF NOT EXISTS custom_foods (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Daily log entries
	CREATE TABLE IF NOT EXISTS daily_logs (
		id INTEGER PRIMARY KEY,
		date DATE NOT NULL UNIQUE,
		day_type TEXT CHECK(day_type IN ('rest', 'workout')) DEFAULT 'rest',
		weight_kg REAL,
		notes TEXT,
		synced_to_obsidian BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Individual meal entries
	CREATE TABLE IF NOT EXISTS meal_entries (
		id INTEGER PRIMARY KEY,
		daily_log_id INTEGER REFERENCES daily_logs(id),
		meal_type TEXT CHECK(meal_type IN ('breakfast', 'lunch', 'dinner', 'snacks', 'shake')),
		food_id INTEGER REFERENCES foods(id),
		custom_food_id INTEGER REFERENCES custom_foods(id),
		quantity REAL DEFAULT 1,
		saved_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Offline sync queue (future use)
	CREATE TABLE IF NOT EXISTS sync_queue (
		id INTEGER PRIMARY KEY,
		action TEXT,
		payload TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		synced_at DATETIME
	);
	`

	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	return nil
}
