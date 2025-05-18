package queue

import (
	"context"
	"log"
)

func StartWorker(ctx context.Context, taskQueue <-chan Task) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("worker shutting down")
		case task := <-taskQueue:
			payload := task.Payload
			handler := HandlerRegistry[task.Type]
			err := handler(payload)
			if err != nil {
				return err
			}
		}
	}
}
