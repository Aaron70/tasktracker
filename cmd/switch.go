package cmd

import (
	"fmt"
	"time"

	"github.com/aaron70/task/models"
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

func switchTask(service services.TaskService, name string, date time.Time) error {
	started, stopped, err := service.Switch(name, date)
	if err != nil {
		return err
	}

	fmt.Println("Started Task")
	tui.PrintTable([]models.Task{started})

	if stopped.Id != "" {
		fmt.Println("Stopped Task")
		tui.PrintTable([]models.Task{stopped})
	}

	return nil
}
