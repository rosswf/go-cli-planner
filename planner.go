package planner

type Task struct {
	Name     string
	Complete bool
}

type Tasks []Task

func (t *Tasks) Add(name string) {
	*t = append(*t, Task{Name: name, Complete: false})
}

func (t Tasks) GetAll() []Task {
	return t
}
