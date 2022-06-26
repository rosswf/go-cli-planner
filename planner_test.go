package planner_test

import (
	"reflect"
	"testing"

	planner "github.com/rosswf/go-cli-planner"
)

func TestTasks(t *testing.T) {
	taskList := planner.Tasks{}

	t.Run("A task is added to the task list", func(t *testing.T) {
		taskList.Add("Task 1")

		got := taskList.GetAll()
		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("A task is marked as completed", func(t *testing.T) {
		err := taskList.Complete(1)
		if err != nil {
			t.Fatal(err)
		}

		got := taskList.GetAll()
		want := []planner.Task{{Id: 1, Name: "Task 1", Complete: true}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("Test error when task id doesn't exist", func(t *testing.T) {
		err := taskList.Complete(100)

		if err == nil {
			t.Errorf("Expected an error")
		}
	})

	t.Run("Test only show incomplete", func(t *testing.T) {
		taskList.Add("Task 2")
		taskList.Add("Task 3")
		taskList.Add("Task 4")
		err := taskList.Complete(3)
		if err != nil {
			t.Fatal(err)
		}

		got := taskList.GetOutstanding()
		want := []planner.Task{
			{Id: 2, Name: "Task 2", Complete: false},
			{Id: 4, Name: "Task 4", Complete: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
