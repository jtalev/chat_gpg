package queue

type Producer interface {
	Enqueue(taskType string, payload any, maxRetries int) error
}
