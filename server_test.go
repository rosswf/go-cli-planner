package todo_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rosswf/go-todo-cli"
)

func TestGETTasks(t *testing.T) {
	storage := CreateMockStorage()
	taskList := todo.CreateTaskList(storage)
	server := todo.NewTaskServer(taskList)

	t.Run("test /tasks returns status 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got status %d, want %d", got, want)
		}
	})
}
