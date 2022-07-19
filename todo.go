package todo

type TaskStorage interface {
	Add(*Task) error
	GetAll() ([]Task, error)
	GetTask(TaskId) (*Task, error)
	ToggleStatus(TaskId) error
	GetOutstanding() ([]Task, error)
	Delete(TaskId) error
}

type TaskId int64

type Task struct {
	Id       TaskId `json:"id"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

type TaskList struct {
	storage TaskStorage
}

func CreateTaskList(storage TaskStorage) *TaskList {
	return &TaskList{storage: storage}
}

func (t *TaskList) Add(name string) error {
	task := Task{Name: name, Complete: false}
	err := t.storage.Add(&task)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskList) GetAll() ([]Task, error) {
	return t.storage.GetAll()
}

func (t *TaskList) ToggleStatus(task *Task) error {
	err := t.storage.ToggleStatus(task.Id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskList) GetOutstanding() ([]Task, error) {
	return t.storage.GetOutstanding()
}

func (t *TaskList) Delete(task *Task) error {
	err := t.storage.Delete(task.Id)
	return err
}

func (t *TaskList) GetOne(id TaskId) (Task, error) {
	task, err := t.storage.GetTask(id)

	if err != nil {
		return Task{}, err
	}

	return *task, nil
}
