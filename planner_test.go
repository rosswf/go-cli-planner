package planner_test

import (
	"reflect"
	"testing"

	planner "github.com/rosswf/go-cli-planner"
)

func TestTasks(t *testing.T) {
	t.Run("A task is added to the task list", func(t *testing.T) {
		taskList := planner.Tasks{}
		taskList.Add("Task 1")

		got := taskList.GetAll()
		want := []planner.Task{{Name: "Task 1", Complete: false}}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}
