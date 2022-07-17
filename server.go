package todo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TaskServer struct {
	taskList *TaskList
	http.Handler
}

func NewTaskServer(taskList *TaskList) *TaskServer {
	p := new(TaskServer)
	p.taskList = taskList

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", p.tasksHandler)
		r.Get("/incomplete", p.incompleteHandler)
		r.Get("/{taskID:[0-9]+}", p.taskHandler)
	})

	p.Handler = r
	return p
}

func (p *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskList.GetAll()
	if err != nil {
		log.Printf("Could not get tasks %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasks)
}

func (p *TaskServer) incompleteHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskList.GetOutstanding()
	if err != nil {
		log.Printf("Could not get tasks %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, tasks)
}

func (p *TaskServer) taskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid taskID given %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	taskId := TaskId(id)

	task, err := p.taskList.GetOne(taskId)
	if err != nil {
		log.Printf("Could not get task with id %d, %v", taskId, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		log.Printf("Could not encoide json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func writeJSON(w http.ResponseWriter, tasks []Task) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(tasks)

	if err != nil {
		log.Printf("Could not encode json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
