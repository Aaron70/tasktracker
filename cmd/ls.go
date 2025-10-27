package cmd

import (
	"time"

	"github.com/aaron70/task/models"
	"github.com/aaron70/task/services"
	"github.com/aaron70/task/tui"
	"github.com/spf13/cobra"
)

func newListTasksCommand(taskServices services.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use: "ls",
		Short: "List the tasks started in specific day and those in status 'InProgress', 'ToDo'.",
		RunE: func(cmd *cobra.Command, args []string) error {

			all, err := cmd.Flags().GetBool("all")
			if err != nil {
				return err
			}

			if all {
				allTasks, err :=  taskServices.List()
				if err != nil { return err }
				tui.PrintTable(allTasks)
				return nil
			}

			day, err := cmd.Flags().GetTime("date")
			if err != nil { return err }

			status, err := cmd.Flags().GetStringArray("status")
			if err != nil { return err }
			filter := services.FindTasksFilter{
				Day: day,
				Status: status,
			}

			tasks, err := taskServices.Find(filter)
			if err != nil {
				return err
			}

			tui.PrintTable(tasks)
			return nil
		},
	}

	cmd.Flags().BoolP("all", "a", false, "List all tasks without any filter applied.")
	cmd.Flags().TimeP("date", "d", time.Now(), DateFormats , "The date to filter the tasks by the StartedAt field.")
	cmd.Flags().StringArrayP("status", "s", []string{ string(models.TODO), string(models.IN_PROGRESS) }, "The status to filter the tasks by")
	// TODO: Add support for filter by tags
	// cmd.Flags().StringArrayP("tag", "t", []string{}, "The tag to filter the tasks by")

	return cmd
}
