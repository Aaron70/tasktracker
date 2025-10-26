package cmd

import (
	"fmt"
	"time"

	"github.com/aaron70/task/models"
	"github.com/aaron70/task/services"
	"github.com/aaron70/task/tui"
	"github.com/spf13/cobra"
)

func newStopTaskCommand(taskService services.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stops a task and mark it as done. But you can start it again.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			date, err := cmd.Flags().GetTime("date")
			if err != nil {
				return err
			}

			err = taskService.Stop(args[0], date)
			if err != nil {
				return err
			}

			task, err := taskService.Get(args[0], date)
			if err != nil {
				return err
			}

			fmt.Println("Task stopped:")
			tui.PrintTable([]models.Task{task})

			return nil
		},
	}

	cmd.Flags().TimeP("date", "d", time.Now(), []string{"02/01/2006"}, "The date when the task was created. Defaults to 'today'.")

	return cmd
}
