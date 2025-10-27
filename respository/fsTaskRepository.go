package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/aaron70/task/models"
)

type fsTaskRepository struct {
	Path string
}

func NewFsTaskRepository(path string) fsTaskRepository {
	return fsTaskRepository{ Path: path }
}

func (r fsTaskRepository) Save(task models.Task) error {
	tasks, err := r.GetAll()
	if err != nil { return err }

	index := slices.IndexFunc(tasks, func(t models.Task) bool { return t.Id == task.Id })
	if index == -1 {
		tasks = append(tasks, task)
	} else {
		return fmt.Errorf("The task %q already exists", task.Id)
	}

	return saveTasks(filepath.Join(r.Path, "tasks.json"), tasks)
}

func (r fsTaskRepository) Update(id string, task models.Task) error {
	tasks, err := r.GetAll()
	if err != nil { return err }

	index := slices.IndexFunc(tasks, func(t models.Task) bool { return t.Id == task.Id })
	if index == -1 {
		return fmt.Errorf("%w: The task %q wasn't found", TaskNotFoundError, task.Id)
	} else {
		task.Name = tasks[index].Name
		task.CreatedAt = tasks[index].CreatedAt
		tasks[index] = task
	}

	saveTasks(filepath.Join(r.Path, "tasks.json"), tasks)

	return nil
}

func (r fsTaskRepository) Get(id string) (models.Task, error) {
	tasks, err := r.GetAll()
	if err != nil {
		return models.Task{}, err
	}
	i := slices.IndexFunc(tasks, func(t models.Task) bool { return t.Id == id })
	if i <= -1 {
		return models.Task{}, fmt.Errorf("%w: Task %q not found", TaskNotFoundError, id)
	}
	return tasks[i], nil
}

func (r fsTaskRepository) GetAll() ([]models.Task, error) {
	content, err := os.ReadFile(filepath.Join(r.Path, "tasks.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Task(nil), nil
		}
		return []models.Task(nil), err
	}

	var tasks []models.Task
  err = json.Unmarshal(content, &tasks)
  if err != nil {
		return []models.Task(nil), err
  }

	return tasks, nil
}

func (r fsTaskRepository) Delete(id string) error {
	tasks, err := r.GetAll()
	if err != nil {
		return err
	}
	i := slices.IndexFunc(tasks, func(t models.Task) bool { return t.Id == id })
	if i <= -1 {
		return fmt.Errorf("%w: Task %q not found", TaskNotFoundError, id)
	}

	tasks = slices.Delete(tasks, i, i+1)
	saveTasks(filepath.Join(r.Path, "tasks.json"), tasks)

	return nil
}

func (r fsTaskRepository) SelectTask(task models.Task) error {
	return  saveTasks(filepath.Join(r.Path, "selected.json"), []models.Task{ task })
}


func (r fsTaskRepository) GetSelectedTask() (models.Task, error) {
	content, err := os.ReadFile(filepath.Join(r.Path, "selected.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return models.Task{}, nil
		}
		return models.Task{}, err
	}

	var tasks []models.Task
  err = json.Unmarshal(content, &tasks)
  if err != nil {
		return models.Task{}, err
  }

	if len(tasks) <= 0 {
		return models.Task{}, nil
	}

	task := tasks[0]

	task, err = r.Get(models.HashID(task.Name, task.CreatedAt))
	if err != nil {
		if errors.Is(err, TaskNotFoundError) {
			return models.Task{}, nil
		}
		return models.Task{}, err
	}

	return task, nil
}

func saveTasks(path string, tasks []models.Task) error {
	payload, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, payload, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

