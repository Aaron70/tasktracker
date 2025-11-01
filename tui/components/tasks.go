package components

import (
	"fmt"
	"strings"

	"github.com/aaron70/task/models"
	"github.com/charmbracelet/bubbles/table"
)

func NewTasksTable(rows ...table.Row) table.Model {
	columns := []table.Column{
		{Title: "N", Width: 5},
		{Title: "Name", Width: 45},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 15},
		{Title: "Created At", Width: 20},
		{Title: "Started At", Width: 20},
		{Title: "Finished At", Width: 20},
		{Title: "Tags", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(len(rows)+1),
	)

	return t
} 

func ConvertTaskToRow(index int, task models.Task) table.Row {

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

func ConvertTaskToRows(tasks []models.Task) []table.Row {
	rows := make([]table.Row, 0, len(tasks))
	for i, task := range tasks {
		rows = append(rows, ConvertTaskToRow(i, task))
	}
	return rows
}
