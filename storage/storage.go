package todo_storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	todo "github.com/rosswf/go-todo"
)

type Sqlite3TaskStorage struct {
	conn *sql.DB
}

func CreateSqlite3TaskStorage(location string) (*Sqlite3TaskStorage, error) {
	db, err := sql.Open("sqlite3", location)
	if err != nil {
		return nil, err
	}
	sqlStmt := `CREATE TABLE IF NOT EXISTS tasks
(id INTEGER not null primary key, name TEXT, complete BOOL);`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	return &Sqlite3TaskStorage{db}, nil
}

func (s *Sqlite3TaskStorage) Close() {
	s.conn.Close()
}

func (s *Sqlite3TaskStorage) Add(task *todo.Task) (todo.TaskId, error) {
	sqlStmt := "INSERT INTO tasks(name, complete) values(?, ?)"
	result, err := s.conn.Exec(sqlStmt, task.Name, task.Complete)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return todo.TaskId(id), err
}

func (s *Sqlite3TaskStorage) GetAll() ([]todo.Task, error) {
	tasks := []todo.Task{}
	rows, err := s.conn.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id todo.TaskId
		var name string
		var complete bool
		_ = rows.Scan(&id, &name, &complete)
		tasks = append(tasks, todo.Task{Id: id, Name: name, Complete: complete})
	}
	return tasks, nil
}

func (s *Sqlite3TaskStorage) GetTask(id todo.TaskId) (*todo.Task, error) {
	row := s.conn.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var taskId todo.TaskId
	var name string
	var complete bool
	err := row.Scan(&taskId, &name, &complete)
	if err != nil {
		return nil, err
	}
	return &todo.Task{Id: taskId, Name: name, Complete: complete}, nil
}

func (s *Sqlite3TaskStorage) ToggleStatus(id todo.TaskId) error {
	sqlStmt := `UPDATE tasks SET complete = CASE WHEN complete = true 
THEN false ELSE true END WHERE id=?`

	_, err := s.conn.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sqlite3TaskStorage) GetOutstanding() ([]todo.Task, error) {
	tasks := []todo.Task{}
	rows, err := s.conn.Query("SELECT * FROM tasks WHERE complete = false")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id todo.TaskId
		var name string
		var complete bool
		err = rows.Scan(&id, &name, &complete)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, todo.Task{Id: id, Name: name, Complete: complete})
	}
	return tasks, nil
}

func (s *Sqlite3TaskStorage) Delete(id todo.TaskId) error {
	sqlStmt := "DELETE FROM tasks WHERE id=?"
	_, err := s.conn.Exec(sqlStmt, id)
	if err != nil {
		return err
	}

	return nil
}
