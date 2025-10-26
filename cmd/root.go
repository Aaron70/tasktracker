package cmd

import (
	"github.com/aaron70/task/services"
	"github.com/spf13/cobra"
)

func newRootCommand(taskService services.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use: "task",
		Short: "CLI tool to keep track of daily tasks.",
		Long: "Task is a CLI tool that helps you to keep track of the daily tasks and the time you spend working on the task.",
		Example: "task --help",
	}

	cmd.AddCommand(newNewTaskCommand(taskService))
	cmd.AddCommand(newListTasksCommand(taskService))
	cmd.AddCommand(newSwitchTaskCommand(taskService))
	cmd.AddCommand(newStopTaskCommand(taskService))

	return cmd
}

func Run(taskService services.TaskService) error {
	cmd := newRootCommand(taskService)

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}

