package planner_test

import (
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	planner "github.com/rosswf/go-cli-planner"
	storage "github.com/rosswf/go-cli-planner/storage"
)

type MockTaskStorage struct {
	taskList []planner.Task
}

func (m *MockTaskStorage) Add(task *planner.Task) error {
	id := len(m.taskList) + 1
	task.Id = planner.TaskId(id)
	m.taskList = append(m.taskList, *task)
	return nil
}

func (m *MockTaskStorage) GetAll() ([]planner.Task, error) {
	return m.taskList, nil
}

func (m *MockTaskStorage) ToggleStatus(id planner.TaskId) error {
	task, _ := m.GetTask(id)
	if task.Complete {
		task.Complete = false
	} else {
		task.Complete = true
	}
	return nil
}

func (m *MockTaskStorage) GetTask(id planner.TaskId) (*planner.Task, error) {
	id-- // slice is 0 indexed
	return &m.taskList[id], nil
}

func (m *MockTaskStorage) GetOutstanding() ([]planner.Task, error) {
	outstanding := make([]planner.Task, 0)

	for _, task := range m.taskList {
		if !task.Complete {
			outstanding = append(outstanding, task)
		}
	}
	return outstanding, nil
}

func (m *MockTaskStorage) Delete(id planner.TaskId) error {
	id-- // slice is 0 indexed
	m.taskList = append(m.taskList[:id], m.taskList[id+1:]...)
	return nil
}

func CreateMockStorage() *MockTaskStorage {
	return &MockTaskStorage{}
}

func TestTasks(t *testing.T) {
	storage := CreateMockStorage()

	taskList := planner.CreateTaskList(storage)

	t.Run("A task is added to the task list", func(t *testing.T) {
		err := taskList.Add("Task 1")
		AssertNoError(t, err)

		got, err := taskList.GetAll()
		AssertNoError(t, err)

		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %#v, want %#v", got, want)
		}
	})

	t.Run("A task is marked as completed", func(t *testing.T) {
		tasks, err := taskList.GetAll()
		AssertNoError(t, err)

		err = taskList.ToggleStatus(&tasks[0])
		AssertNoError(t, err)

		got, err := taskList.GetAll()
		AssertNoError(t, err)

		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: true}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("A task is marked as incomplete", func(t *testing.T) {
		tasks, err := taskList.GetAll()

		err = taskList.ToggleStatus(&tasks[0])
		AssertNoError(t, err)

		got, err := taskList.GetAll()

		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("Test only show incomplete", func(t *testing.T) {
		taskList.Add("Task 2")
		taskList.Add("Task 3")
		taskList.Add("Task 4")

		tasks, err := taskList.GetAll()

		err = taskList.ToggleStatus(&tasks[0])
		err = taskList.ToggleStatus(&tasks[2])

		got, err := taskList.GetOutstanding()
		AssertNoError(t, err)

		want := []planner.Task{
			{Id: 2, Name: "Task 2", Complete: false},
			{Id: 4, Name: "Task 4", Complete: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("Delete a task", func(t *testing.T) {
		tasks, err := taskList.GetAll()

		err = taskList.Delete(&tasks[1])
		AssertNoError(t, err)

		got, err := taskList.GetAll()
		AssertNoError(t, err)

		want := []planner.Task{
			{Id: 1, Name: "Task 1", Complete: true},
			{Id: 3, Name: "Task 3", Complete: true},
			{Id: 4, Name: "Task 4", Complete: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

func TestSqlite3TaskStorage(t *testing.T) {
	storage, err := storage.CreateSqlite3TaskStorage(":memory:")
	AssertNoError(t, err)

	t.Run("Add a new task to sqlite storage", func(t *testing.T) {
		task := planner.Task{Name: "Task 1", Complete: false}

		err := storage.Add(&task)
		AssertNoError(t, err)

		got, err := storage.GetAll()
		AssertNoError(t, err)

		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("Get a single task by ID", func(t *testing.T) {
		task := planner.Task{Name: "Task 2", Complete: true}
		storage.Add(&task)

		got, err := storage.GetTask(2)
		AssertNoError(t, err)

		want := &planner.Task{Id: 2, Name: "Task 2", Complete: true}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("Toggle task status", func(t *testing.T) {
		err := storage.ToggleStatus(1)
		AssertNoError(t, err)

		got, _ := storage.GetTask(1)

		want := &planner.Task{Id: 1, Name: "Task 1", Complete: true}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
