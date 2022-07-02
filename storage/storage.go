package planner_storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	planner "github.com/rosswf/go-cli-planner"
)

type Sqlite3TaskStorage struct {
	conn *sql.DB
}

func CreateSqlite3TaskStorage(location string) (*Sqlite3TaskStorage, error) {
	db, _ := sql.Open("sqlite3", location) // Come back to this future me
	sqlStmt := `CREATE TABLE IF NOT EXISTS tasks (id INTEGER not null primary key, name TEXT, complete BOOL);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}
	return &Sqlite3TaskStorage{db}, nil
}

func (s *Sqlite3TaskStorage) Close() {
	s.conn.Close()
}

func (s *Sqlite3TaskStorage) Add(task *planner.Task) error {
	tx, _ := s.conn.Begin()
	stmt, _ := tx.Prepare("INSERT INTO tasks(name, complete) values(?, ?)")
	defer stmt.Close()
	result, _ := stmt.Exec(task.Name, task.Complete)
	err := tx.Commit()
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	task.Id = planner.TaskId(id)
	return err
}

func (s *Sqlite3TaskStorage) GetAll() ([]planner.Task, error) {
	tasks := []planner.Task{}
	rows, err := s.conn.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id planner.TaskId
		var name string
		var complete bool
		_ = rows.Scan(&id, &name, &complete)
		tasks = append(tasks, planner.Task{Id: id, Name: name, Complete: complete})
	}
	return tasks, nil
}

func (s *Sqlite3TaskStorage) GetTask(id planner.TaskId) (*planner.Task, error) {
	row := s.conn.QueryRow("SELECT * FROM tasks WHERE id = ?", id)

	var taskId planner.TaskId
	var name string
	var complete bool
	err := row.Scan(&taskId, &name, &complete)
	if err != nil {
		return nil, err
	}
	return &planner.Task{Id: taskId, Name: name, Complete: complete}, nil
}

func (s *Sqlite3TaskStorage) ToggleStatus(id planner.TaskId) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE tasks SET complete = CASE WHEN complete = true THEN false ELSE true END WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sqlite3TaskStorage) GetOutstanding() ([]planner.Task, error) {
	tasks := []planner.Task{}
	rows, err := s.conn.Query("SELECT * FROM tasks WHERE complete = false")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id planner.TaskId
		var name string
		var complete bool
		_ = rows.Scan(&id, &name, &complete)
		tasks = append(tasks, planner.Task{Id: id, Name: name, Complete: complete})
	}
	return tasks, nil
}

func (s *Sqlite3TaskStorage) Delete(id planner.TaskId) error {
	tx, err := s.conn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM tasks WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
