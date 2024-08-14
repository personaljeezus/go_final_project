package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTask(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}
func DbCheck() (*sqlx.DB, error) {
	godotenv.Load("ENV_PATH")
	appPath := os.Getenv("DATABASE_PATH")
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	if appPath == "" {
		dbFile = "./scheduler.db"
	}
	_, err := os.Stat(dbFile)
	var install bool
	if os.IsNotExist(err) {
		install = true
	}
	db, err := sqlx.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal("Failed to open db")
	}
	if install {
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (id INTEGER PRIMARY KEY AUTOINCREMENT, date CHAR(8) NOT NULL, title TEXT NOT NULL, comment TEXT NOT NULL, repeat VARCHAR(128) NOT NULL)`)
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(`CREATE INDEX IF NOT EXISTS task_date ON scheduler(date)`)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
