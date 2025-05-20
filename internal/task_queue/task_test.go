package task_queue

import (
	"database/sql"
	"encoding/json"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func TestInitTask_ValidInput(t *testing.T) {
	tests := []struct {
		taskType string
		handler  string
		payload  any

		expectedPayload any
	}{
		{"one_time", "send_email", EmailPayload{
			To:      "test",
			Subject: "test",
			Body:    "test",
		}, EmailPayload{
			To:      "test",
			Subject: "test",
			Body:    "test",
		},
		},
		{"one_time", "send_email", EmailPayload{}, EmailPayload{}},
	}

	for _, tt := range tests {
		expected, _ := json.Marshal(tt.expectedPayload)
		task, err := initTask(tt.taskType, tt.handler, tt.payload)
		if err != nil {
			t.Errorf("error inititalizing task. want=%v, got=%v", expected, err)
		}
		if string(expected) != string(task.Payload) {
			t.Errorf("error marshalling payload. want=%v, got=%v", expected, task.Payload)
		}
	}
}

func TestEnque(t *testing.T) {
	tests := []struct {
		taskType string
		handler  string
		payload  any

		expectedOneTimeLen   int
		expectedScheduledLen int
	}{
		{"one_time", "send_email", EmailPayload{}, 1, 0},
		{"scheduled", "send_email", EmailPayload{}, 0, 1},
	}

	db, err := initTestDb()
	if err != nil {
		t.Fatalf("error initializing db: %v", err)
	}

	for _, tt := range tests {
		oneTimeQueue := make(chan<- Task, 1)
		scheduledQueue := make(chan<- Task, 1)
		tp := TaskProducer{
			Db:             db,
			OneTimeQueue:   oneTimeQueue,
			ScheduledQueue: scheduledQueue,
		}
		err = tp.Enqueue(tt.taskType, tt.handler, tt.payload)
		if err != nil {
			t.Errorf("error enqueueing task. want=%v, %v, got=%v, %v",
				tt.expectedOneTimeLen, tt.expectedScheduledLen,
				len(tp.OneTimeQueue), len(tp.ScheduledQueue))
		}

		if len(tp.OneTimeQueue) != tt.expectedOneTimeLen || len(tp.ScheduledQueue) != tt.expectedScheduledLen {
			t.Errorf("tasks enqueued but len of queues wrong. want=%v, %v, got=%v, %v",
				tt.expectedOneTimeLen, tt.expectedScheduledLen,
				len(tp.OneTimeQueue), len(tp.ScheduledQueue))
		}
	}

	closeTestDb(db)
}

func initTestDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func closeTestDb(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func TestPostTask(t *testing.T) {
	db, err := initTestDb()
	if err != nil {
		t.Fatalf("error initializing test db: %v", err)
	}

	tasksBeforeTest, err := GetTasks(db)
	if err != nil {
		t.Fatalf("error getting all tasks from db: %v", err)
	}

	task, err := initTask(
		"one_time",
		"send_email",
		EmailPayload{
			To: "test", Subject: "test", Body: "test",
		},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = PostTask(task, db)
	if err != nil {
		t.Errorf("error posting task: expected: %v, got: %v", nil, err)
	}

	tasksAfterTest, err := GetTasks(db)
	if err != nil {
		t.Fatalf("error getting tasks: %v", err)
	}

	if len(tasksAfterTest) != len(tasksBeforeTest)+1 {
		t.Errorf("should be len(tasks)++ after test: expected: %v, got: %v", len(tasksBeforeTest)+1, len(tasksAfterTest))
	}

	DeleteTask(task.UUID, db)
	closeTestDb(db)
}
