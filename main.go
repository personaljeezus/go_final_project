package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if nowParam == "" || dateParam == "" || repeatParam == "" {
		http.Error(w, "Поля параметров пусты", http.StatusBadRequest)
		return
	}

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Неверный параметр", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("20060102", dateParam)
	if err != nil {
		http.Error(w, "Неверный параметр", http.StatusBadRequest)
		return
	}

	nextDate, err := nextWeekday(now, date, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(nextDate.Format("20060102")))

}

func THandler(w http.ResponseWriter, r *http.Request) {
	dateParam := r.FormValue("date")
	titleParam := r.FormValue("title")
	repeatParam := r.FormValue("repeat")

	if titleParam == " " {
		http.Error(w, "Не указан заголовок задачи", http.StatusRequestURITooLong)
		return
	}
	if dateParam == " " {
		date, _ := time.Parse("20060102", dateParam)
		http.Error(w, "Неверная дата задачи", http.StatusTeapot)
		nextWeekday(time.Now(), date, repeatParam)
		return
	}
}
func main() {
	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "7540"
		fmt.Printf("Сервер слушает порт %s\n", defaultPort)
	}
	defer db.Close()
	dir := http.Dir("./web/")
	webFile := http.FileServer(dir)
	mux := http.NewServeMux()
	mux.Handle("/", webFile)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/task", THandler)
	err := http.ListenAndServe(":"+defaultPort, mux)
	if err != nil {
		panic(err)
	}
}
