package todo

import "net/http"

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
	w.WriteHeader(http.StatusOK)
}
