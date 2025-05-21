package task_queue

import (
	"database/sql"
	"log"
	"time"
)

type Scheduler struct {
	db    *sql.DB
	queue chan Task
}

func initScheduler(db *sql.DB, queue chan Task) *Scheduler {
	return &Scheduler{
		db:    db,
		queue: queue,
	}
}

func (s *Scheduler) startSchedulerLoop(task Task) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handler, ok := HandlerRegistry[task.Handler]
				if !ok {
					log.Printf("%s handler not found", task.Handler)
					err := DeleteTask(task.UUID, s.db)
					if err != nil {
						log.Printf("error deleting task: %v", err)
						return
					}
					return
				}
				err := handler.ProcessTask(task, s.queue, s.db)
				if err != nil {
					log.Printf("error processing scheduled task with handler %s: %v", task.Handler, err)
					return
				}
			}
		}
	}()
}
