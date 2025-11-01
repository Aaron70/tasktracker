package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/aaron70/task/models"
	"github.com/aaron70/task/services"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableTask struct {
	service    services.TaskService
	tableModel table.Model
}

func NewTableTask(service services.TaskService) (TableTask, error) {
	tasks, err := service.List()
	if err != nil {
		return TableTask{}, err
	}
	columns := []table.Column{
		{Title: "N", Width: 5},
		{Title: "Name", Width: 55},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 15},
		{Title: "Created At", Width: 20},
		{Title: "Started At", Width: 20},
		{Title: "Finished At", Width: 20},
		{Title: "Tags", Width: 20},
	}

	rows := tasksToRows(tasks)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	return TableTask{
		service:    service,
		tableModel: t,
	}, nil
}

func (t TableTask) Init() tea.Cmd {
	return nil
}

func (t TableTask) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return t, tea.Quit
		case "up", "k":
			t.tableModel.MoveUp(1)
		case "down", "j":
			t.tableModel.MoveDown(1)
		case "s":
			name, date := t.getSelectedRowNameDate()
			_, _, err := t.service.Switch(name, date)
			if err != nil {
				panic(err)
			}
		case "S":
			name, date := t.getSelectedRowNameDate()
			err := t.service.Start(name, date)
			if err != nil {
				panic(err)
			}
		case "d":
			name, date := t.getSelectedRowNameDate()
			err := t.service.Stop(name, date)
			if err != nil {
				panic(err)
			}
		case "D":
			name, date := t.getSelectedRowNameDate()
			err := t.service.Delete(name, date)
			if err != nil {
				panic(err)
			}
		}
	}

	tasks, err := t.service.List()
	if err != nil {
		panic(err)
	}
	t.tableModel.SetRows(tasksToRows(tasks))
	return t, nil
}

func (t TableTask) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	view := strings.Builder{}
	view.WriteString(style.Render(t.tableModel.View()) + "\n")
	view.WriteString("Task: " + t.tableModel.SelectedRow()[1] + "\n")
	return view.String()
}

func (t TableTask) getSelectedRowNameDate() (string, time.Time) {
	row := t.tableModel.SelectedRow()
	createdAt, err := time.Parse("02-01-2006 15:04:05", row[4])
	if err != nil {
		panic(err)
	}
	return row[1], createdAt
}

func taskToRow(index int, task models.Task) table.Row {

	finishedAt := "---"
	if !task.FinishedAt.IsZero() {
		finishedAt = task.FinishedAt.Format("02-01-2006 15:04:05")
	}

	startedAt := "---"
	if !task.StartedAt.IsZero() {
		startedAt = task.StartedAt.Format("02-01-2006 15:04:05")
	}

	createdAt := "---"
	if !task.CreatedAt.IsZero() {
		createdAt = task.CreatedAt.Format("02-01-2006 15:04:05")
	}

	tags := make([]string, len(task.Tags))
	for i, tag := range task.Tags {
		tags[i] = tag.Name
	}

	row := table.Row{
		fmt.Sprint(index),
		task.Name,
		string(task.Status),
		task.Duration.String(),
		createdAt,
		startedAt,
		finishedAt,
		strings.Join(tags, "|"),
	}
	return row
}

func tasksToRows(tasks []models.Task) []table.Row {
	rows := make([]table.Row, 0, len(tasks))
	for i, task := range tasks {
		rows = append(rows, taskToRow(i, task))
	}
	return rows
}
