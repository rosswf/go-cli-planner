package planner

import "errors"

type Task struct {
	Name     string
	Complete bool
}

type Tasks []Task

func (t *Tasks) Add(name string) {
	*t = append(*t, Task{Name: name, Complete: false})
}

func (t Tasks) Complete(id int) error {
	id-- // slice is zero indexed

	if id > len(t) {
		return errors.New("No task with that id")
	}

	t[id].Complete = true
	return nil
}

func (t Tasks) GetAll() []Task {
	return t
}
