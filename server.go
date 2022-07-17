package todo

import (
	"encoding/json"
	"log"
	"net/http"
)

type TaskServer struct {
	taskList *TaskList
	http.Handler
}

func NewTaskServer(taskList *TaskList) *TaskServer {
	p := new(TaskServer)
	p.taskList = taskList

	router := http.NewServeMux()
	router.Handle("/tasks", http.HandlerFunc(p.tasksHandler))

	p.Handler = router

	return p
}

func (p *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, _ := p.taskList.GetAll()

	encoder := json.NewEncoder(w)
	err := encoder.Encode(tasks)
	if err != nil {
		log.Printf("Could not encode json %v", err)
	}
}
