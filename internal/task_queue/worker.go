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
			// if task.type = one_time
			handler, ok := HandlerRegistry[task.Handler]
			if !ok {
				log.Println("handler required for this task not supported")
				continue
			}
			log.Println("worker running")
			// if task.Type == scheduled
			// pass task to scheduler
			err := handler.ProcessTask(task, queue, w.db)
			if err != nil {
				log.Printf("error processing task: %v", err)
				return err
			}
		}
	}
}
