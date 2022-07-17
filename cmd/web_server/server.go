package main

import (
	"log"
	"net/http"

	"github.com/rosswf/go-todo-cli"
	storage "github.com/rosswf/go-todo-cli/storage"
)

func main() {
	storage, _ := storage.CreateSqlite3TaskStorage("tasks.db")
	taskList := todo.CreateTaskList(storage)

	server := todo.NewTaskServer(taskList)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
