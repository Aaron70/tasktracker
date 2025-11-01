package services

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aaron70/task/internal/stringutils"
	"github.com/aaron70/task/models"
	repositories "github.com/aaron70/task/respository"
)

type FindTasksFilter struct {
	Day    time.Time
	Status []string
	Tags   []string
}

type TaskService interface {
	Save(*models.Task) error
	Update(string, time.Time, *models.Task) error
	Delete(string, time.Time) error
	Get(string, time.Time) (models.Task, error)
	Find(FindTasksFilter) ([]models.Task, error)
	List() ([]models.Task, error)
	Start(string, time.Time) error
	Stop(string, time.Time) error
	SelectTask(models.Task) error
	GetSelectedTask() (models.Task, error)
	Switch(string, time.Time) (startedTask models.Task, stoppedTask models.Task, err error)
}

func validateTask(task *models.Task) error {
	if stringutils.IsBlank(task.Name) {
		return fmt.Errorf("Task name can not be blank")
	}

	if task.CreatedAt.IsZero() {
		task.CreatedAt = time.Now()
	}

	if stringutils.IsBlank(task.Id) {
		task.Id = models.HashID(task.Name, task.CreatedAt)
	}

	return nil
}

type taskService struct {
	repository repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) taskService {
	return taskService{repository: repo}
}

func (s taskService) Save(task *models.Task) error {
	t, err := s.Get(task.Name, task.CreatedAt)
	if err == nil && t.Name != "" {
		return fmt.Errorf("Task %q already exists", task.Name)
	}

	if err := validateTask(task); err != nil {
		return err
	}
	return s.repository.Save(*task)
}

func (s taskService) Update(name string, date time.Time, task *models.Task) error {
	if name != task.Name {
		return fmt.Errorf("You can't change the name of the task: %s", name)
	}

	if task.CreatedAt.Format("02/01/2006") != date.Format("02/01/2006") {
		return fmt.Errorf("You can't change the creation date of the task: %s", name)
	}

	if err := validateTask(task); err != nil {
		return err
	}

	return s.repository.Update(models.HashID(name, date), *task)
}

func (s taskService) Get(name string, date time.Time) (models.Task, error) {
	task, err := s.repository.Get(models.HashID(name, date))
	if err != nil {
		if errors.Is(err, repositories.TaskNotFoundError) {
			return models.Task{}, fmt.Errorf("%w: Task %q not found from date %q", repositories.TaskNotFoundError, name, date.Format("02-01-2006"))
		}
		return models.Task{}, err
	}

	// TODO: Maybe move this to the ls command and only applies for one task and requires a flag?
	// if !task.InProgress.IsZero() {
	// 	task.Duration += time.Since(task.InProgress)
	// }

	return task, err
}

func (s taskService) List() ([]models.Task, error) {
	return s.repository.GetAll()
}

func (s taskService) Delete(name string, date time.Time) error {
	return s.repository.Delete(models.HashID(name, date))
}

func (s taskService) Start(name string, date time.Time) error {
	task, err := s.Get(name, date)
	if err != nil {
		return err
	}

	now := time.Now()
	if task.StartedAt.IsZero() {
		task.StartedAt = now
	}

	task.FinishedAt = time.Time{}
	task.Status = models.IN_PROGRESS
	task.InProgress = now

	return s.Update(name, date, &task)
}

func (s taskService) Stop(name string, date time.Time) error {
	task, err := s.Get(name, date)
	if err != nil {
		return err
	}

	if !task.InProgress.IsZero() {
		task.Duration += time.Since(task.InProgress)
	}
	task.InProgress = time.Time{}
	task.Status = models.DONE
	task.FinishedAt = time.Now()

	return s.Update(name, date, &task)
}

func (s taskService) Find(filter FindTasksFilter) ([]models.Task, error) {
	allTasks, err := s.List()
	if err != nil {
		return []models.Task(nil), err
	}

	filteredTasks := make([]models.Task, 0, len(allTasks))
	for _, task := range allTasks {
		dayFilter := true
		statusFilter := true

		if !filter.Day.IsZero() {
			dayFilter = task.StartedAt.Format("02/01/2006") == filter.Day.Format("02/01/2006")
		}

		if len(filter.Status) > 0 {
			statusFilter = slices.Contains(filter.Status, string(task.Status))
		}

		if len(filter.Tags) > 0 {
			// TODO: Filter tasks by their tags
		}

		if dayFilter && statusFilter {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return filteredTasks, nil
}

func (s taskService) SelectTask(task models.Task) error {
	return s.repository.SelectTask(task)
}

func (s taskService) GetSelectedTask() (models.Task, error) {
	return s.repository.GetSelectedTask()
}

func (s taskService) Switch(name string, date time.Time) (startedTask, stoppedTask models.Task, err error) {
	err = s.Start(name, date)
	if err != nil {
		return models.Task{}, models.Task{}, err
	}

	startedTask, err = s.Get(name, date)
	if err != nil {
		return models.Task{}, models.Task{}, err
	}

	stoppedTask, err = s.GetSelectedTask()
	if err != nil {
		return startedTask, models.Task{}, err
	}

	err = s.SelectTask(startedTask)
	if err != nil {
		return startedTask, stoppedTask, err
	}

	if stoppedTask.Name != "" && stoppedTask.Name != startedTask.Name {
		err = s.Stop(stoppedTask.Name, stoppedTask.CreatedAt)
		if err != nil {
			return startedTask, stoppedTask, err
		}
	}

	return startedTask, stoppedTask, nil
}

// 	startedTask, err = s.GetSelectedTask()
// 	if err != nil {
// 		return models.Task{}, models.Task{}, err
// 	}
//
// 	if startedTask.Name != "" {
// 		err := s.Stop(startedTask.Name, startedTask.CreatedAt)
// 		if err != nil {
// 			return models.Task{}, models.Task{}, err
// 		}
// 	}
//
// 	if models.HashID(name, date) == models.HashID(startedTask.Name, startedTask.CreatedAt) {
// 		return startedTask, models.Task{}, nil
// 	}
//
// 	err = s.Start(name, date)
// 	if err != nil {
// 		return models.Task{}, models.Task{}, err
// 	}
//
// 	task, err := s.Get(name, date)
// 	if err != nil {
// 		return models.Task{}, models.Task{}, err
// 	}
// 	err = s.SelectTask(task)
// 	if err != nil {
// 		return models.Task{}, models.Task{}, err
// 	}
//
// 	stoppedTask, err = s.Get(startedTask.Name, startedTask.CreatedAt)
// 	if err != nil {
// 		if errors.Is(err, repositories.TaskNotFoundError) {
// 			return models.Task{}, nil
// 		}
// 		return models.Task{}, models.Task{}, err
// 	}
//
// 	return stoppedTask, nil
// }
