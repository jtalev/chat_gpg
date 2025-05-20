package task_queue

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func TestInitTask_ValidInput(t *testing.T) {
	payload := EmailPayload{
		To:      "user@example.com",
		Subject: "Hello",
		Body:    "This is a test message.",
	}
	taskType := "send_email"
	maxRetries := 3

	task, err := initTask(taskType, payload, maxRetries)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Type != taskType {
		t.Errorf("expected task type %s, got %s", taskType, task.Type)
	}

	if task.Status != "pending" {
		t.Errorf("expected status 'pending', got %s", task.Status)
	}

	if task.MaxRetries != maxRetries {
		t.Errorf("expected maxRetries %d, got %d", maxRetries, task.MaxRetries)
	}

	if task.UUID == "" {
		t.Error("expected a UUID to be generated")
	}

	var decodedPayload EmailPayload
	err = json.Unmarshal(task.Payload, &decodedPayload)
	if err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if decodedPayload.To != payload.To {
		t.Errorf("expected payload To to be 'user@example.com', got %s", decodedPayload.To)
	}
	if decodedPayload.Subject != payload.Subject {
		t.Errorf("expected payload To to be 'Hello', got %s", decodedPayload.Subject)
	}
	if decodedPayload.Body != payload.Body {
		t.Errorf("expected payload To to be 'This is a test message.', got %s", decodedPayload.Body)
	}
}

func TestInitTask_InvalidInput(t *testing.T) {
	_, err := initTask("", nil, 1)
	if err == nil {
		t.Fatal("expected error for empty taskType and nil payload, got nil")
	}
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
		"send_email",
		EmailPayload{
			To: "test", Subject: "test", Body: "test",
		},
		0,
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

func TestEnqueue(t *testing.T) {
	db, err := initTestDb()
	if err != nil {
		t.Fatalf("error init db: %v", err)
	}
	p := TaskProducer{
		Db:    db,
		Queue: make(chan Task, 1),
	}
	qLengthBefore := len(p.Queue)
	taskType := "send_email"
	payload := EmailPayload{To: "test", Subject: "test", Body: "test"}
	maxRetries := 0
	err = p.Enqueue(taskType, payload, maxRetries)
	if err != nil {
		t.Errorf("error enqueueing task: expected: nil, got: %v", err)
	}
	qLengthAfter := len(p.Queue)
	if qLengthAfter != qLengthBefore+1 {
		t.Errorf("task not added to queue: expected length %v, got %v", qLengthBefore+1, qLengthAfter)
	}

	select {
	case task := <-p.Queue:
		DeleteTask(task.UUID, db)
	default:
		t.Errorf("no task found in queue")
	}
	closeTestDb(db)
}
