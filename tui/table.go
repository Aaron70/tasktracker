package tui

import (
	"fmt"
	"strings"

	"github.com/aaron70/task/models"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func PrintTable(tasks []models.Task) {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 20},
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

		tags := make([]string, len(task.Tags))
		for i, tag := range task.Tags {
			tags[i] = tag.Name
		}

		row := table.Row{ 
			task.Name, 
			string(task.Status), 
			task.Duration.String(), 
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
	fmt.Println(style.Render(m))
}
