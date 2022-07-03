package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	planner "github.com/rosswf/go-cli-planner"
	storage "github.com/rosswf/go-cli-planner/storage"
)

type model struct {
	tasks       []planner.Task
	taskStorage planner.TaskList
	cursor      int
	taskInput   string
	toggle      bool
}

func initialModel(taskList *planner.TaskList) model {
	tasks, _ := taskList.GetAll()
	return model{
		tasks:       tasks,
		taskStorage: *taskList,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			if m.taskInput != "" {
				m.taskStorage.Add(m.taskInput)
				m.taskInput = ""
			}

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}
		case "left", "right":
			if len(m.tasks) == 0 {
				break
			}
			if err := m.taskStorage.ToggleStatus(&m.tasks[m.cursor]); err != nil {
				fmt.Println(err)
			}
		case "tab":
			if m.toggle {
				m.toggle = false
			} else {
				m.toggle = true
			}

		case "backspace":
			if m.taskInput != "" {
				m.taskInput = m.taskInput[:len(m.taskInput)-1]
			}
		default:
			m.taskInput += msg.String()
		}

	}
	if m.toggle {
		m.tasks, _ = m.taskStorage.GetAll()
	} else {
		m.tasks, _ = m.taskStorage.GetOutstanding()
	}
	if m.cursor > len(m.tasks)-1 {
		m.cursor = 0
	}
	return m, nil
}

func (m model) View() string {
	s := "Here is your task list:\n\n"

	for i, choice := range m.tasks {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		complete := "✖"
		if choice.Complete {
			complete = "✓"
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, complete, choice.Name)
	}

	s += fmt.Sprintf("\nAdd a new task > %s█", m.taskInput)
	s += "\nPress tab to mark completed. Press ctrl+c to quit.\n"

	return s
}

func main() {
	storage, _ := storage.CreateSqlite3TaskStorage("tasks.db")

	taskList := planner.CreateTaskList(storage)

	p := tea.NewProgram(initialModel(taskList))
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
