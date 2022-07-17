package todo_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rosswf/go-todo-cli"
)

func TestGETTasks(t *testing.T) {
	data := []todo.Task{
		{
			Id:       1,
			Name:     "Task 1",
			Complete: false,
		},
		{
			Id:       2,
			Name:     "Task 2",
			Complete: true,
		},
	}

	storage := CreateMockStorage(data)
	taskList := todo.CreateTaskList(storage)
	server := todo.NewTaskServer(taskList)

	t.Run("test /tasks returns a list of tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		got := decodeTaskList(t, response.Body)
		want := data

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})

	t.Run("test /tasks/incomplete returns a list of tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks/incomplete", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		got := decodeTaskList(t, response.Body)
		want := []todo.Task{
			{
				Id:       1,
				Name:     "Task 1",
				Complete: false,
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}

func decodeTaskList(t testing.TB, taskList *bytes.Buffer) []todo.Task {
	var got []todo.Task
	err := json.NewDecoder(taskList).Decode(&got)

	if err != nil {
		t.Fatalf("Could not decode json, %v", err)
	}
	return got
}
