// discord @domco00 for help
package db

import (
	"database/sql"
	"devfix/internal/utils"
	"fmt"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {
	appDir := utils.GetAppDir()
	dbPath := filepath.Join(appDir, "devfix.db")

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		action TEXT NOT NULL,
		details TEXT,
		freed_bytes INTEGER DEFAULT 0,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func LogAction(action, details string, freedBytes int64) {
	if DB == nil {
		return
	}
	_, err := DB.Exec("INSERT INTO history (action, details, freed_bytes, timestamp) VALUES (?, ?, ?, ?)",
		action, details, freedBytes, time.Now())
	if err != nil {
		utils.PrintError("Failed to log action to DB: %v", err)
	}
}

func GetHistory() (*sql.Rows, error) {
	if DB == nil {
		return nil, fmt.Errorf("DB not initialized")
	}
	return DB.Query("SELECT action, details, freed_bytes, timestamp FROM history ORDER BY timestamp ASC")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

