package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

var (
	db *sqlx.DB
)

func count(db *sqlx.DB) (int, error) {
	var count int
	return count, db.Get(&count, "SELECT count(id) FROM scheduler")
}

func openDB() (*sqlx.DB, error) {
	dbFile := os.Getenv("DATABASE_PATH")
	if dbFile == "" {
		dbFile = "../scheduler.db"
	}
	db, err := sqlx.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func dbCheck() (*sqlx.DB, error) {
	mu.Lock()
	defer mu.Unlock()
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}
	dbFile := os.Getenv("DATABASE_PATH")
	appDir := filepath.Dir(appPath)
	dbFile = filepath.Join(appDir, dbFile)
	_, err = os.Stat(dbFile)
	if err != nil {
		log.Fatal()
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXIST scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Date CHAR(8) NOT NULL, 
		Title TEXT NOT NULL,
		Сomment TEXT NOT NULL, 
		Repeat VARCHAR(128) NOT NULL)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE INDEX task_date ON scheduler(date)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}
