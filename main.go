package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aaron70/task/cmd"
	repositories "github.com/aaron70/task/respository"
	"github.com/aaron70/task/services"
)

func main() {
	path, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Couldn't get the task tracker folder path", err)
		os.Exit(1)
	}

	folderPath := filepath.Join(path, ".tasktracker")

	if _, err := os.Stat(folderPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(folderPath, os.ModePerm)
			if err != nil {
				fmt.Println("Couldn't create the task tracker folder")
				os.Exit(1)
			}
		}
	}

	taskRepository := repositories.NewFsTaskRepository(folderPath)
	taskService := services.NewTaskService(taskRepository)

	err = cmd.Run(taskService)
	if err != nil {
		os.Exit(1)
	}
}
