package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Tasks struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var mu sync.Mutex

func TaskHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	defer mu.Unlock()
	var post Tasks
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAlreadyReported)
		return 0, nil
	}
	if r.Method == http.MethodPost {
		var task Tasks
		res, err := db.Exec("INSERT INTO clients (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
			sql.Named("date", task.Date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat))
		if err != nil {
			return 0, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		return int(id), nil
	}
	if err = json.Unmarshal(buf.Bytes(), &post); err != nil {
		log.Fatal("Ошибка десериализации JSON:", err)
		return 0, nil
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return 0, nil
}
