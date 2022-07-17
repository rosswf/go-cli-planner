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

var dummyData = []todo.Task{
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

func TestGETTasks(t *testing.T) {
	storage := CreateMockStorage(dummyData)
	taskList := todo.CreateTaskList(storage)
	server := todo.NewTaskServer(taskList)

	t.Run("test /tasks returns a list of tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		got := decodeTaskList(t, response.Body)
		want := dummyData

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})

	t.Run("test /tasks/incomplete returns a list of tasks", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks/incomplete", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONContentType(t, response)

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

	t.Run("test /tasks/2 returns the correct task", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks/2", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertJSONContentType(t, response)

		got := decodeTask(t, response.Body)
		want := todo.Task{
			Id:       2,
			Name:     "Task 2",
			Complete: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})

	t.Run("test /tasks/0 return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/tasks/0", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestPOSTTasks(t *testing.T) {
	data := []todo.Task{}

	storage := CreateMockStorage(data)
	taskList := todo.CreateTaskList(storage)
	server := todo.NewTaskServer(taskList)

	t.Run("test POST to /tasks adds a task", func(t *testing.T) {
		jsonData := []byte(`{"Name": "New Task"}`)

		request, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonData))
		request.Header.Set("content-type", "application/json")

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		// Get all
		request, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		got := decodeTaskList(t, response.Body)
		want := []todo.Task{{Id: 1, Name: "New Task", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})

	t.Run("test POST to /tasks/1 marks a task as complete", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/tasks/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		// Get all
		request, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		got := decodeTaskList(t, response.Body)
		want := []todo.Task{{Id: 1, Name: "New Task", Complete: true}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got response %+v, want %+v", got, want)
		}
	})

}

func assertJSONContentType(t testing.TB, response *httptest.ResponseRecorder) {
	t.Helper()

	if response.Result().Header.Get("content-type") != "application/json" {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}

func decodeTaskList(t testing.TB, taskList *bytes.Buffer) []todo.Task {
	t.Helper()
	var got []todo.Task
	err := json.NewDecoder(taskList).Decode(&got)

	if err != nil {
		t.Fatalf("Could not decode json, %v", err)
	}
	return got
}

func decodeTask(t testing.TB, task *bytes.Buffer) todo.Task {
	t.Helper()
	var got todo.Task
	err := json.NewDecoder(task).Decode(&got)

	if err != nil {
		t.Fatalf("Could not decode json, %v", err)
	}
	return got
}
