package queue

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	UUID       string    `json:"uuid"`
	Type       string    `json:"type"`
	Payload    []byte    `json:"payload"`
	Status     string    `json:"status"`
	Retries    int       `json:"retries"`
	MaxRetries int       `json:"max_retries"`
	CreatedAt  time.Time `json:"created_at"`
}

type TaskProducer struct {
	Db    *sql.DB
	Queue chan Task
}

func GetTasks(db *sql.DB) ([]Task, error) {
	q := `
	select * from tasks;
	`

	rows, err := db.Query(q)
	if err != nil {
		log.Printf("error gettings tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	var t Task
	for rows.Next() {
		if err := rows.Scan(
			&t.UUID,
			&t.Type,
			&t.Payload,
			&t.Status,
			&t.Retries,
			&t.MaxRetries,
			&t.CreatedAt,
		); err != nil {
			log.Printf("error scanning row: %v", err)
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func PostTask(task Task, db *sql.DB) error {
	q := `
	insert into tasks (uuid, type, payload, status, retries, max_retries)
	values ($1, $2, $3, $4, $5, $6)
	`
	_, err := db.Exec(q, task.UUID, task.Type, task.Payload, task.Status,
		task.Retries, task.MaxRetries)
	if err != nil {
		log.Printf("error posting task: %v", err)
		return err
	}
	return nil
}

func DeleteTask(uuid string, db *sql.DB) error {
	q := `
	delete from tasks where uuid = ?;
	`

	_, err := db.Exec(q, uuid)
	if err != nil {
		log.Printf("error deleting task: %v", err)
		return err
	}
	return nil
}

func initTask(taskType string, payload any, maxRetries int) (Task, error) {
	var task Task
	if taskType == "" || payload == nil {
		return task, errors.New("taskType and payload must be initialized")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return task, err
	}
	task = Task{
		UUID:       uuid.NewString(),
		Type:       taskType,
		Payload:    data,
		Status:     "pending",
		Retries:    0,
		MaxRetries: maxRetries,
	}

	return task, nil
}

func (t *TaskProducer) Enqueue(taskType string, payload any, maxRetries int) error {
	if t.Queue == nil {
		t.Queue = make(chan Task, 100)
	}

	task, err := initTask(taskType, payload, maxRetries)
	if err != nil {
		return err
	}

	err = PostTask(task, t.Db)
	if err != nil {
		return err
	}

	t.Queue <- task

	return nil
}
