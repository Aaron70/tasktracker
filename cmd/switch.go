package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/aaron70/task/models"
	repositories "github.com/aaron70/task/respository"
	"github.com/aaron70/task/services"
	"github.com/aaron70/task/tui"
	"github.com/spf13/cobra"
)

func newSwitchTaskCommand(taskService services.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "switch",
		Short:   "Starts the given task. Stops the previous InProgress task if any.",
		Long:    "This command will start to keep track of the time spend in the given task. If there is other task InProgress it will stopped and mark as DONE. (You can start it later again)",
		Aliases: []string{"start"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			date, err := cmd.Flags().GetTime("date")
			if err != nil {
				return err
			}

			return switchTask(taskService, args[0], date)
		},
	}

	cmd.Flags().TimeP("date", "d", time.Now(), DateFormats, "The date when the task was created. Defaults to 'today'.")

	return cmd
}

func switchTask(taskService services.TaskService, name string, date time.Time) error {
	current, err := taskService.GetSelectedTask()
	if err != nil {
		return err
	}

	if current.Name != "" {
		err := taskService.Stop(current.Name, current.CreatedAt)
		if err != nil {
			return err
		}
	}

	if models.HashID(name, date) == models.HashID(current.Name, current.CreatedAt) {
		fmt.Printf("The task %q is already started.\n", name)
		return nil
	}

	err = taskService.Start(name, date)
	if err != nil {
		return err
	}

	task, err := taskService.Get(name, date)
	if err != nil {
		return err
	}
	err = taskService.SelectTask(task)
	if err != nil {
		return err
	}

	fmt.Println("Task started:")
	tui.PrintTable([]models.Task{task})

	oldTask, err := taskService.Get(current.Name, current.CreatedAt)
	if err != nil {
		if errors.Is(err, repositories.TaskNotFoundError) {
			return nil
		}
		return err
	}

	fmt.Println("Task stopped:")
	tui.PrintTable([]models.Task{oldTask})

	return nil
}
