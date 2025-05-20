package task_queue

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshalDbBackupPayload(t *testing.T) {
	payloads := []DbBackupPayload{
		{NextRunAt: time.Date(2025, time.November, 01, 0, 0, 0, 0, time.Local)},
	}
	testPayloads := [][]byte{}
	for _, p := range payloads {
		data, err := json.Marshal(p)
		if err != nil {
			t.Fatalf("error setting up test cases: %v", err)
		}
		testPayloads = append(testPayloads, data)
	}

	tests := []struct {
		payload  []byte
		expected DbBackupPayload
	}{
		{testPayloads[0], payloads[0]},
	}

	for _, tt := range tests {
		d, err := unmarshalDbBackupPayload(tt.payload)
		if err != nil {
			t.Errorf("incompatible struct types. want=%v, got=%v", tt.expected, err)
		}
		if d != tt.expected {
			t.Errorf("returned type incorrect. want=%v, got=%v", tt.expected, d)
		}
	}
}
