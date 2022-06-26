package planner

import "errors"

type Task struct {
	Id       int
	Name     string
	Complete bool
}

type Tasks []Task

func (t *Tasks) Add(name string) {
	nextId := len(*t) + 1
	*t = append(*t, Task{Id: nextId, Name: name, Complete: false})
}

func (t Tasks) ToggleStatus(id int) error {
	id-- // slice is zero indexed

	if id > len(t) {
		return errors.New("No task with that id")
	}

	if t[id].Complete {
		t[id].Complete = false
	} else {
		t[id].Complete = true
	}
	return nil
}

func (t Tasks) GetOutstanding() []Task {
	outstanding := []Task{}
	for _, task := range t {
		if !task.Complete {
			outstanding = append(outstanding, task)
		}
	}
	return outstanding
}
