package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type TaskStatus string

const (
	TODO        TaskStatus = "todo"
	IN_PROGRESS TaskStatus = "inprogress"
	DONE        TaskStatus = "done"
	WONTDO      TaskStatus = "wontdo"
)

type Task struct {
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	Status     TaskStatus    `json:"status"`
	Duration   time.Duration `json:"duration"`
	CreatedAt  time.Time     `json:"createdAt"`
	StartedAt  time.Time     `json:"startedAt"`
	FinishedAt time.Time     `json:"finishedAt"`
	InProgress time.Time     `json:"inProgress"`
	Tags       []Tag         `json:"tags"`
}

func NewTask(name string) Task {
	createdAt := time.Now()
	name = strings.TrimSpace(name)
	return Task{
		Id:        HashID(name, createdAt),
		Name:      name,
		Status:    TODO,
		CreatedAt: createdAt,
	}
}

func HashID(name string, createdAt time.Time) string {
	data := fmt.Sprintf("%s%s", strings.ToLower(strings.TrimSpace(name)), createdAt.Format("02/01/2006"))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
