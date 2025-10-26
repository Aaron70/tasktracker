package cmd

import (
	"github.com/aaron70/task/models"
	"github.com/aaron70/task/services"
	"github.com/aaron70/task/tui"
	"github.com/spf13/cobra"
)

func newNewTaskCommand(taskService services.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Creates a new task with the given name.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			task := models.NewTask(args[0])

			tags, err := cmd.Flags().GetStringArray("tag")
			if err != nil {
				return err
			}

			task.Tags = make([]models.Tag, 0, len(tags))
			for _, tag := range tags {
				task.Tags = append(task.Tags, models.Tag{Name: tag})
			}

			err = taskService.Save(&task)
			if err != nil {
				return err
			}

			shouldStart, err := cmd.Flags().GetBool("start")
			if err != nil {
				return err
			}

			if shouldStart {
				err = taskService.Start(task.Name, task.CreatedAt)
				if err != nil {
					return err
				}
				switchTask(taskService, task.Name, task.CreatedAt)
				return  nil
			}

			task, err = taskService.Get(task.Name, task.CreatedAt)
			if err != nil {
				return err
			}

			tui.PrintTable([]models.Task{task})
			return nil
		},
	}

	cmd.Flags().BoolP("start", "s", false, "Starts the task after creation.")
	cmd.Flags().StringArrayP("tag", "t", []string{}, "Add a tag to the task.")

	return cmd
}
