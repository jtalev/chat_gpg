package repository

import (
	"sort"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jtalev/chat_gpg/models"
)

func Test_GetLeaveRequests(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Mock expected SQL query
	rows := sqlmock.NewRows([]string{"request_id", "employee_id", "first_name", "last_name", "leave_type", "from_date", "to_date", "note", "is_approved"}).
		AddRow(1, "123", "John", "Doe", "Sick", "2024-01-01", "2024-01-05", "Feeling unwell", true).
		AddRow(3, "123", "John", "Doe", "Sick", "2024-01-01", "2024-01-05", "Feeling unwell", true).
		AddRow(2, "124", "Jane", "Smith", "Vacation", "2024-02-01", "2024-02-10", "Family trip", false)

	mock.ExpectQuery(`select \* from leave_request order by employee_id asc`).
		WillReturnRows(rows)

	// Call the function
	result, err := GetLeaveRequests(true, db)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Sort the result by EmployeeId to ensure correct order
	sort.Slice(result, func(i, j int) bool {
		return result[i].EmployeeId < result[j].EmployeeId
	})

	// Verify result
	if len(result) != 3 {
		t.Fatalf("Expected 3 leave requests, got %d", len(result))
	}

	// Ensure the first two are John and the third is Jane (order by employee_id)
	if result[0].FirstName != "John" || result[1].FirstName != "John" || result[2].FirstName != "Jane" {
		t.Fatalf("Expected leave requests for John and Jane in correct order, got %v", result)
	}

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unmet expectations: %v", err)
	}
}

func Test_PostLeaveRequest(t *testing.T) {
	// Initialize mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	// Define the leave request to be inserted
	leaveRequest := models.LeaveRequest{
		EmployeeId: "123",
		Type:       "Sick",
		From:       "2024-01-01",
		To:         "2024-01-05",
		Note:       "Feeling unwell",
	}

	// Mock expected SQL insert
	mock.ExpectExec(`INSERT INTO leave_request`).
		WithArgs(leaveRequest.EmployeeId, leaveRequest.Type, leaveRequest.From, leaveRequest.To, leaveRequest.Note).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mocking the result of the insert

	// Call the function
	result, err := PostLeaveRequest(leaveRequest, db)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify result
	if result.EmployeeId != leaveRequest.EmployeeId {
		t.Fatalf("Expected EmployeeId %v, got %v", leaveRequest.EmployeeId, result.EmployeeId)
	}
	if result.Type != leaveRequest.Type {
		t.Fatalf("Expected Type %v, got %v", leaveRequest.Type, result.Type)
	}

	// Ensure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unmet expectations: %v", err)
	}
}
