package queue

var HandlerRegistry = map[string]func([]byte) error{}

func Register(taskType string, handler func([]byte) error) {
	HandlerRegistry[taskType] = handler
}
