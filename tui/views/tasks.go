package views

import (
	"time"

	"github.com/aaron70/task/services"
	"github.com/aaron70/task/tui/components"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanefung/bubble-datepicker"
)

var style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)

type TasksModel struct {
	date        time.Time
	taskService services.TaskService

	tableModel      table.Model
}

func NewTasksModel(taskService services.TaskService) (TasksModel, error) {
	day := time.Time{}
	tasks, err := taskService.Find(services.FindTasksFilter{Day: day})
	if err != nil {
		return TasksModel{}, err
	}

	tableModel := components.NewTasksTable(components.ConvertTaskToRows(tasks)...)
	datepickerModel := datepicker.New(time.Now())
	datepickerModel.SelectDate()
	return TasksModel{
		date:            day,
		taskService:     taskService,
		tableModel:      tableModel,
	}, nil
}

func (t TasksModel) Init() tea.Cmd {
	return nil
}

func (t TasksModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return t.updateKeyMsg(msg)
	}
	return t, cmd
}

func (t TasksModel) updateKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m, cmd := t.tableModel.Update(msg)
	t.tableModel = m
	switch msg.String() {
	case "ctrl+c", "q":
		return t, tea.Quit
	case "s":
		_, _, err := t.taskService.Switch(t.getSelectedNameDate())
		if err != nil {
			panic(err)
		}
		t.refreshTasks()
	case "S":
		err := t.taskService.Start(t.getSelectedNameDate())
		if err != nil {
			panic(err)
		}
		t.refreshTasks()
	case "d":
		err := t.taskService.Stop(t.getSelectedNameDate())
		if err != nil {
			panic(err)
		}
		t.refreshTasks()
	case "D":
		err := t.taskService.Delete(t.getSelectedNameDate())
		if err != nil {
			panic(err)
		}
		t.refreshTasks()
	default:

	}
	return t, cmd
}

func (t TasksModel) getSelectedNameDate() (string, time.Time) {
	row := t.tableModel.SelectedRow()

	if len(row) <= 0 {
		return "", time.Time{}
	}

	createdAt, err := time.Parse("02-01-2006 15:04:05", row[4])
	if err != nil {
		panic(err)
	}
	return row[1], createdAt
}

func (t *TasksModel) refreshTasks() {
	tasks, err := t.taskService.Find(services.FindTasksFilter{Day: t.date})
	if err != nil {
		panic(err)
	}
	t.tableModel.SetRows(components.ConvertTaskToRows(tasks))
}

func (t TasksModel) View() string {
	task, _ := t.getSelectedNameDate()
	table := style.Render(t.tableModel.View()) + "\n"

	view := ""
	view += table + "\n"
	view += "Task: " + task + "\n"

	return view
}
