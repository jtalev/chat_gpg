package task_queue

import (
	"context"
	"database/sql"
	"log"
)

type Worker struct {
	ctx context.Context
	db  *sql.DB
}

func InitWorker(ctx context.Context, db *sql.DB) Worker {
	return Worker{
		ctx: ctx,
		db:  db,
	}
}

func (w *Worker) StartWorkerLoop(queue chan Task) {
	go func() {
		err := w.StartWorker(queue)
		if err != nil {
			log.Printf("worker stopped running: %v", err)
			go w.StartWorker(queue)
		}
	}()
}

func (w *Worker) StartWorker(queue chan Task) error {
	for {
		select {
		case <-w.ctx.Done():
			log.Println("worker shutting down")
			return nil
		case task := <-queue:
			if task.Type == "scheduled" {
				log.Println("scheduled worker")
			}
			if task.Type == "one_time" {
				log.Println("one_time worker")
			}
			switch task.Type {
			case "one_time":
				if task.Retries >= task.MaxRetries {
					DeleteTask(task.UUID, w.db)
					continue
				}

				handler, ok := HandlerRegistry[task.Handler]
				if !ok {
					log.Println("handler required for this task not supported")
					continue
				}

				err := handler.ProcessTask(task, queue, w.db)
				if err != nil {
					log.Printf("error processing task: %v", err)
					return err
				}
			case "scheduled":
				scheduler := initScheduler(w.db, queue)
				scheduler.startSchedulerLoop(task)
			default:
				log.Printf("%s task type not supported", task.Type)
			}
		}
	}
}
