package repositories

import (
	"errors"

	"github.com/aaron70/task/models"
)

var TaskNotFoundError = errors.New("TaskNotFoundError")

type TaskRepository interface {
	Save(task models.Task) error
	Update(id string, task models.Task) error
	Get(id string) (models.Task, error)
	GetAll() ([]models.Task, error)
	Delete(id string) error
	SelectTask(models.Task) error
	GetSelectedTask() (models.Task, error)
}

