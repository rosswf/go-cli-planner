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
		r.Use(setJsonContentType)
		r.Get("/", p.tasksHandler)
		r.Post("/", p.newTaskHandler)
		r.Get("/incomplete", p.incompleteHandler)
		r.Get("/{taskID:^[1-9][0-9]*}", p.taskHandler)
		r.Post("/{taskID:^[1-9][0-9]*}", p.taskStatusToggleHandler)
		r.Delete("/{taskID:^[1-9][0-9]*}", p.taskDeleteHandler)
	})

	p.Handler = r
	return p
}

func setJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (p *TaskServer) tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskList.GetAll()
	if err != nil {
		log.Printf("Could not get tasks %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeTasksJSON(w, tasks)
}

func (p *TaskServer) incompleteHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskList.GetOutstanding()
	if err != nil {
		log.Printf("Could not get tasks %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeTasksJSON(w, tasks)
}

func (p *TaskServer) taskHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid taskID given %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskId := TaskId(id)

	task, err := p.taskList.GetOne(taskId)
	if err != nil {
		log.Printf("Could not get task with id %d, %v", taskId, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(task)

	if err != nil {
		log.Printf("Could not encode json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *TaskServer) newTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		log.Printf("Could not decode json %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)

	err = p.taskList.Add(task.Name)
	if err != nil {
		log.Printf("Could not add task %v, %v", task, err)
	}
}

func (p *TaskServer) taskStatusToggleHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid taskID given %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskId := TaskId(id)

	task, err := p.taskList.GetOne(taskId)
	if err != nil {
		log.Printf("Could not get task with id %d, %v", taskId, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = p.taskList.ToggleStatus(&task)
	if err != nil {
		log.Printf("Could not toggle status of task %v, %v", task, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *TaskServer) taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "taskID")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Invalid taskID given %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskId := TaskId(id)

	task, err := p.taskList.GetOne(taskId)
	if err != nil {
		log.Printf("Could not get task with id %d, %v", taskId, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = p.taskList.Delete(&task)
	if err != nil {
		log.Printf("Could not delete task %v, %v", task, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func writeTasksJSON(w http.ResponseWriter, tasks []Task) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(tasks)

	if err != nil {
		log.Printf("Could not encode json %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
