package task_queue

import (
	"context"
	"database/sql"
	"log"
)

func StartWorker(ctx context.Context, taskQueue chan Task, db *sql.DB) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("worker shutting down")
			return nil
		case task := <-taskQueue:
			if task.Status == "failed" || task.Status == "complete" {
				continue
			}
			payload := task.Payload
			handler := HandlerRegistry[task.Type]
			err := handler(payload)
			if err != nil {
				p := TaskProducer{
					Db:    db,
					Queue: taskQueue,
				}
				task.Retries++
				if task.Retries > task.MaxRetries {
					task.Status = "failed"
					err = UpdateTask(task, db)
					if err != nil {
						return err
					}
				} else {
					err = p.Retry(task)
					if err != nil {
						return err
					}
				}
			}
			task.Status = "complete"
			err = UpdateTask(task, db)
			if err != nil {
				return err
			}
		}
	}
}
