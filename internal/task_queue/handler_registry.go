package task_queue

import "database/sql"

type HandlerStrategy interface {
	ProcessTask(task Task, queue chan Task, db *sql.DB) error
}

var HandlerRegistry = map[string]HandlerStrategy{}

func RegisterTaskHandler(taskType string, handler HandlerStrategy) {
	HandlerRegistry[taskType] = handler
}
