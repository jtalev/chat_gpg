package task_queue

import (
	"database/sql"
	"log"
	"time"
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

func UpdateTask(task Task, db *sql.DB) error {
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
