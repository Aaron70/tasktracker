package tui

import (
	"strings"

	"github.com/aaron70/task/models"
	"github.com/aaron70/task/services"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TasksTable struct {
	Service services.TaskService
	tableModel table.Model
}

func NewTaskTable(service services.TaskService) TasksTable {
	model := TasksTable{ Service: service }

	tasks, err := service.List()
	if err != nil {
		panic(err)
	}

	columns := []table.Column{
		{Title: "Name", Width: 55},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 15},
		{Title: "Created At", Width: 20},
		{Title: "Started At", Width: 20},
		{Title: "Finished At", Width: 20},
		{Title: "Tags", Width: 20},
	}

	rows := make([]table.Row, 0, len(tasks))
	for _, task := range tasks {

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
			task.Name, 
			string(task.Status), 
			task.Duration.String(), 
			createdAt,
			startedAt, 
			finishedAt, 
			strings.Join(tags, "|"),
		}
		rows = append(rows, row)
	}

	model.tableModel = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows) + 1),
	)

	return model
}

func (t TasksTable) Init() tea.Cmd {
	return nil
}

func (t TasksTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return t, tea.Quit
		case "up", "k":
			t.tableModel.MoveUp(1)
		case "down", "j":
			t.tableModel.MoveDown(1)
		}
	}

	return t, nil
}

func (t TasksTable) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	return style.Render(t.View())
}


func PrintTable(tasks []models.Task) string {
	columns := []table.Column{
		{Title: "Name", Width: 55},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 15},
		{Title: "Created At", Width: 20},
		{Title: "Started At", Width: 20},
		{Title: "Finished At", Width: 20},
		{Title: "Tags", Width: 20},
	}

	rows := make([]table.Row, 0, len(tasks))
	for _, task := range tasks {

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
			task.Name, 
			string(task.Status), 
			task.Duration.String(), 
			createdAt,
			startedAt, 
			finishedAt, 
			strings.Join(tags, "|"),
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows) + 1),
	)

	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)

	m := t.View()
	return style.Render(m)
}
