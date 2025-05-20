package task_queue

import (
	"database/sql"
	"errors"
	"log"
)

type TaskProducer struct {
	Db             *sql.DB
	OneTimeQueue   chan<- Task
	ScheduledQueue chan<- Task

	enqueueStrategy EnqueueStrategy
}

func (t *TaskProducer) InitQueues() error {
	tasks, err := GetTasks(t.Db)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		switch task.Type {
		case "one_time":
			if task.Status != "pending" {
				continue
			}
			t.OneTimeQueue <- task
		case "scheduled":
			t.ScheduledQueue <- task
		default:
			log.Println("task type not supported")
		}
	}

	return nil
}

func (t *TaskProducer) Enqueue(taskType, handler string, payload any) error {
	if taskType == "" || handler == "" || payload == nil {
		return errors.New("task not initialized correctly")
	}

	switch taskType {
	case "one_time":
		t.enqueueStrategy = &OneTimeEnqueueStrategy{db: t.Db}
		err := t.enqueueStrategy.Enqueue(taskType, handler, payload, t.OneTimeQueue)
		if err != nil {
			return err
		}
	case "scheduled":
		t.enqueueStrategy = &ScheduledEnqueueStrategy{db: t.Db}
		err := t.enqueueStrategy.Enqueue(taskType, handler, payload, t.ScheduledQueue)
		if err != nil {
			return err
		}
	default:
		return errors.New("enqueueing for this task type is not supported")
	}

	return nil
}

type EnqueueStrategy interface {
	Enqueue(taskType, handler string, payload any, queue chan<- Task) error
}

type ScheduledEnqueueStrategy struct {
	db *sql.DB
}

func (s ScheduledEnqueueStrategy) Enqueue(taskType, handler string, payload any, queue chan<- Task) error {
	task, err := initTask(taskType, handler, payload)
	if err != nil {
		return err
	}
	task.Status = "scheduled"

	err = PostTask(task, s.db)
	if err != nil {
		return err
	}

	queue <- task

	return nil
}

type OneTimeEnqueueStrategy struct {
	db *sql.DB
}

func (o OneTimeEnqueueStrategy) Enqueue(taskType, handler string, payload any, queue chan<- Task) error {
	task, err := initTask(taskType, handler, payload)
	if err != nil {
		return err
	}
	task.Status = "pending"

	err = PostTask(task, o.db)
	if err != nil {
		return err
	}

	queue <- task

	return nil
}
