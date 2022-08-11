package main

import (
	"log"
	"net/http"

	"github.com/rosswf/go-todo"
	storage "github.com/rosswf/go-todo/storage"
)

func main() {
	storage, _ := storage.CreateSqlite3TaskStorage("tasks.db")
	taskList := todo.CreateTaskList(storage)

	server := todo.NewTaskServer(taskList)

	log.Println("Listening on port 5000...")
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
