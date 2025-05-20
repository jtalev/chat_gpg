package task_queue

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	UUID       string    `json:"uuid"`
	Type       string    `json:"type"`    // currently supporting 'one_off' and 'scheduled'
	Handler    string    `json:"handler"` // eg. 'db_backup', 'send_email'
	Payload    []byte    `json:"payload"`
	Status     string    `json:"status"`
	Retries    int       `json:"retries"`
	MaxRetries int       `json:"max_retries"`
	CreatedAt  time.Time `json:"created_at"`
}

func initTask(taskType, handler string, payload any) (Task, error) {
	var task Task
	data, err := json.Marshal(payload)
	if err != nil {
		return task, err
	}
	task = Task{
		UUID:       uuid.NewString(),
		Type:       taskType,
		Handler:    handler,
		Payload:    data,
		Retries:    0,
		MaxRetries: 3,
	}

	return task, nil
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
			&t.Handler,
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
	insert into tasks (uuid, type, handler, payload, status, retries, max_retries)
	values ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := db.Exec(q, task.UUID, task.Type, task.Handler, task.Payload, task.Status,
		task.Retries, task.MaxRetries)
	if err != nil {
		log.Printf("error posting task: %v", err)
		return err
	}
	return nil
}

func UpdateTask(task Task, db *sql.DB) error {
	q := `
	update tasks 
	set type = $1, handler = $2, payload = $3, status = $4,
	retries = $5, max_retries = $6
	where uuid = $7;
	`
	_, err := db.Exec(q, task.Type, task.Handler, task.Payload,
		task.Status, task.Retries, task.MaxRetries, task.UUID)
	if err != nil {
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
