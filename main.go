package main

import (
	"os"

	"github.com/aaron70/task/cmd"
	repositories "github.com/aaron70/task/respository"
	"github.com/aaron70/task/services"
)

func main() {
	taskRepository := repositories.NewFsTaskRepository("./")
	taskService := services.NewTaskService(taskRepository)

	err := cmd.Run(taskService)
	if err != nil {
		os.Exit(1)
	}
}
